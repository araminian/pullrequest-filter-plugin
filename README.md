# Pull Request Filter : ArgoCD ApplicationSet Generator Plugin

## Overview

[ArgoCD Plugin Generator](https://github.com/binboum/argo-cd/blob/1f7ff874c4961edd7e50595602c174936619299c/docs/operator-manual/applicationset/Generators-Plugin.md) was introduced recently.

This feature let's us to create our own `generator` and extend the functionality of ArgoCD ApplicationSet. The plugin can be written in any language and it should be able to run in a container.

It would make adding new feature to ArgoCD ApplicationSet much easier. `plugins` can be used in `matrix` and `merge` generators to extend functionality of other generators.

For more information the [sample repo](https://github.com/argoproj-labs/applicationset-hello-plugin) can be checked.

## Background

`PullRequest` generator is a plugin that can generate an `Application` for each `PR` in a `Github` repository. It checked for `labels` in `PRs` and if the labels in `PR` matches with the `labels` in `ApplicationSet` it will create an `Application` for that `PR`.

```yaml
generators:
  - pullRequest:
      github:
      # The GitHub organization or user.
        owner: araminian
          # The Github repository
        repo: BRApplication
          # Labels is used to filter the PRs that you want to target. (optional)
        labels:
          - advertisers
        requeueAfterSeconds: 180
```

For instance here, We are generating an `Application` for each `PR` in `BRApplication` repository that has `advertisers` label. The `labels` is optional and if it is not provided it will generate an `Application` for each `PR` in the repository.

## Our Problem

In our setup we have an `Preview Environment ApplicationSet` for each `backend service` which checks if there is a `PR` with `[ServiceName]` as label and if there is a `PR` with that label it will create an `Application` for that `PR` and deploy it to `Preview Environment`.

Like following example that deploy a preview environment for `advertisers` service if there is a `PR` with `advertisers` label.


```yaml
generators:
  - pullRequest:
      github:
      # The GitHub organization or user.
        owner: araminian
          # The Github repository
        repo: BRApplication
          # Labels is used to filter the PRs that you want to target. (optional)
        labels:
          - advertisers
        requeueAfterSeconds: 180
```

A `PR` can have more than one `service` label and we want to deploy a `Preview Environment` for each `service` label in that `PR`. For instance if a `PR` has `advertisers` and `campaigns` labels we want to deploy a `Preview Environment` for each of them.

The `label` assignment is done by `Github Actions`. If the `PR` has a change related to a service it will add that service name as a label to the `PR`.

There are some scenarios that we don't want a preview environment to be deployed for a `PR`:
- Saving resources in non-working hours
- No need for preview environment for a `PR` that has a change in `README.md` or documents
- No preview environment needed

## Solution

Our solution relies on new `plugin generator`. The idea is to check the `PR` labels and if the `PR` has a label defined like `no-preview`, it will generate `Application` for that `PR` and services but not point to the real `K8S manifest` instead it will point to a `blackhole` directory that has no `K8S manifest` in it.

The reason for using `blackhole` direcotry is that we can't limit `PullRequest generator` functionallity, but we can limit the `path` that is used for finding `manifests`.

The `blackhole` directory should be created in `gitops` repo. [Here](https://github.com/araminian/BRApplication/tree/main/manifests/blackhole) is an example of `blackhole` directory.


Let see an exmaple of `ApplicationSet` that uses `PullRequestGenerator` and `applicationset-pr-filter-plugin` plugin:


```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: br-pr
  namespace: argocd
spec:
  goTemplate: true
  generators:
    - matrix:
        generators:
          - pullRequest:
              github:
                # The GitHub organization or user.
                owner: araminian
                # The Github repository
                repo: BRApplication
                # Labels is used to filter the PRs that you want to target. (optional)
                labels:
                  - advertisers
              requeueAfterSeconds: 180
          - plugin:
              configMapRef:
                name: applicationset-pr-filter-plugin
              input:
                parameters:
                  labels: "{{.labels}}"
                  excludeLabel: "no-preview"
                  path: "manifests/app/"
                  blackHole: "manifests/blackhole/"
                  number: "{{.number}}"
  template:
    metadata:
      name: "br-{{.branch}}-{{.number}}"
    spec:
      source:
        repoURL: "https://github.com/araminian/BRApplication.git"
        targetRevision: HEAD
        path: "{{.generatedPath}}"
        directory:
          include: '{*.yml,*.yaml}'
      project: "default"
      destination:
        server: https://kubernetes.default.svc
        namespace: "default"
      syncPolicy:
        syncOptions:
          - CreateNamespace=true
```


The important part is the `plugin` generator:

```yaml
          - plugin:
              configMapRef:
                name: applicationset-pr-filter-plugin
              input:
                parameters:
                  labels: "{{.labels}}"
                  excludeLabel: "no-preview"
                  path: "manifests/app/"
                  blackHole: "manifests/blackhole/"
                  number: "{{.number}}"

```

The `configMapRef` is defining the `configMap` that `ApplicationSet Controller` will use to find the `plugin` and how it can call it.

The `input.parameters` is specifying the `input` that will be passed to the `plugin` and the `plugin` will use it to generate the `output`.

- The `labels` will be provided by `Pull Request Generator`. It will be a list of labels that are assigned to the `PR`.

- The `excludeLabel` is the label that if it exists in `PR` we don't need a preview environment and use `blackhole`.

- The `path` is the real path that we have our service manifests in it.

- The `blackHole` is the path that we want to use if we don't want a preview environment for a `PR`. And it's empty directory.

- The `number` is the `PR` number that we want to generate `Application` for it.

Our logic for `plugin`:

1. Check all `labels` in `PR` and if `excludeLabel` exists in `PR`, set `generatedPath` to `blackHole` path.
2. If `excludeLabel` doesn't exist in `PR`, set `generatedPath` to `path` path which is the real path that we have our service manifests in it.
3. Output the `generatedPath` as `output` of `plugin`, which will be available in `ApplicationSet` template.


So adding `no-preview` label to `PRs` will generate `Application` for that `PR` but it will point to `blackhole` directory which has no `K8S manifest` in it, which means we don't deploy any `Kubernetes` resources for that `PR`.

```yaml
  template:
    metadata:
      name: "br-{{.branch}}-{{.number}}"
    spec:
      source:
        repoURL: "https://github.com/araminian/BRApplication.git"
        targetRevision: HEAD
        path: "{{.generatedPath}}"
        directory:
          include: '{*.yml,*.yaml}'
      project: "default"
      destination:
        server: https://kubernetes.default.svc
        namespace: "default"
      syncPolicy:
        syncOptions:
          - CreateNamespace=true
```
the `{{.generatedPath}}` will be replaced with the `output` of `plugin` which is the `path` that we want to use for generating `Application`. It can be either `blackHole` or `path`.


## How to use Plugin

1. Install new version of `argoCD`. `Plugin Generator` feature is not available in `stable` version yet!!

```bash
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/master/manifests/core-install.yaml
```

2. Install `plugin` and it's cofigurations:

for more information about `Plugin Generator` check [here](https://github.com/binboum/argo-cd/blob/1f7ff874c4961edd7e50595602c174936619299c/docs/operator-manual/applicationset/Generators-Plugin.md)

```bash
kubectl apply -f ./manifests
```

3. Create `ApplicationSet` with `plugin` generator:

```yaml
          - plugin:
              configMapRef:
                name: applicationset-pr-filter-plugin
              input:
                parameters:
                  labels: "{{.labels}}"
                  excludeLabel: "no-preview"
                  path: "manifests/app/"
                  blackHole: "manifests/blackhole/"
                  number: "{{.number}}"
```