<p align="center">
  <img  src="https://raw.githubusercontent.com/zhilyaev/helmwave/main/img/logo.png" style="max-height:100%;" height="220px" />
</p>

<h1 align="center"> HelmWave</h1>

<p align="center">
  <a href="https://github.com/zhilyaev/helmwave/actions?query=workflow%3Arelease"><img src="https://github.com/zhilyaev/helmwave/workflows/release/badge.svg" /></a>
  <a href="https://t.me/helmwave" ><img src="https://img.shields.io/badge/telegram-chat-179cde.svg?logo=telegram" /></a>
  <img alt="GitHub" src="https://img.shields.io/github/license/zhilyaev/helmwave">
  <img alt="GitHub tag (latest SemVer)" src="https://img.shields.io/github/v/tag/zhilyaev/helmwave?label=latest">
</p>


🌊 HelmWave is **[helm](https://github.com/helm/helm/)-native** tool for deploy your Helm Charts via **GitOps**. HelmWave is like docker-compose for helm.



- Keep a directory of chart value files and maintain changes in version control.
- Apply CI/CD to configuration changes
- Template values
- Aggregate your application


## Comparison
 🚀 Features  | 🌊 HelmWave   | helmfile 
-------------| :------------:|:-----------:
Docker | ![Docker Image Size helmwave (latest by date)](https://img.shields.io/docker/image-size/diamon/helmwave) | ![Docker Image Size helmfile (latest by date)](https://img.shields.io/docker/image-size/chatwork/helmfile)
[Kubedog](https://github.com/werf/kubedog) |✅|❌
Without helm binary |✅|❌
All options helm|✅|partially
Helm 3 |✅|✅
Helm 2 |❌|✅
Parallel helm install/upgrade |✅|❌
Repository Skipping|✅|❌
Install only needs repositories|✅|❌
Tags|✅| You can use labels
Store|✅| You can use labels
Planfile|✅|❌
Sprig | ✅|✅
Call helm | via Golang Module | Shell Executor
Speed of deploy <sup>[*]</sup> | 10 sec | 2 min

`*` - WIP

## 📥 Installation ![GitHub all releases](https://img.shields.io/github/downloads/zhilyaev/helmwave/total) ![Docker Pulls](https://img.shields.io/docker/pulls/diamon/helmwave)

- Download one of [releases](https://github.com/zhilyaev/helmwave/releases)
    - `$ wget -c https://github.com/zhilyaev/helmwave/releases/download/0.8.0/helmwave-0.8.0-linux-amd64.tar.gz -O - | tar -xz && cp -f helmwave /usr/local/bin/`
- Run as a container
    - `$ docker run diamon/helmwave:0.8.0`
    - `$ docker run --entrypoint=ash -it --rm --name helmwave diamon/helmwave:0.8.0`

### Build ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zhilyaev/helmwave)

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
version: 0.8.0


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
redis-a                 my-namespace    1               2020-10-31 17:05:35.829766 +0300 MSK    deployed        redis-11.2.3            6.0.9      
redis-b                 my-namespace    1               2020-10-31 17:05:39.437556 +0300 MSK    deployed        redis-11.2.3            6.0.9  

$ k get po -n my-namespace                                                                                                                         
NAME               READY   STATUS    RESTARTS   AGE
redis-a-master-0   1/1     Running   0          64s
redis-a-slave-0    1/1     Running   0          31s
redis-a-slave-1    1/1     Running   0          62s
redis-b-master-0   1/1     Running   0          59s
redis-b-slave-0    1/1     Running   0          32s
redis-b-slave-1    1/1     Running   0          51s
```

## Articles

### RU
- [HelmWave v0.5.0 – GitOps для твоего Kubernetes](https://habr.com/ru/post/532596/)
-  HelmWave v0.8.0 – Kubedog рядом

## EN
- WIP

## [Documentation](https://zhilyaev.github.io/helmwave/)


### Examples
  - [How to pass `image.tag` to release?](docs/examples/CI_COMMIT_TAG/README.md)
  - [How to pass `podAnnotations.gitCommit` to release?](docs/examples/CI_COMMIT_SHORT_SHA/README.md)
  - [How to use environments for your release?](docs/examples/CI_ENVIRONMENT_NAME/README.md)

### [🔰 Store](docs/store.md)
It allows pass you custom values to render release. 

### [🔰 Tags](docs/tags.md)
Use tags for choose specific releases

### [🧬 Full `helmwave.yml` config](docs/helmwave.yml.md)
All Options

## 🛠 CLI Reference

```console                                                                                                                                     (k8s-sbs/stage)
NAME:
   helmwave - composer for helm

USAGE:
   helmwave [global options] command [command options] [arguments...]

VERSION:
   0.8.0

DESCRIPTION:
   🏖 This tool helps you compose your helm releases!

AUTHOR:
   💎 Dmitriy Zhilyaev <helmwave+zhilyaev.dmitriy@gmail.com>

COMMANDS:
   yml                           📄 Render helmwave.yml.tpl -> helmwave.yml
   planfile, plan                📜 Generate planfile to plandir
   deploy, apply, sync, release  🛥 Deploy your helmwave!
   help, h                       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --tpl value              Main tpl file (default: "helmwave.yml.tpl") [$HELMWAVE_TPL_FILE]
   --file value, -f value   Main yml file (default: "helmwave.yml") [$HELMWAVE_FILE, $HELMWAVE_YAML_FILE, $HELMWAVE_YML_FILE]
   --plan-dir value         It keeps your state via planfile (default: ".helmwave/") [$HELMWAVE_PLAN_DIR]
   --tags value, -t value   It allows you choose releases for sync. Example: -t tag1 -t tag3,tag4 [$HELMWAVE_TAGS]
   --parallel helm install  It allows you call helm install in parallel mode  (default: true) [$HELMWAVE_PARALLEL]
   --log-format value       You can set: [ text | json | pad | emoji ] (default: "emoji") [$HELMWAVE_LOG_FORMAT]
   --log-level value        You can set: [ debug | info | warn | panic | fatal | trace ] (default: "info") [$HELMWAVE_LOG_LEVEL, $HELMWAVE_LOG_LVL]
   --log-color              Force color (default: true) [$HELMWAVE_LOG_COLOR]
   --kubedog                Enable/Disable kubedog (default: true) [$HELMWAVE_KUBEDOG]
   --help, -h               show help (default: false)
   --version, -v            print the version (default: false)
```

### yml

Transform helmwave.yml.tpl to helmwave.yml

Suppose the `helmwave.yml.tpl` looks like:

```yaml
project: {{ env "CI_PROJECT_NAME" }}
version: 0.8.0


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
$ helmwave yml
```

Once applied, your `helmwave.yml` will look like:

```yaml
project: my-project
version: 0.8.0


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

### planfile, plan

This command will generate helmwave.plan.

`helmwave.plan` is an object save to yaml.

<details>
  <summary>helmwave.plan looks like</summary>
  
  ```yaml
  project: my-project
  version: 0.8.0
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
- `hasKey` get true if field is exists
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
  version: 0.8.0
  
  
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
$ helmwave yml
```

<details>
  <summary>helmwave.yml</summary>
  
  ```yaml
  project: my
  version: 0.8.0
  
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







