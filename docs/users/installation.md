# Installation

To install the MySQL operator, you must perform a manual installation.
The previous installation methods have been removed due to the project
not being maintained anymore

## Manual Installation

To install OLM manually, you need to have `make`, `kustomize`, `kubectl`, `git`
and `go` installed. We assume the controller has been built and is available
from `quay.io/blaqkube/mysql-controller`. If you are using another registry and
versioning scheme, you would have to change the `IMG` value accordingly.

To deploy the operator, run:

```shell
git clone https://github.com/blaqkube/mysql-operator.git
cd mysql-operator/mysql-operator
# checkout the version of your choice
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
