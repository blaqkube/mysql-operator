# Testing the operator locally

Testing the operator is key to develop successfully. To do so, you need to
access a Kubernetes cluster. Using Kind is usually a good way to work. Once
you have access to the cluster, you can:

## Install the CRDs

```shell
cd deploy/crds
kubectl apply -f mysql.blaqkube.io_backups_crd.yaml
kubectl apply -f mysql.blaqkube.io_stores_crd.yaml
kubectl apply -f mysql.blaqkube.io_instances_crd.yaml
```

## Run the controllers

```shell
operator-sdk run --local
```

