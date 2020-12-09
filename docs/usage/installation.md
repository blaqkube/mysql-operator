# Installation

To install the MySQL operator from blaqkube, you can perform a manual
installation from the repository or you can rely on the Operator Lifecycle
Manager (OLM). The later option is the preferred way to proceed as it will not
only proceed with the installation but will also help with upgrading the tool
later on.

## Installation with OLM

To install `blaqkube/mysql-operator` with OLM, you should have it installed on
your cluster. Once done, you should be able to connect to the blaqkube registry
and subscribe to the Operator.

### Installing OLM

The simplest way to install OLM on your cluster is:

- to [Install the Operator SDK CLI](https://sdk.operatorframework.io/docs/install-operator-sdk/)
- Make sure your can access a Kubernetes cluster with cluster admin role
- Run `operator-sdk olm` like below

```shell
operator-sdk olm install
operator-sdk olm status
```

> Note: In order for `blaqkube/mysql-operator` to update correctly, you should
> use OLM 0.15+. This is due to issue
> [#1347](https://github.com/operator-framework/operator-lifecycle-manager/issues/1347)

### Connecting to the registry

In the following section, we will assume we work in the default namespace. If you want
to use an other namespace, change the configuration accordingly.

Once OLM installed, create a `CatalogSource` that references the blaqkube registry. The
registry contains the references to the operator. To create the `CatalogSource`, run
the command below:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: blaqkube-catalog
spec:
  sourceType: grpc
  displayName: blaqkube Operators
  address: registry.blaqkube.io:50051
  publisher: blaqkube.io
EOF
```

### Subscribing to the operator

You are almost ready to install the Operator. Before you proceed, you should
create a `OperatorGroup` that specify which namespaces the Operator can
manage. Here again, we assume the targeted namespace is the default one and
you can change it to fit your needs:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: mysql-operatorgroup
spec:
  targetNamespaces:
  - default
EOF
```

To subscribe to the operator, create a `Subscription` with `mysql-operator`
and the `CatalogSource` created previously. The `OperatorGroup` does not
have to be referenced, there should only be one per namespace:

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

Manual installation consists in cloning the git repository first; install the
Custom Resource Definitions, grant the required permission and starting the
controller.

### Cloning the project

Use `git` to clone the project like below:

```shell
git clone https://github.com/blaqkube/mysql-operator.git
cd mysql-operator
```

### Installing CRDs

Once done, you should install the CRDs like below:

```shell
cd $(git rev-parse --show-toplevel)/deploy/crds
kubectl apply -f mysql.blaqkube.io_backups_crd.yaml
kubectl apply -f mysql.blaqkube.io_stores_crd.yaml
kubectl apply -f mysql.blaqkube.io_instances_crd.yaml
```

### Permissions and Controller

You should then be able to create a role with the right permissions and deploy
the operator:

```shell
cd $(git rev-parse --show-toplevel)/deploy/crds
kubectl apply -f role.yaml
kubectl apply -f role_binding.yaml
kubectl apply -f service_account.yaml
kubectl apply -f operator.yaml
```

Once done with the installation, you can create the resources as described in
[next section](resources) of the documentation
