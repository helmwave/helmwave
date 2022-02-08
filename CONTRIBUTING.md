# Contributing Guidelines

The 🌊 HelmWave project accepts contributions via GitHub pull requests. \
This document outlines the process to help get your contribution accepted.

## Milestones

We use milestones to track progress of specific planned releases.

## Versioning

We use [semver](https://semver.org/) 

## Developing flow

*fix/feature-branch --> release-$SEMVER --> main*


**Example:**

- my-feature --> release-0.17.0 --> main
- my-fix --> release-0.17.1 --> main

### Non product update

When don't affect any `*.go` files we use [githubFlow](https://docs.github.com/en/get-started/quickstart/github-flow). 

`some branch --> main`

### Documentation update

After merging new release branch (minor or major) a new branch and PR in [helmwave/docs](https://github.com/helmwave/docs) will be created automatically. You will need to update documentation (if necessary) and merge this PR into main branch.

## How to build?

https://helmwave.github.io/docs/0.17.x/install/#compile-from-source

### Pre commit

We use https://pre-commit.com for git hooks