# 🔰 Store

It allows pass you custom values to render release.

#### `helmwave.yml`:

```yaml 
project: my-project
version: 0.5.0

releases:
  - name: backend
    chart: my/backend
    options:
      install: true
    store:
      secret:
        type: vault
        path: secret/my/frontend
    values:
      - my-custom-values.yml

  - name: frontend
    chart: my/frontend
    options:
      install: true
    store:
      secret:
        type: vault
        path: secret/my/frontend
    values:
      - my-custom-values.yml
```

#### `my-custom-values.yml`:

```yaml
secretForApp:
  kind: {{ .Release.Store.secret.type }}
  path: {{ .Release.Store.secret.path | quote }}
```
