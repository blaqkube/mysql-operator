# Installation

To install the MySQL operator, you should rely on the Operator Lifecycle
Manager (OLM) and, at least for now, a `CatalogSource` that references
the blaqkube registry. OLM provides an easy way to discover operators and
helps with release management.

## Operator-sdk Lifecycle Manager (OLM)

The simplest way to install OLM on your cluster is to follow the
[Install the Operator SDK CLI](https://sdk.operatorframework.io/docs/installation/)
section of the documentation. Once done, run `operator-sdk olm` like below:

```shell
operator-sdk olm install
operator-sdk olm status
```

> Note: In order for `blaqkube/mysql-operator` to update correctly, you should
> use OLM 0.15+. This is due to issue
> [#1347](https://github.com/operator-framework/operator-lifecycle-manager/issues/1347)

## Declare the CatalogSource

Once OLM installed, create a `CatalogSource` that references the blaqkube
registry. The registry contains references to the operator. To create the
`CatalogSource`, run the command below:

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
cat <<EOF | kubectl delete -f -
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

It should create a statefulset with your instance. To delete the instance, run
`kubectl delete instance red`.