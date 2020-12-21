# APIs

APIs are made of Custom Resource Definitions or CRDs. operator sdk relies on
kubebuilder to scaffold the directory structures. Generators are provided to
build derivative constructs like DeepCopy methods and the YAML manifests.

The way it works is straightforward:

- First, create the API skeleton in the `api/<groupversion>` directory for
  the operator. API are made of Go `struct`s that you can extend to adapt to
  your needs
- Second, generate the Copy methods that are required on your API for the
  controller to be implemented. If you wonder why this needs to be done, the
  short answer is because, there is a need to move parts by copying them AND
  code generation is the Go alternative to generics.
- Third, create the YAML manifest files for your APIs so that you can install
  them with Kubernetes and be handler with the operator lifecycle manager. You
  can adapt those files with examples and some metadata properties.

The sections below provide some directions to help with these 3 tasks.

## Create and modify the API structure

API are made of `struct` in the `api` subdirectories. Most of the work should
be modifying those ; however, you might want to create a new API. If that is
the case, you can run the `operator-sdk create api` command like the one below
that was executed to create the `Backup` API:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
operator-sdk create api --version v1alpha1 --group mysql --kind Backup
Create Resource [y/n]
y
Create Controller [y/n]
n
```

Once done, modify the structure of the resource with what you intend to make
it do:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
vi api/v1alpha1/backup_types.go
```

## Generate the Copy methods

Generate the code that is necessary for the API should be as simple as the
command below. It has to be done if you modify the `struct` associated with
the API:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make generate
```

## Update CRDs and improve them

To bundle and deliver the API, manifests should be generated. Note that part
of those are regenerated everytime, like the CRD yaml files and some are
not like custom resources samples. To regenerate the manifests, you should
run:

```shell
cd $(git rev-parse --show-toplevel)
cd mysql-operator
make manifests
```

The you can edit some of these to improve them like for the sample resources:

```shell
vi config/samples/mysql_v1alpha1_backup.yaml
```

## More...

If you want to know more about API generation and the integration between
`kubebuilder` and `operator-sdk`, have a look at the
[kubebuilder designs](https://github.com/kubernetes-sigs/kubebuilder/tree/master/designs). Other than that:

- operator-sdk [Create a new API and Controller](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/#create-a-new-api-and-controller)
- kubebuilder [Adding a new API](https://book.kubebuilder.io/cronjob-tutorial/new-api.html)
