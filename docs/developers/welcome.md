---
title: Welcome Contributors
---

Thank you for considering contributing to `blaqkube/mysql-operator`! This
section is written to help you. 

## Overview

A Kubernetes operator is a set of custom resource definitions, the API and
controllers packaged together to provide a service for an application. This
MySQL operator fit the definition and provide you the ability to install,
backup and restore MySQL databases.

The project contains a number of components that are used to manage the
MySQL instances.

- `agent` contains a MySQL agent that are installed within the MySQL
  StatefulSet as a sidecar and is used by the controller to perform
  database commands like backups
- `mysql-operator` contains both the APIs in `pkg/apis` and
  the controller in `pkg/controller`.
- `registry` contains the registry for the application
- `docs` contains the documentation

## Tools

In order to develop, you would need a number of tools, including `go`,
`operator-sdk`, `kubectl` and a Kubernetes cluster. What follows provide
some hints about those.

### Go

The operator is written in Go, you can get it from the
[downloads](https://golang.org/dl/) page. Install it as well as some
tools to help with it, like a code editor and git.

### operator-sdk

This operator relies on
[operator-framework/operator-sdk](https://github.com/operator-framework/operator-sdk).
You will very likely need it. Installing `operator-sdk` is usually as simple as:

- downloading the executable from the
  [release page](https://github.com/operator-framework/operator-sdk/releases);
- making sure you can execute the file
- adding the file to your `PATH`

To use `operator-sdk`, you might need to set `GOROOT` to workaround
[Issue #1480](https://github.com/operator-framework/operator-sdk/issues/1480)
like below:

```shell
export GOROOT=$(go env GOROOT)
```

If you do not set GOROOT; you might 

```text
F0524 12:39:30.302991  168940 deepcopy.go:885] Hit an unsupported type invalid
type for invalid type
```

> Note: To know more about operator-sdk and how to install it, see
> [Operator SDK documentation](https://sdk.operatorframework.io/docs)

### Kind, Docker and Minikube

To test the operator, you would need a Kubernetes cluster. An easy way to
get one with multiple nodes is to rely on [Kind](https://kind.sigs.k8s.io/).
You can also rely on
[Docker Desktop](https://docs.docker.com/docker-for-mac/kubernetes/)
or [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/).

