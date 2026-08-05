# helmwave

Go CLI. Deploys Helm charts from `helmwave.yml`. Think docker-compose for helm.
Docs: https://docs.helmwave.app (source: `helmwave/docs` repo; reference pages: `yaml.md`, `cli.md`, `tpl.md`).

## Commands

- Build: `go build ./cmd/helmwave`
- Unit tests: `go test ./...` — never touches a cluster; cluster tests skip without `HELMWAVE_TEST_CLUSTER`
- Integration tests: KinD cluster from `tests/kind-config.yaml`, then `HELMWAVE_TEST_CLUSTER=~/.kube/config go test --tags=integration -timeout=20m ./...`
- Lint: `golangci-lint run` (config enables almost all linters; builds with `integration` tag)
- JSON schema: `go run ./cmd/helmwave schema` (from `jsonschema` struct tags; CI attaches to releases)

## Layout

- `cmd/helmwave` — entrypoint (urfave/cli v2)
- `pkg/action` — one file per CLI command (build, up, down, diff, ...)
- `pkg/plan` — core: plan build/import/apply
- `pkg/release`, `pkg/repo`, `pkg/registry` — helm releases, repos, OCI
- `pkg/template` — values templating (gomplate v3/v4, sprig, sops)
- `pkg/kubedog`, `pkg/monitor` — live resource tracking
- `tests/` — integration fixtures (`NN_helmwave.yml`)

## Rules

- Branching: Go changes go `feature/x` → `release/$SEMVER` → `main`. Non-Go changes PR straight to `main`.
- Conventional commits. Enforced by pre-commit (go-fmt, golangci-lint, go-mod-tidy, hadolint).
- Every user-facing change needs a changelog entry: `changie new` → `.changes/unreleased/`. CI fails without it.
- Tests: testify suites. Internal tests: `*_internal_test.go`. Integration tests: `//go:build integration` tag.
- Logging: logrus.
- Go version bump: update both `go.mod` and `Dockerfile` (`GOLANG_VERSION`).
