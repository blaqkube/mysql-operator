# Welcome

Thank you for using and considering contributing to the
`blaqkube/mysql-operator`! This section is written to help developers
starting with the project. Leave an issue on
[Github Project](https://github.com/blaqkube/mysql-operator/issues) if needed.

## Overview

Kubernetes operators are custom resource definitions and controllers packaged
together to manage an application.
[blaqkube/mysql-operator](https://github.com/blaqkube/mysql-operator) can be
used to install backup and restore MySQL databases.

The project contains a number of components that are used to manage the
MySQL instances.

- `docker-gally` contains an artifact used on CircleCI to build and deploy
  the right component
- `agent` contains a MySQL agent that are installed within the MySQL
  StatefulSet as a sidecar and is used by the controller to perform
  database commands like backups
- `mysql-operator` contains both the APIs in `api` and the controller in the
  `controller` subdirectory.
- `registry` contains the registry for the application
- `docs` contains the documentation

## Development environment

You need a number of tools to develop. They include `go`, `operator-sdk`,
`kubectl`, `gcc`, `make` and a Kubernetes cluster, like `kind` or `minikube`.
We will assume you have setup and configured all those tools so that you can
run `kubectl` and you can manage the cluster.

## Running the operator manually

The operator relies on the
[Golang version of operator-sdk](https://sdk.operatorframework.io/docs/building-operators/golang/).
To run it from outside of the cluster:

- Clone the project with `git clone https://github.com/blaqkube/mysql-operator`
- Go into the operator subdirectory `cd mysql-operator/mysql-operator`
- Install the CRDs to your default namespace `make install`
- Run controllers outside of your cluster `make run ENABLE_WEBHOOKS=false`

```shell
cat <<EOF | kubectl apply -f -
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Instance
metadata:
  name: blue
spec:
  database: blue
```
