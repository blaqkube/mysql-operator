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
export http_proxy=http://localhost:3128
make run ENABLE_WEBHOOKS=false
```

## Running the Operator in a Container 

Assuming an OCI image for the controller has been created in the CI, you can test it
instead of running the manager outside the cluster. To proceed:

- Clone the project with `git clone https://github.com/blaqkube/mysql-operator`
- Go into the operator subdirectory `cd mysql-operator/mysql-operator`
- Install the CRDs to your default namespace `make install`
- Set `IMG` with the version of the controller OCI that is expected
- Run controllers in your cluster `make deploy`

The script below does the steps above:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
export IMG=quay.io/blaqkube/mysql-controller:$(\
    git log --format='%H' -1 . | cut -c1-16)
echo $IMG
make deploy
```

## kuttl for testing

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

`kuttl` provides different tests. To run those tests, you should have your
operator up and running:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make install
kubectl apply -f .ci/squid.yaml
kubectl port-forward squid 3128 &
export http_proxy=http://localhost:3128
make run ENABLE_WEBHOOKS=false
```

Running `kuttl` tests should be as simple as 

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator/tests
kubectl kuttl test
```

## Integration tests with kuttl

We also provide some integration tests that you can also run with `kuttl`. In
order for those tests to work, you should set a secret with your AWS/GCP
credentials. The associated secrets are in 
`mysql-operator/tests/integration/store/secrets.yaml`. Change them to meet your
requirements.

You can simply apply those:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator/tests
kybectl apply -f integration/store/secrets.yaml
```

Then, you need to change the bucket from the kuttl manifests. You can do it
with a GNU `sed` command like below. Update `MYBUCKET` variable with the bucket
of your choice:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator/tests
export MYBUCKET=yours3bucket
sed -i s/bucket\.blaqkube\.io/$MYBUCKET/ integration/store/s3/00-install.yaml
export MYBUCKET=yourgcpbucket
sed -i s/bucket\.blaqkube\.io/$MYBUCKET/ integration/store/gcp/00-install.yaml
```

You can now start the operator with your cluster:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make install
kubectl apply -f .ci/squid.yaml
kubectl port-forward squid 3128 &
export http_proxy=http://localhost:3128
make run ENABLE_WEBHOOKS=false
```

Running the test should then be as easy as running:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator/tests
kubectl kuttl test integration/store --test s3
```

Testing is key to the project. There is a lot to improve so do not hesitate to
provide some feedback.