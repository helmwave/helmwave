<p align="center">
  <img  src="https://raw.githubusercontent.com/zhilyaev/helmwave/main/img/logo.png" style="max-height:100%;" height="220px" />
</p>

<h1 align="center"> HelmWave</h1>

<p align="center">
  <a href="https://github.com/helmwave/helmwave/actions?query=workflow%3Arelease"><img src="https://github.com/helmwave/helmwave/workflows/release/badge.svg" /></a>
  <a href="https://t.me/helmwave" ><img src="https://img.shields.io/badge/telegram-chat-179cde.svg?logo=telegram" /></a>
  <img alt="GitHub" src="https://img.shields.io/github/license/zhilyaev/helmwave">
  <img alt="GitHub tag (latest SemVer)" src="https://img.shields.io/github/v/tag/zhilyaev/helmwave?label=latest">
</p>


🌊 HelmWave is **[helm](https://github.com/helm/helm/)-native** tool for deploy your Helm Charts via **GitOps**.
HelmWave is like docker-compose for helm.

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

## Todo:

- buy a domain
- make docs

## 📥 Installation

### Download one of [releases](https://github.com/helmwave/helmwave/releases) ![GitHub all releases](https://img.shields.io/github/downloads/zhilyaev/helmwave/total)

```bash
$ wget -c https://github.com/helmwave/helmwave/releases/download/0.9.1/helmwave_0.9.1_linux_amd64.tar.gz -O - | tar -xz
$ mv helmwave /usr/local/bin/
```

### Install with go ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zhilyaev/helmwave)

```bash
$ GO111MODULE=on go get github.com/helmwave/helmwave/cmd/helmwave@0.9.1
```

### Run as a container ![Docker Pulls](https://img.shields.io/docker/pulls/diamon/helmwave)

```
$ docker run diamon/helmwave:0.9.1
$ docker run --entrypoint=ash -it --rm --name helmwave diamon/helmwave:0.9.1
```

### Build with ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zhilyaev/helmwave)

```bash
$ git clone git@github.com:zhilyaev/helmwave.git
$ cd helmwave
$ go build ./cmd/helmwave
$ mv helmwave /usr/local/bin/
```

## 🔰 Getting Started

Let's start with a simple **helmwave** and gradually improve it to fit your use-case!

Suppose the `helmwave.yml.tpl` representing the desired state of your helm releases looks like:

```yaml
project: my-project
version: 0.9.1


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
- HelmWave v0.8.3 – Kubedog рядом

## EN

- WIP

## [Documentation](https://zhilyaev.github.io/helmwave/)

### Annotations

> inspired by [werf annotations](https://werf.io/documentation/reference/deploy_annotations.html)

```yaml
  template:
    metadata:
      {{- with .Values.podAnnotations }}
      annotations:
        helmwave.dev/show-service-messages: true
      {{- toYaml . | nindent 8 }}
      {{- end }}
```

#### `helmwave.dev/track-termination-mode`

Defines a condition when helmwave should stop tracking of the resource:

- `WaitUntilResourceReady` (default) — the entire deployment process would monitor and wait for the readiness of the
  resource having this annotation. Since this mode is enabled by default, the deployment process would wait for all
  resources to be ready.
- `NonBlocking` — the resource is tracked only if there are other resources that are not yet ready.

#### helmwave.dev/fail-mode

Defines how helmwave will handle a resource failure condition which occured after failures threshold has been reached
for the resource during deploy process:

- `FailWholeDeployProcessImmediately` (default) — the entire deploy process will fail with an error if an error occurs
  for some resource.
- `HopeUntilEndOfDeployProcess` — when an error occurred for the resource, set this resource into the “hope” mode, and
  continue tracking other resources. If all remained resources are ready or in the “hope” mode, transit the resource
  back to “normal” and fail the whole deploy process if an error for this resource occurs once again.
- `IgnoreAndContinueDeployProcess` — resource errors do not affect the deployment process.

#### helmwave.dev/failures-allowed-per-replica

By default, one error per replica is allowed before considering the whole deployment process unsuccessful. This setting
defines a threshold of failures after which resource will be considered as failed and helmwave will handle this
situation using fail mode.

- NUMBER

#### helmwave.dev/log-regex

Defines a Re2 regex template that applies to all logs of all containers of all Pods owned by a resource with this
annotation. helmwave would show only those log lines that fit the specified regex template. By default, helmwave shows
all log lines.

- RE2_REGEX

#### helmwave.dev/log-regex-for-{container}

Defines a Re2 regex template that applies to all logs of specified container of all Pods owned by a resource with this
annotation. helmwave would show only those log lines that fit the specified regex template. By default, helmwave shows
all log lines.

- RE2_REGEX

#### helmwave.dev/skip-logs

Set to "true" to turn off printing logs of all containers of all Pods owned by a resource with this annotation. This
annotation is disabled by default.

- "true"|"false"

#### helmwave.dev/skip-logs-for-containers

Turn off printing logs of specified containers of all Pods owned by a resource with this annotation. This annotation is
disabled by default.

- string with `,` as a separator

#### helmwave.dev/show-logs-only-for-containers

Turn off printing logs of all containers except specified of all Pods owned by a resource with this annotation. This
annotation is disabled by default.

- string with `,` as a separator

#### helmwave.dev/show-service-messages

Set to "true" to enable additional real-time debugging info (including Kubernetes events) for a resource during
tracking. By default, helmwave would show these service messages only if the resource has failed the entire deploy
process.

- "true"|"false"

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
   0.9.1

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
   --tpl value                      Main tpl file (default: "helmwave.yml.tpl") [$HELMWAVE_TPL_FILE]
   --file value, -f value           Main yml file (default: "helmwave.yml") [$HELMWAVE_FILE, $HELMWAVE_YAML_FILE, $HELMWAVE_YML_FILE]
   --plan-dir value                 It keeps your state via planfile (default: ".helmwave/") [$HELMWAVE_PLAN_DIR]
   --tags value, -t value           It allows you choose releases for sync. Example: -t tag1 -t tag3,tag4 [$HELMWAVE_TAGS]
   --parallel helm install          It allows you call helm install in parallel mode  (default: true) [$HELMWAVE_PARALLEL]
   --log-format value               You can set: [ text | json | pad | emoji ] (default: "emoji") [$HELMWAVE_LOG_FORMAT]
   --log-level value                You can set: [ debug | info | warn  | fatal | panic | trace ] (default: "info") [$HELMWAVE_LOG_LEVEL, $HELMWAVE_LOG_LVL]
   --log-color                      Force color (default: true) [$HELMWAVE_LOG_COLOR]
   --kubedog                        Enable/Disable kubedog (default: true) [$HELMWAVE_KUBEDOG, $HELMWAVE_KUBEDOG_ENABLED]
   --kubedog-status-interval value  Interval of kubedog status messages (default: 5s) [$HELMWAVE_KUBEDOG_STATUS_INTERVAL]
   --kubedog-start-delay value      Delay kubedog start (default: 1s) [$HELMWAVE_KUBEDOG_START_DELAY]
   --kubedog-timeout value          Timout of kubedog multitrackers (default: 5m0s) [$HELMWAVE_KUBEDOG_TIMEOUT]
   --help, -h                       show help (default: false)
   --version, -v                    print the version (default: false)
```

### yml

Transform helmwave.yml.tpl to helmwave.yml

Suppose the `helmwave.yml.tpl` looks like:

```yaml
project: {{ env "CI_PROJECT_NAME" }}
version: 0.9.1


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
version: 0.9.1


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
  version: 0.9.1
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

Helmwave supports all built-in functions, [Sprig library](https://godoc.org/github.com/Masterminds/sprig), and several
advanced functions:

- `toYaml` marshals a map into a string
- `fromYaml` reads a golang string and generates a map
- `readFile` get file as string
- `hasKey` get true if field is exists
- `get` (Sprig's original `get` is available as `sprigGet`)
- `setValueAtPath` PATH NEW_VALUE traverses a golang map, replaces the value at the PATH with NEW_VALUE
- `requiredEnv` The requiredEnv function allows you to declare a particular environment variable as required for
  template rendering. If the environment variable is unset or empty, the template rendering will fail with an error
  message.

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
  version: 0.9.1
  
  
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
  version: 0.9.1
  
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







