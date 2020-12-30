# Releasing the Operator

There are 4 components to release as part of the operator:

- The `mysql-agent` is located in the `agent` directory. It is started as a
  sidecar of the instance StatefulSet and provides advanced features like backup
  and restore. It is released as a container image
- `Controllers` are pieces of logic that interact with Kubernetes when the
  configuration changes. It is also released as a container image
- The `Operator` glues things together. It provides the manifests, including
  CRDs, roles, deployment for `Controllers` and more. It is also packaged and
  released as an OCI image
- In addition to the operator, the project also provides an `Index` that
  references the operator versions and ease deployments and updates.

In adidtion to these components, the project also includes:

- An OCI image named `docker-gally` that is used in CircleCI to build and
  release the other components
- A set of documentation

## Building the Agent

The agent is built as part the CI. There is nothing to do to make it happen.
Every change to the `agent` directory generates a version that is tagged with
a part of the Git commit on the directory.

## Releasing the Agent

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

> Note: the client generator requires node and 
> [@openapitools/openapi-generator-cli](https://www.npmjs.com/package/@openapitools/openapi-generator-cli)
> to be generated.

## Building Controllers

The controller is built as part the CI. Just pay attention to the fact that,
unless you have explicitly release the agent, it will use a previous version of
the agent.

## Releasing Controllers

Controllers are released with the Operator.

## Releasing the Operator

Building and releasing the operator is actually the same process. Releasing the
operator will also release controllers. To proceed, you should

- prepare the release
- create a tag on the repository

### Preparation

Preparing for a release of an operator implies the tasks below that are done
as part of pull requests on the repository:

- Update the `VERSION` and `PREV_VERSION` variables in the `.gally.yml` file of
  the `mysql-operator` subdirectory. These variables are used to control the
  versioning. Set `PREV_VERSION` to the previous version and `VERSION` to the
  version yet to come.
- Rebuild the bundle files with the `make bundle` command in the `mysql-operator`
  directory.
- Update the `CHANGELOG.md` with the list of changed in the release
- Run all the tests that are part of the version. In particular, all the
  `mysql-operator/tests` integration and end-to-end tests that are built on
  `kuttl` and not yet part of the CI
- Make sure the agent latest version is used. The operator will not be released
  if it does not include the latest agent release. Make sure you run `make api`
  from `mysql-operator/agent`. This should regenerate the client API for the
  agent and add the tag for the AGENT in the API

Once all the steps have been performed and merged on the `main` branch, you are
ready to perform the release

### Releasing

The release actually consists in adding a tag `v${VERSION}` on the `main`
branch to Github. CircleCI does the rest.

> Note: the versioning need to comply to [semver](https://semver.org/)

## Documentation

The documentation website is part of a separate repository, even if all the
content is actually from `docs`. There is a separate process.

## Build and Release the Index

The `index` is available as an OCI image. It is built with the operator. It
is released as part of a separate process allowing to perform additional tests. 

## Check the Operator

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
