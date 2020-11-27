# Modify APIs

The mysql-operator API is defines by Custom Resource Definitions or CRDs.

## Create a new CRD

```shell
operator-sdk add api --api-version=mysql.blaqkube.io/v1alpha1 --kind=Store
```

## Modify the structure from the resources

```shell
cd pkg/apis/mysql/v1alpha1
vi backup_types.go
```

## Update the generated code

```shell
operator-sdk generate k8s
```

## Update CRDs

```shell
operator-sdk generate crds
```

## Change the CR example

```shell
vi deploy/crds/mysql.blaqkube.io_v1alpha1_backup_cr.yaml
operator-sdk generate crds
```

