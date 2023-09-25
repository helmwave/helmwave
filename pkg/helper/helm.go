package helper

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/fs"
	"net/url"
	"os"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	"k8s.io/apimachinery/pkg/version"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

//nolint:gochecknoglobals // TODO: get rid of globals
var (
	// Helm is an instance of helm CLI.
	Helm = helm.New()

	// Default logLevel for helm logs.
	helmLogLevel = log.Debugf

	// HelmRegistryClient  is an instance of helm registry client.
	HelmRegistryClient *registry.Client

	HelmFS fsimpl.WriteableFS
)

func init() {
	var err error
	HelmRegistryClient, err = registry.NewClient(
		registry.ClientOptDebug(Helm.Debug),
		registry.ClientOptWriter(log.StandardLogger().Writer()),
		registry.ClientOptCredentialsFile(Helm.RegistryConfig),
	)
	if err != nil {
		log.Fatal(err)
	}

	helmROFS, err := filefs.New(&url.URL{Scheme: "file", Path: "/"})
	if err != nil {
		log.Fatal(err)
	}
	HelmFS = helmROFS.(fsimpl.WriteableFS) //nolint:forcetypeassert
}

func wrapConfigFn(client *rest.Config) *rest.Config {
	client.QPS = 100   // default is 5.0
	client.Burst = 100 // default is 10

	return client
}

// NewCfg creates helm internal configuration for provided namespace and kubecontext.
func NewCfg(ns, kubecontext string) (*action.Configuration, error) {
	cfg := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER") // TODO: get rid of getenv in runtime
	config := genericclioptions.NewConfigFlags(true)
	config.WrapConfigFn = wrapConfigFn
	config.Namespace = &ns
	if kubecontext != "" {
		config.Context = &kubecontext
	} else {
		config.Context = &Helm.KubeContext
	}

	if Helm.Debug {
		helmLogLevel = log.Infof
	}
	err := cfg.Init(config, ns, helmDriver, helmLogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm configuration for %s namespace: %w", ns, err)
	}

	cfg.RegistryClient = HelmRegistryClient

	return cfg, nil
}

// NewHelm is a hack to create an instance of helm CLI and specifying namespace without environment variables.
func NewHelm(ns string) (*helm.EnvSettings, error) {
	env := helm.New()
	flagset := &pflag.FlagSet{}
	env.AddFlags(flagset)
	flag := flagset.Lookup("namespace")

	if err := flag.Value.Set(ns); err != nil {
		return nil, fmt.Errorf("failed to set namespace %s for helm: %w", ns, err)
	}

	return env, nil
}

// GetKubernetesVersion returns kubernetes server version.
//
//nolint:wrapcheck
func GetKubernetesVersion(cfg *action.Configuration) (*version.Info, error) {
	clientSet, err := cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	return clientSet.Discovery().ServerVersion()
}

type helmFSGetter struct {
	baseFS fs.FS
}

//nolint:wrapcheck
func (g *helmFSGetter) Get(u string, _ ...getter.Option) (*bytes.Buffer, error) {
	f, err := g.baseFS.Open(u)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func newHelmFSGetterWrapped(baseFS fs.FS) getter.Constructor {
	return func(_ ...getter.Option) (getter.Getter, error) {
		return &helmFSGetter{baseFS: baseFS}, nil
	}
}

func GetHelmFSProvider(baseFS fs.FS) getter.Providers {
	return []getter.Provider{
		{
			Schemes: []string{"", "file"},
			New:     newHelmFSGetterWrapped(baseFS),
		},
	}
}

func LocateChart(c *action.ChartPathOptions, baseFS fs.FS, name string, settings *helm.EnvSettings) (string, error) {
	if IsExists(baseFS, name) {
		return name, nil
	}

	return c.LocateChart(name, settings) //nolint:wrapcheck
}

func ChartFSLoad(baseFS fs.StatFS, path string) (*chart.Chart, error) {
	fi, err := baseFS.Stat(path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	if fi.IsDir() {
		return chartFSDirLoad(baseFS, path)
	}

	return chartFSFileLoad(baseFS, path)
}

//nolint:wrapcheck,funlen
func chartFSDirLoad(baseFS fs.FS, path string) (*chart.Chart, error) {
	c := &chart.Chart{}

	// rules := ignore.Empty()
	// ifile := filepath.Join(topdir, ignore.HelmIgnore)
	// if _, err := os.Stat(ifile); err == nil {
	//	r, err := ignore.ParseFile(ifile)
	//	if err != nil {
	//		return c, err
	//	}
	//	rules = r
	// }
	// rules.AddDefaults()

	files := []*loader.BufferedFile{}

	walk := func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if de.IsDir() {
			// Directory-based ignore rules should involve skipping the entire
			// contents of that directory.
			// if rules.Ignore(n, fi) {
			//	return filepath.SkipDir
			// }
			return nil
		}

		// If a .helmignore file matches, skip this file.
		// if rules.Ignore(n, fi) {
		//	return nil
		// }

		fi, err := de.Info()
		if err != nil {
			return err
		}

		// Irregular files include devices, sockets, and other uses of files that
		// are not regular files. In Go they have a file mode type bit set.
		// See https://golang.org/pkg/os/#FileMode for examples.
		if !fi.Mode().IsRegular() {
			return fmt.Errorf("cannot load irregular file %s as it has file mode type bits set", path)
		}

		f, err := baseFS.Open(path)
		if err != nil {
			return errors.Wrapf(err, "error reading %s", path)
		}
		defer f.Close() //nolint:errcheck

		buf := &bytes.Buffer{}
		_, err = buf.ReadFrom(f)
		if err != nil {
			return errors.Wrapf(err, "error reading %s", path)
		}
		data := buf.Bytes()

		data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

		files = append(files, &loader.BufferedFile{Name: path, Data: data})

		return nil
	}
	if err := fs.WalkDir(baseFS, path, walk); err != nil {
		return c, err
	}

	return loader.LoadFiles(files)
}

//nolint:wrapcheck
func chartFSFileLoad(baseFS fs.FS, path string) (*chart.Chart, error) {
	raw, err := baseFS.Open(path)
	if err != nil {
		return nil, err
	}

	err = ensureArchive(raw)
	if err != nil {
		return nil, err
	}
	_ = raw.Close()

	raw, err = baseFS.Open(path)
	if err != nil {
		return nil, err
	}
	defer raw.Close() //nolint:errcheck

	return loader.LoadArchive(raw)
}

//nolint:wrapcheck
func ensureArchive(raw fs.File) error {
	r, err := gzip.NewReader(raw)
	if err != nil {
		return err
	}

	return r.Close()
}
