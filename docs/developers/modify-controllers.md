# Modify Controllers

Controllers are programs that manages resources created via the API. 

## Create a new Controller

To create a new controller, you should use the `operator-sdk create api`
command like below:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
operator-sdk create api --version v1alpha1 --group mysql --kind Backup
Create Resource [y/n]
n
Create Controller [y/n]
y
```

Controllers are created inside the `mysql-operator/controllers` directory.
You then should be able to write the logic that the controller implements.
