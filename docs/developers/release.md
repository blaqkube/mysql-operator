# Releasing the operator

There are 4 components to release as part of the operator:

- The `mysql-agent` is located in the agent directory. It is started as a
  sidecar of the database StatefulSet and provide advanced features like backup
  and restore. It is released as a container image
- `Controllers` are logic that interact with Kubernetes from changes in the
  configuration. It is released as a container image
- The `Operator` reference the controller and includes the manifests including
  CRDs, permissions and the schedule of `Controllers`. It is also packaged and
  released as a container image
- The `Registry` contains references to all the operator version. It is 
  also built as a container image and deployed on the blaqkube infrastructure.

> Note: Releasing the operator requires access to the different components
> and the change to be correclty reviewed and tested. As a result only
> project owners can trigger a new release.

## Building the agent

The agent is built as part the CI. There is nothing to do to make it happen.

## Releasing the agent

Releasing the agent is actually part of the operator built process. As a matter
of fact, the agent is referenced in the `mysql-operator/main.go` file. To
proceed:

- Change the agent client part in the `mysql-operator/agent` directory
- Change the agent version in the `mysql-operator/main.go`.

There is a Makefile script that does the change for you. To run that script,

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator/agent
make api
```

## Building Controllers

The controller is built as part the CI. Just pay attention to the fact that,
unless you have explicitly release the agent, it will be built with the last
release agent version.

## Releasing Controllers

Controllers are released with the operator.

## Building and Releasing the Operator

Building and releasing the operator is part of the same process that also
embeds the releasing of Controllers. To proceed, you should plan it.

### Preparation

- Update VERSION in .gally.yml
- CHANGELOG
- RUN ALL the tests
- Make sure the agent latest version is used
- Merge to Master

### Releasing

Tag v${VERSION} on main


## Build and release the Registry

see `make registry` in the `registry` directory. This commands update the
index.db file with the new operator. It creates a docker image to serve the
different versions. The publication of the registry is done outside of this
repository.

## Check the Operator get updated

Once the new registry published, you can check your subscription and control
the Operator has been updated. The subscription should show the version is
the last one:

```shell
kubectl get sub -n mysql -o yaml
```

> Note: there is a bug [1347](https://github.com/operator-framework/operator-lifecycle-manager/issues/1347)
> with OLM previous than 0.15.0 that prevents `clusterServiceVersionNames` to
> be set correctly and prevent the upgrade of the operator. Make sure you use
> OLM 0.15.0+
