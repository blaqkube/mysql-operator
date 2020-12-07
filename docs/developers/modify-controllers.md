# Modify Controllers

Controllers are programs that manages resources created via the API. 

## Create a new Controller

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
operator-sdk create api --version v1alpha1 --group mysql --kind Backup
Create Resource [y/n]
n
Create Controller [y/n]
y
```
