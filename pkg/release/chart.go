package release

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v4/pkg/action"
	chartiface "helm.sh/helm/v4/pkg/chart"
	chart "helm.sh/helm/v4/pkg/chart/v2"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"helm.sh/helm/v4/pkg/downloader"
	"helm.sh/helm/v4/pkg/getter"
	"helm.sh/helm/v4/pkg/helmpath"
	"helm.sh/helm/v4/pkg/registry"
	"helm.sh/helm/v4/pkg/repo/v1"
)

// Chart is a structure for chart download options.
//
//nolint:lll
type Chart struct {
	Name                  string `yaml:"name" json:"name" jsonschema:"required,description=Name of the chart,example=bitnami/nginx,example=oci://ghcr.io/helmwave/unit-test-oci"`
	CaFile                string `yaml:"ca_file" json:"ca_file" jsonschema:"description=Verify certificates of HTTPS-enabled servers using this CA bundle"`
	CertFile              string `yaml:"cert_file" json:"cert_file" jsonschema:"description=Identify HTTPS client using this SSL certificate file"`
	KeyFile               string `yaml:"key_file" json:"key_file" jsonschema:"description=Identify HTTPS client using this SSL key file"`
	Keyring               string `yaml:"keyring" json:"keyring" jsonschema:"description=Location of public keys used for verification"`
	RepoURL               string `yaml:"repo_url" json:"repo_url" jsonschema:"description=Chart repository url"`
	Username              string `yaml:"username" json:"username" jsonschema:"description=Chart repository username"`
	Password              string `yaml:"password" json:"password" jsonschema:"description=Chart repository password"`
	Version               string `yaml:"version" json:"version" jsonschema:"description=Chart version"`
	InsecureSkipTLSVerify bool   `yaml:"insecure" json:"insecure" jsonschema:"description=Connect to server with an insecure way by skipping certificate verification"`
	Verify                bool   `yaml:"verify" json:"verify" jsonschema:"description=Verify the provenance of the chart before using it"`
	PassCredentialsAll    bool   `yaml:"pass_credentials" json:"pass_credentials" jsonschema:"description=Pass credentials to all domains"`
	PlainHTTP             bool   `yaml:"plain_http" json:"plain_http" jsonschema:"description=Connect to server with plain http and not https,default=false"`
	SkipDependencyUpdate  bool   `yaml:"skip_dependency_update" json:"skip_dependency_update" jsonschema:"description=Skip updating and downloading dependencies,default=false"`
	SkipRefresh           bool   `yaml:"skip_refresh,omitempty" json:"skip_refresh,omitempty" jsonschema:"description=Skip refreshing repositories,default=false"`
}

// CopyOptions is a helper for copy options from Chart to ChartPathOptions.
func (c *Chart) CopyOptions(cpo *action.ChartPathOptions) {
	// I hate private field without normal New(...Options)
	cpo.CaFile = c.CaFile
	cpo.CertFile = c.CertFile
	cpo.KeyFile = c.KeyFile
	cpo.InsecureSkipTLSVerify = c.InsecureSkipTLSVerify
	cpo.PlainHTTP = c.PlainHTTP
	cpo.Keyring = c.Keyring
	cpo.Password = c.Password
	cpo.PassCredentialsAll = c.PassCredentialsAll
	cpo.RepoURL = c.RepoURL
	cpo.Username = c.Username
	cpo.Verify = c.Verify
	cpo.Version = c.Version
}

// UnmarshalYAML flexible config.
func (c *Chart) UnmarshalYAML(node *yaml.Node) error {
	type raw Chart
	var err error

	switch node.Kind {
	case yaml.ScalarNode, yaml.AliasNode:
		err = node.Decode(&(c.Name))
	case yaml.MappingNode:
		err = node.Decode((*raw)(c))
	default:
		err = ErrUnknownFormat
	}

	if err != nil {
		return fmt.Errorf("failed to decode chart %q from YAML at %d line: %w", node.Value, node.Line, err)
	}

	return nil
}

func (c *Chart) IsRemote() bool {
	return !helper.IsExists(filepath.Clean(c.Name))
}

// applyOCIRegistryClient swaps the shared TLS-only registry client for a per-release one when the
// chart is an OCI reference that asked for plain HTTP or a skipped TLS verification. helm v4's OCI
// getter uses an injected registry client as-is and never applies these flags to it, so the shared
// client can't reach such a registry. For every other chart the shared client stays the default and
// setClient is left untouched.
func (rel *config) applyOCIRegistryClient(setClient func(*registry.Client)) {
	c := rel.Chart()
	if !registry.IsOCI(c.Name) || (!c.PlainHTTP && !c.InsecureSkipTLSVerify) {
		return
	}

	rc, err := helper.NewRegistryClient(c.PlainHTTP, c.InsecureSkipTLSVerify)
	if err != nil {
		rel.Logger().WithError(err).
			Error("failed to build a dedicated OCI registry client, falling back to the shared TLS client")

		return
	}

	setClient(rc)
}

func (rel *config) LocateChartWithCache() (string, error) {
	if !rel.Chart().IsRemote() {
		return rel.Chart().Name, nil
	}

	ch, err := rel.findChartInHelmCache()
	if err == nil {
		rel.Logger().WithField("path", ch).Info("❎ found chart in helm cache, using it")

		return ch, nil
	}

	rel.Logger().WithError(err).Debug("haven't found chart in helm cache, need to download it")

	// nice action bro
	client := rel.newInstall()

	ch, err = client.ChartPathOptions.LocateChart(rel.Chart().Name, rel.Helm())
	if err != nil {
		return "", fmt.Errorf("failed to locate chart %s: %w", rel.Chart().Name, err)
	}

	return ch, nil
}

func (rel *config) getDownloader() downloader.ChartDownloader {
	settings := rel.Helm()
	client := rel.newInstall()

	return downloader.ChartDownloader{
		Getters: getter.All(settings),
		Options: []getter.Option{
			getter.WithPassCredentialsAll(client.ChartPathOptions.PassCredentialsAll),
			getter.WithTLSClientConfig(
				client.ChartPathOptions.CertFile,
				client.ChartPathOptions.KeyFile,
				client.ChartPathOptions.CaFile,
			),
			getter.WithInsecureSkipVerifyTLS(client.ChartPathOptions.InsecureSkipTLSVerify),
		},
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
		ContentCache:     settings.ContentCache,
		RegistryClient:   client.GetRegistryClient(),
	}
}

// Helm doesn't use its own charts cache, it only stores charts there. So we copypaste some code from
// *downloader.ChartDownloader to find already downloaded charts in our cache.
// We also check chart file digest in case of any collision.
func (rel *config) findChartInHelmCache() (string, error) {
	settings := rel.Helm()

	dl := rel.getDownloader()

	// The first return is the content digest, which we don't use: for the repository charts
	// below we compare against the digest from the repository index instead.
	_, u, err := dl.ResolveChartVersion(rel.Chart().Name, rel.Chart().Version)
	if err != nil {
		return "", NewChartCacheError(err)
	}

	name := filepath.Base(u.Path)
	if u.Scheme == registry.OCIScheme {
		idx := strings.LastIndexByte(name, ':')
		name = fmt.Sprintf("%s-%s.tgz", name[:idx], name[idx+1:])

		rel.Logger().Debug("digest validation is not supported for OCI charts, skipping it")

		chartFile := filepath.Join(settings.RepositoryCache, name)

		_, err := os.Stat(chartFile)
		if err != nil {
			return "", NewChartCacheError(err)
		}

		return chartFile, nil
	}

	chartFile := filepath.Join(settings.RepositoryCache, name)

	ch, err := rel.getChartRepoEntryFromIndex(u.String(), settings.RepositoryCache)
	if err != nil {
		return "", NewChartCacheError(err)
	}

	digest := ch.Digest
	hasher := sha256.New()

	f, err := os.Open(chartFile)
	if err != nil {
		return "", NewChartCacheError(err)
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = io.Copy(hasher, f)
	if err != nil {
		return "", NewChartCacheError(err)
	}

	hashSum := hex.EncodeToString(hasher.Sum(nil))

	if hashSum != digest {
		return "", NewChartCacheError(ErrDigestNotMatch)
	}

	return chartFile, nil
}

func (rel *config) getChartRepoEntryFromIndex(url, repositoryCache string) (*repo.ChartVersion, error) { //nolint:gocognit
	repoName := strings.SplitN(rel.Chart().Name, "/", 2)[0]
	repositoryIndex, err := repo.LoadIndexFile(filepath.Join(repositoryCache, helmpath.CacheIndexFile(repoName)))
	if err != nil {
		return nil, fmt.Errorf("no cached repo found: %w", err)
	}

	file, err := repo.LoadFile(rel.getDownloader().RepositoryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load repository file: %w", err)
	}

	var baseUrl string
	for _, r := range file.Repositories {
		if r.Name == repoName {
			baseUrl = r.URL

			break
		}
	}

	for _, e := range repositoryIndex.Entries {
		for _, ver := range e {
			for _, u := range ver.URLs {
				if baseUrl != "" {
					u, err = repo.ResolveReferenceURL(baseUrl, u)
					if err != nil {
						return nil, fmt.Errorf("failed to resolve reference URL %s: %w", u, err)
					}
				}
				if u == url {
					return ver, nil
				}
			}
		}
	}

	return nil, errors.New("repo not found")
}

func (rel *config) GetChart() (*chart.Chart, error) {
	ch, err := rel.LocateChartWithCache()
	if err != nil {
		return nil, err
	}

	c, err := loader.Load(ch)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart %s: %w", rel.Chart().Name, err)
	}

	if err := rel.chartCheck(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (rel *config) chartCheck(ch *chart.Chart) error {
	if req := ch.Metadata.Dependencies; req != nil {
		// action.CheckDependencies takes the version-agnostic chart.Dependency interface,
		// so widen the slice.
		deps := helper.SlicesMap(req, func(d *chart.Dependency) chartiface.Dependency {
			return d
		})

		if err := action.CheckDependencies(ch, deps); err != nil {
			return fmt.Errorf("failed to check chart %s dependencies: %w", ch.Name(), err)
		}
	}

	if ch.Metadata.Type != "" && ch.Metadata.Type != "application" {
		rel.Logger().Warnf("%s charts are not installable", ch.Metadata.Type)
	}

	if ch.Metadata.Deprecated {
		rel.Logger().Warnf("⚠️ Chart %s is deprecated. Please update your chart.", ch.Name())
	}

	return nil
}

func (rel *config) ChartDepsUpd() error {
	if rel.Chart().IsRemote() {
		rel.Logger().Info("❎ skipping updating dependencies for remote chart")

		return nil
	}

	if rel.Chart().SkipDependencyUpdate {
		rel.Logger().Info("❎ forced skipping updating dependencies for local chart")

		return nil
	}

	settings := rel.Helm()

	client := action.NewDependency()
	man := &downloader.Manager{
		Out:              log.StandardLogger().Writer(),
		ChartPath:        filepath.Clean(rel.Chart().Name),
		Keyring:          client.Keyring,
		RegistryClient:   helper.HelmRegistryClient,
		SkipUpdate:       rel.Chart().SkipRefresh,
		Getters:          getter.All(settings),
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
		ContentCache:     settings.ContentCache,
		Debug:            settings.Debug,
	}
	if client.Verify {
		man.Verify = downloader.VerifyAlways
	}

	if err := man.Update(); err != nil {
		return fmt.Errorf("failed to update %s chart dependencies: %w", rel.Chart().Name, err)
	}

	return nil
}

func (rel *config) DownloadChart(tmpDir string) error {
	if !rel.Chart().IsRemote() {
		rel.Logger().Info("❎ chart is local, skipping exporting")

		return nil
	}

	destDir := path.Join(tmpDir, "charts", rel.Uniq().String())
	if err := os.MkdirAll(destDir, 0o750); err != nil {
		return fmt.Errorf("failed to create temporary directory for chart: %w", err)
	}

	ch, err := rel.LocateChartWithCache()
	if err != nil {
		return err
	}

	return helper.CopyFile(ch, destDir)
}

func (rel *config) SetChartName(name string) {
	rel.lock.Lock()
	rel.ChartF.Name = name
	rel.lock.Unlock()
}
