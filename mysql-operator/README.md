# A few notes about the operator

## About testing

Kubebuilder comes with `envtest`, an easy way to run integration tests, without
having a Kubernetes cluster. To install the testing environment, run the
following commands:

```shell
curl -sSLo setup_envtest.sh \
  https://raw.githubusercontent.com/kubernetes-sigs/kubebuilder/master/scripts/setup_envtest_bins.sh 
chmod +x setup_envtest.sh
./setup_envtest.sh v1.18.6 v3.4.10
```

You should also modify the `Makefile` as described in
[envtest setup](https://master.sdk.operatorframework.io/docs/building-operators/golang/references/envtest-setup/)
in the [operator-sdk documentation](https://master.sdk.operatorframework.io/docs/).
This change has already been done.

