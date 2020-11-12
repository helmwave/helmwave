<p align="center">
  <img  src="logo.svg" style="max-height:100%;" height="220px" />
</p>

<h1 align="center"> HelmWave </h1>

<p align="center">
  <img  src="https://github.com/zhilyaev/helmwave/workflows/.github/workflows/release.yml/badge.svg" />
</p>


🏖 HelmWave is **[helm](https://github.com/helm/helm/)-native** tool for deploy your chart.
It helps you compose your helm releases!


> Inspired by the [helmfile](https://github.com/roboll/helmfile)

## Сomparison
 item        | 🌊 HelmWave | helmfile 
-------------| ------------| ----------- 
Docker image | 23 mb       | 190 mb
Without helm binary |✅|❌
All options helm|✅| partially
Helm version | only 3 | 2 and 3
Parallel helm install/upgarde |✅|❌
Repository Skipping|✅|❌
Install only needs repositories|✅|❌
Tags|✅| labels
Planfile|✅|❌
Sprig | ✅|✅
Call helm | via Golang Module | shell executor


## 📥 Installation

- Download one of [releases](https://github.com/zhilyaev/helmwave/releases)
    - `$ wget -c https://github.com/zhilyaev/helmwave/releases/download/0.1.6/helmwave-0.1.6-linux-amd64.tar.gz -O - | tar -xz && cp -f helmwave /usr/local/bin/`
- Run as a container
    - `$ docker run diamon/helmwave:0.1.6`
    - `$ docker run --entrypoint=ash -it --rm --name helmwave diamon/helmwave:0.1.6`

### Build

> golang:1.15

```bash
$ export GO111MODULE=on
$ git clone git@github.com:zhilyaev/helmwave.git $GOPATH/src/github.com/zhilyaev/helmwave
$ cd $GOPATH/src/github.com/zhilyaev/helmwave
$ go build github.com/zhilyaev/helmwave/cmd/helmwave
$ mv helmwave /usr/local/bin 
```

## 🔰 Getting Started 

Let's start with a simple **helmwave** and gradually improve it to fit your use-case!

Suppose the `helmwave.yml.tpl` representing the desired state of your helm releases looks like:

```yaml
project: my-project
version: 0.1.6


repositories:
  - name: bitnami
    url: https://charts.bitnami.com/bitnami


.options: &options
  install: true
  namespace: my-namespace


releases:
  - name: redis-a
    chart: bitnami/redis
    options:
      <<: *options

  - name: redis-b
    chart: bitnami/redis
    options:
      <<: *options
```

```shell script
$ helmwave deploy
```

Congratulations! 
```shell script
$ helm list -n my-namespace
NAME                    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
redis-a                 my-namespace    2               2020-10-31 17:05:35.829766 +0300 MSK    deployed        redis-11.2.3            6.0.9      
redis-b                 my-namespace    16              2020-10-31 17:05:39.437556 +0300 MSK    deployed        redis-11.2.3            6.0.9  

$ k get po -n my-namespace                                                                                                                            (k8s-sbs-dev/stage)
NAME               READY   STATUS    RESTARTS   AGE
redis-a-master-0   1/1     Running   0          64s
redis-a-slave-0    1/1     Running   0          31s
redis-a-slave-1    1/1     Running   0          62s
redis-b-master-0   1/1     Running   0          59s
redis-b-slave-0    1/1     Running   0          32s
redis-b-slave-1    1/1     Running   0          51s
```

## 🔰 Tags

Suppose the `helmwave.yml.tpl` looks like:

```yaml
project: my-project
version: 0.1.6


repositories:
  - name: bitnami
    url: https://charts.bitnami.com/bitnami


.options: &options
  install: true
  namespace: my-namespace


releases:
  - name: redis-a
    chart: bitnami/redis
    tags:
      - a
      - redis
    options:
      <<: *options

  - name: redis-b
    chart: bitnami/redis
    tags:
      - b
      - redis
    options:
      <<: *options

  - name: memcached-a
    chart: bitnami/memcached
    tags:
      - a
      - memcached
    options:
      <<: *options

  - name: memcached-b
    chart: bitnami/memcached
    tags:
      - b
      - memcached
    options:
      <<: *options
```

This command will deploy only `redis-a` & `memcached-a`
```shell script
$ helmwave -t a deploy
```

This command will deploy only `redis-a` & `redis-b`
```shell script
$ helmwave -t redis deploy
```


This command will deploy only `redis-b`
```shell script
$ helmwave -t redis,b deploy
```

## 🛠 CLI Reference

```console
NAME:
   🌊 HelmWave - composer for helm

USAGE:
   helmwave [global options] command [command options] [arguments...]

VERSION:
   0.1.6

DESCRIPTION:
   🏖 This tool helps you compose your helm releases!

AUTHOR:
   💎 Dmitriy Zhilyaev <helmwave+zhilyaev.dmitriy@gmail.com>

COMMANDS:
   render, r                        📄 Render tpl -> yml
   planfile, p, plan                📜 Generate planfile
   repos, rep, repo                 🗄 Sync repositories
   deploy, d, apply, sync, release  🛥 Deploy your helmwave!
   help, h                          🚑 Help me!

GLOBAL OPTIONS:
   --tpl value                 Main tpl file (default: "helmwave.yml.tpl") [$HELMWAVE_TPL_FILE]
   --file value, -f value      Main yml file (default: "helmwave.yml") [$HELMWAVE_FILE, $HELMWAVE_YAML_FILE, $HELMWAVE_YML_FILE]
   --planfile value, -p value  (default: "helmwave.plan") [$HELMWAVE_PLANFILE]
   --tags value, -t value      Chose tags: -t tag1 -t tag3,tag4 [$HELMWAVE_TAGS]
   --debug, -d                 Debug helmwave (default: false) [$HELMWAVE_DEBUG]
   --parallel                  Parallel mode (default: true) [$HELMWAVE_PARALLEL]
   --version, -v               print the version (default: false)

```

### Render, r
Transform helmwave.yml.tpl to helmwave.yml

Suppose the `helmwave.yml.tpl` looks like:

```yaml
project: {{ env "CI_PROJECT_NAME" }}
version: 0.1.6


repositories:
- name: your-private-git-repo-hosted-charts
  url: https://{{ env "GITHUB_TOKEN"}}@raw.githubusercontent.com/foo/bar/master/


.options: &options
  install: true
  namespace: {{ env "NS" }}


releases:
  - name: redis-a
    chart: bitnami/redis
    options:
      <<: *options
```

This command will render `helmwave.yml.tpl` to `helmwave.yml`
```shell script
$ export NS=stage
$ export CI_PROJECT_NAME=my-project
$ export GITHUB_TOKEN=my-secret-token
$ helmwave render
```

Once applied, your `helmwave.yml` will look like:

```yaml
project: my-project
version: 0.1.6


repositories:
- name: your-private-git-repo-hosted-charts
  url: https://my-secret-token@raw.githubusercontent.com/foo/bar/master/


.options: &options
  install: true
  namespace: stage


releases:
  - name: redis-a
    chart: bitnami/redis
    options:
      <<: *options
```

### Planfile, plan

This command will generate helmwave.plan.

`helmwave.plan` is an object save to yaml.

<details>
  <summary>helmwave.plan looks like</summary>
  
  ```yaml
  project: my-project
  version: 0.1.6
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
  releases:
  - name: redis-a
    chart: bitnami/redis
    tags: []
    values: []
    options:
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
      install: true
      devel: false
      namespace: my-namespace
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
  - name: redis-b
    chart: bitnami/redis
    tags: []
    values: []
    options:
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
      install: true
      devel: false
      namespace: my-namespace
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
  ```
</details>



## 📄 Templating
HelmWave uses [Go templates](https://godoc.org/text/template) for templating.

Helmwave supports all built-in functions, [Sprig library](https://godoc.org/github.com/Masterminds/sprig), and several advanced functions:
- `toYaml` marshals a map into a string
- `fromYaml` reads a golang string and generates a map
- `readFile` get file as string
- `get` (Sprig's original `get` is available as `sprigGet`)
- `setValueAtPath` PATH NEW_VALUE traverses a golang map, replaces the value at the PATH with NEW_VALUE
- `requiredEnv` The requiredEnv function allows you to declare a particular environment variable as required for template rendering. If the environment variable is unset or empty, the template rendering will fail with an error message.


#### readFile

<details>
  <summary>my-releases.yml</summary>
  
  ```yaml
releases:
  - name: redis
    repo: bitnami
  - name: memcached
    repo: bitnami
  ```
</details>

<details>
  <summary>helmwave.yml.tpl</summary>
  
  ```yaml
  project: my
  version: 0.1.6
  
  
  repositories:
    - name: bitnami
      url: https://charts.bitnami.com/bitnami
  
  
  .global: &global
    install: true
  
  
  releases:
  {{- with readFile "my-releases.yml" | fromYaml | get "releases" }}
    {{- range $v := . }}
    - name: {{ $v | get "name" }}
      chart: {{ $v | get "repo" }}/{{ $v | get "name" }}
      options:
        <<: *global
    {{- end }}
  {{- end }}
  ``` 
  
</details>

```bash
$ helmwave render
```

<details>
  <summary>helmwave.yml</summary>
  
  ```yaml
  project: my
  version: 0.1.6
  
  repositories:
    - name: bitnami
      url: https://charts.bitnami.com/bitnami
  
  .global: &global
    install: true
  
  releases:
    - name: redis
      chart: bitnami/redis
      options:
        <<: *global
    - name: memcached
      chart: bitnami/memcached
      options:
        <<: *global
  ``` 
  
</details>



## 🚀 Features

- [x] Helm release is working via official golang module [helm](https://github.com/helm/helm/tree/master/pkg)
- [x] Minimal docker image without helm binary or another plugins
- [x] Support helm repositories
- [x] Parallel applying helm releases
- [x] Templating helmwave.yml
- [x] Templating values


## 🛸 Coming soon...
- [ ] Dependencies helm release
- [ ] OCI, testing...
- [ ] Force Update repositories
- [ ] Formatting output
- [ ] Applying from planfile
- [ ] Dependencies helmwave
- [ ] More templating functions


## 🧬 Full Configuration

### 🗄 Repository

```yaml
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


### 🛥 Release

```yaml
- name: redis
  chart: bitnami/redis
  tags: []
  values: []
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
    maxhistory: 0 # infinity
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
