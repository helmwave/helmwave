# 🧬 Full Configuration

## General

```yaml
project: my-project
version: 0.5.0

repositories: []
releases: []
```

## 🗄 Repository

```yaml
repositories:
- name: bitnami
  url: https://charts.bitnami.com/bitnami
  username: ""
  password: ""
  certfile: ""
  keyfile: ""
  cafile: ""
  insecureskiptlsverify: false
  force: false
```

## 🛥 Release

```yaml
releases:
- name: redis
  chart: bitnami/redis
  tags: []
  values: []
  store: 
  options:
    install: true
    devel: false
    namespace: b
    skipcrds: false
    timeout: 0s
    wait: false
    disablehooks: false
    dryrun: false
    force: false
    resetvalues: false
    reusevalues: false
    recreate: false
    maxhistory: 0
    atomic: false
    cleanuponfail: false
    subnotes: false
    description: ""
    postrenderer: null
    disableopenapivalidation: false
    chartpathoptions:
      cafile: ""
      certfile: ""
      keyfile: ""
      insecureskiptlsverify: false
      keyring: ""
      password: ""
      repourl: ""
      username: ""
      verify: false
      version: ""
```
