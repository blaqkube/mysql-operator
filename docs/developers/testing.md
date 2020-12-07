# Testing the operator locally

Testing the operator is key to develop successfully. You can run the operator
outside of kubernetes the procedure provided as part of the
[welcome](./welcome.md) section of this documentation.

There is actually more... The project comes with `envtest` tests. `envtest` is
a controlplane/etc simulator that can emulate Kubernetes and is built as part
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
