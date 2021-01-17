# Installation

To install the MySQL operator, you can:

- Rely on the Operator Lifecycle Manager (OLM) and operatorHub.io
- Rely on the Operator Lifecycle Manager (OLM) and a `CatalogSource` that
  references the blaqkube registry.
- Perform a manual installation
  
OLM provides an easy way to discover operators and helps with release
management. This is the preferred way to install the MySQL Operator.

## Install Operator Lifecycle Manager (OLM)

If you plan to install the operator from the operatorHub.io or the blaqkube
registry, you need to have OLM installed on your cluster. The simplest way
proceed is to follow the
[Install the Operator SDK CLI](https://sdk.operatorframework.io/docs/installation/)
section of the documentation. Once done, run `operator-sdk olm` like below:

```shell
operator-sdk olm install
operator-sdk olm status
```

> Note: In order for `blaqkube/mysql-operator` to update correctly, you should
> use OLM 0.15+. This is due to issue
> [#1347](https://github.com/operator-framework/operator-lifecycle-manager/issues/1347)

## Subscribe with OperatorHub.io

To install the operator, you need to create an `OperatorGroup` to declare the
scope of the operator. For now, the only supported scope is clusterwide. Create
an `OperatorGroup` like below :

```shell
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: mysql-operatorgroup
EOF
```

Create a `Subscription` with `mysql-operator` and the OperatorHub Catalog. The
`OperatorGroup` does not have to be referenced, there should only be one per
namespace:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: mysql-operator
spec:
  channel: alpha
  name: mysql-operator
  source: operatorhubio-catalog
  sourceNamespace: olm
EOF
```

> Note: for now on, only the `alpha` channel is available. It is frequently
> updated but might contain imcompatibles changes between releases.

> Note: this subscription assumes the operator is upgraded automatically.
> to perform manual upgrade, use the `installPlanApproval` and `startingCSV`
> properties.

Once done with the installation, you can create the resources as described in
[next section](resources) of the documentation

## Subscribe with the Blaqkube Registry

Installing the operator from the Blaqkube registry is very similar to relying
on the OperatorHub.io catalog except for the `CatalogSource` that you must
create. The sections below detail the required steps.

### Declare the CatalogSource

You must create a `CatalogSource` that references the blaqkube registry. The
registry contains references to the operator. To create the `CatalogSource`,
run the command below:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: blaqkube-catalog
spec:
  sourceType: grpc
  displayName: blaqkube
  address: registry.blaqkube.io:50051
  publisher: blaqkube.io
EOF
```

### Subscribe to the operator

Before you proceed, with installing the operator, you should create an
`OperatorGroup` that defines the fact the operator works at the cluster
level as it is the only supported scope for now:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: mysql-operatorgroup
EOF
```

Create a `Subscription` with `mysql-operator` and the `CatalogSource` created
previously. The `OperatorGroup` does not have to be referenced, there should
only be one per namespace:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: mysql-operator
spec:
  channel: alpha
  name: mysql-operator
  source: blaqkube-catalog
  sourceNamespace: default
EOF
```

> Note: for now on, only the `alpha` channel is available. It is frequently
> updated but might contain imcompatibles changes between releases.

> Note: this subscription assumes the operator is upgraded automatically.
> to perform manual upgrade, use the `installPlanApproval` and `startingCSV`
> properties.

Once done with the installation, you can create the resources as described in
[next section](resources) of the documentation

## Manual Installation

To install OLM manually, you need to have `make`, `kustomize`, `kubectl` and `git` installed. We assume the controller has been
built and is available from `quay.io/blaqkube/mysql-controller`.

To deploy the operator, run:

```shell
git clone https://github.com/blaqkube/mysql-operator.git
git checkout main
cd mysql-operator/mysql-operator
make install
export IMG=quay.io/blaqkube/mysql-controller:$(\
    git log --format='%H' -1 . | cut -c1-16)
echo $IMG
make deploy
```

## Verify the configuration

An easy way to verify the configuration is to create an instance with the
command below:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Instance
metadata:
  name: red
spec:
  database: red
EOF
```

If you need to create a database to the instance, apply another manifest like
the one below:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Database
metadata:
  name: red-blue
spec:
  name: blue
  instance: red
EOF
```

It should create a statefulset with your instance and add a database to it. To
delete the instance and database, run:

```shell
kubectl delete instance red
kubectl delete database red-blue
```
