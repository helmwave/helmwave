# CI_COMMIT_SHORT_SHA [ [source] ](https://github.com/helmwave/helmwave/tree/main/docs/examples/CI_COMMIT_SHORT_SHA)

Project Structure

```
.
├── README.md
├── helmwave.yml.tpl
└── values.yml

```

```yaml
project: my-project # Имя проекта
version: 0.1.6 # Версия helmwave

releases:
  - name: my-release
    chart: my-chart-repo/my-app
    values:
      - values.yml
    options:
      install: true
      namespace: my-namespace

```

```yaml
image:
  tag: {{ env "CI_COMMIT_TAG" }}

podAnnotations:
  gitCommit: {{ requiredEnv "CI_COMMIT_SHORT_SHA" | quote }}
```