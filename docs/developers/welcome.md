# Welcome

Thank you for contributing to the `blaqkube/mysql-operator`! This section is
written to help developers starting with the project. Do not hesitate to leave
an issue on the [Github](https://github.com/blaqkube/mysql-operator/issues)
if needed.

## Overview

Kubernetes operators are custom resource definitions and controllers packaged
together to manage an application.
[blaqkube/mysql-operator](https://github.com/blaqkube/mysql-operator) can be
used to install backup and restore MySQL databases.

The project contains a number of components that are used to manage the
MySQL instances.

- `agent` contains a MySQL agent that are installed within the MySQL
  StatefulSet as a sidecar and is used by the controller to perform
  database commands like backups
- `docker-gally` is an artifact used on CircleCI to build and deploy
  the right component. It is based on
  [missena-corp/gally](https://github.com/missena-corp/gally)
- `docs` contains the project documentation that is available online
  one [docs.blaqkube.io](https://docs.blaqkube.io/)
- `mysql-operator` contains everything related to the operator; this
  includes controllers, API, tests, artifact builds and more.
- `registry` contains a registry for the application

## Development environment

You need a number of tools to develop. They include `go`, `operator-sdk`,
`kubectl`, `gcc`, `make`, `bash` 4+, GNU `sed` and a Kubernetes cluster, like
`kind` or `minikube`. We will assume you have setup and configured all those
tools so that you can run `kubectl` and you can manage the cluster.

> **Note**: If you use a Mac, pay attention to install bash 4+ and GNU sed and
> make sure they are first in your PATH.

Controllers must access the MySQL Pod to perform some maintenance tasks like
create a database. In order to do it, it communicates with an instance sidecar
running the operator MySQL agent. It does that via OpenAPI. To run them
outside the cluster, you should install a proxy server and establish a
local connection to it. The `.ci/squid.yaml` file create a simple pod manifest
to grant access to the cluster network. Before you run the controllers, run
the command below:

```shell
kubectl apply -f .ci/squid.yaml
kubectl port-forward squid 3128
export http_proxy=http://localhost:3128
```

## Running the operator manually

The operator relies on the
[Golang version of operator-sdk](https://sdk.operatorframework.io/docs/building-operators/golang/).
To run it from outside of the cluster, for development purpose:

- Clone the project with `git clone https://github.com/blaqkube/mysql-operator`
- Go into the operator subdirectory `cd mysql-operator/mysql-operator`
- Install the CRDs to your default namespace `make install`
- Make sure you have installed an HTTP proxy as described in the previous
  section
- Run controllers outside of your cluster `make run ENABLE_WEBHOOKS=false`

The operator should start. Once done, you can create a MySQL instance with a
manifest like the one below; run `kubectl apply -f blue-instance.yaml` to
install it:

```yaml
# file blue-instance.yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Instance
metadata:
  name: blue
spec:
  database: blue
```

A statefulset named `blue` should be created as seen below:

```shell
kubectl get sts blue
kubectl get pod blue-0
```

## Clean the configuration

To clean the environment, stop the operator with a `Ctrl+C`. You can remove the instance with `kubectl delete instance blue`. It will remove the statefulset, replicaset and pod.
