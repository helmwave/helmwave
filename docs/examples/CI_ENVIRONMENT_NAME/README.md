# CI_ENVIRONMENT_NAME [ [source] ](https://github.com/helmwave/helmwave/tree/main/docs/examples/CI_ENVIRONMENT_NAME)

#### Project structure

```
.
├── helmwave.yml.tpl
└── values
    ├── _.yml
    ├── prod.yml
    └── stage.yml
```

#### `helmwave.yml.tpl`

```yaml
project: my-project
version: 0.1.6

releases:
  - name: my-release
    chart: my-chart-repo/{{ env "CI_PROJECT_NAME" }}
    values:
      # Default
      - values/_.yml
      # For specific ENVIRONMENT
      - values/{{ env "CI_ENVIRONMENT_NAME" }}.yml
    options:
      install: true
      namespace: {{ env "CI_ENVIRONMENT_NAME" }}
```

#### `_.yml`

```yaml
image:
  tag: {{ env "CI_COMMIT_TAG" }}

podAnnotations:
  gitCommit: {{ requiredEnv "CI_COMMIT_SHORT_SHA" | quote }}
```

#### `prod.yml`

```yaml
replicaCount: 6
```

#### `stage.yml`

```yaml
replicaCount: 2
```