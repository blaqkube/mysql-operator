# Testing the operator locally

Testing the operator is key to develop successfully. You can run the operator
outside of kubernetes with the Envtest Kubernetes simulator. You can also run
it manually on top of a development cluster like `kind`. This document details
some of the various options provided with the project.

## Using Envtest

There is actually more... The project comes with `envtest` tests. `envtest` is
a controlplane/etcd simulator that can emulate Kubernetes and is built as part
of the controller-runtimer project. To know more about `envtest`, read the
[envtest documentation](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest).

The procedure to install `envtest` is part of the operator `Makefile`. To those
tests, execute the following command:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make test
```

The procedure creates a directory named `testbin` in `mysql-operator`; this
directory includes a set of binaries to simulate kubernetes controlplane.

## Running the operator manually

The operator relies on the
[Golang version of operator-sdk](https://sdk.operatorframework.io/docs/building-operators/golang/).
To run it from outside of the cluster, for development purpose:

- Clone the project with `git clone https://github.com/blaqkube/mysql-operator`
- Go into the operator subdirectory `cd mysql-operator/mysql-operator`
- Install the CRDs to your default namespace `make install`
- Install a proxy server so that controllers can access the content of the
  cluster with REST calls
- Run controllers outside of your cluster `make run ENABLE_WEBHOOKS=false`

The script below does the steps above:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make install
kubectl apply -f .ci/squid.yaml
kubectl port-forward squid 3128 &
export HTTP_PROXY=http://localhost:3128
make run ENABLE_WEBHOOKS=false
```

## Using kuttl

The project embeds kuttl tests. In order to run those tests, there are
a few things to consider.

First thing first, you need to install kuttl. On MacOS, it should be as
simple as running `brew install kuttl`. On Linux, you can simply add the
`kubectl-kuttl` binary to your `/usr/local/bin` directory like below:

```shell
export KUTTL=0.7.2
DOWNLOAD=https://github.com/kudobuilder/kuttl/releases/download
sudo su -
cd /usr/local/bin
curl -Lo kubectl-kuttl \
  $DOWNLOAD/v${KUTTL}/kubectl-kuttl_${KUTTL}_linux_x86_64
```

`kuttl` provides integration tests. The problem with integration tests is that
you require consistent data. To proceed, create a bucket. Then create some S3 credentials that you should store in the `store-sample` secret:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: store-sample
type: Opaque
stringData:
  AWS_ACCESS_KEY_ID: AKIA...
  AWS_SECRET_ACCESS_KEY: secret...
  AWS_REGION: us-east-1
EOF
```

Then, you need to change the bucket from the kuttl manifests. You can do it
with a GNU `sed` command like below. Update `MYBUCKET` variable with the bucket
of your choice:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
export MYBUCKET=logs.blaqkube.io
sed -i s/logs.blaqkube.io/$MYBUCKET/ integration/store/00-install.yaml
```

You can now start the operator with your cluster:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make install
kubectl apply -f .ci/squid.yaml
kubectl port-forward squid 3128 &
export HTTP_PROXY=http://localhost:3128
make run ENABLE_WEBHOOKS=false
```

Running the test should then be as easy as running:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
kubectl kuttl test
```
