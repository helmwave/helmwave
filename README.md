<p align="center">
  <img  src="https://raw.githubusercontent.com/helmwave/logo/main/logo.png" style="max-height:100%;" height="220px" />
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


## 📥 Installation

### Download one of [releases](https://github.com/helmwave/helmwave/releases) ![GitHub all releases](https://img.shields.io/github/downloads/zhilyaev/helmwave/total)

```bash
$ wget -c https://github.com/helmwave/helmwave/releases/download/0.9.1/helmwave_0.9.1_linux_amd64.tar.gz -O - | tar -xz
$ mv helmwave /usr/local/bin/
```


### Run as a container ![Docker Pulls](https://img.shields.io/docker/pulls/diamon/helmwave)

```
$ docker run diamon/helmwave
$ docker run --entrypoint=ash -it --rm --name helmwave diamon/helmwave
```

## 📖 [Documentation](https://helmwave.github.io/)

Documentation is hosted live at https://helmwave.github.io

