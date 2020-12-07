# Modify APIs

The mysql-operator API is defines by Custom Resource Definitions or CRDs.

## Create a new API

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
operator-sdk create api --version v1alpha1 --group mysql --kind Backup
Create Resource [y/n]
y
Create Controller [y/n]
n
```

## Modify the structure of the resource

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
cd api/v1alpha1/backup_types.go
```

## Update the generated code

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make generate
```

## Update CRDs

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make manifests
```

## Change the CR example

```shell
vi config/samples/mysql_v1alpha1_backup.yaml
```
