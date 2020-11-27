# Releasing the operator

There are 4 components to release as part of the Operator:

- The mysql-agent is located in the agent directory. It is started as a sidecar
  of the database StatefulSet and provide advanced features like backup and
  restore. It is released as a container image
- `Controllers` are logic that interact with Kubernetes from changes in the
  configuration. It is released as a container image
- The `Operator` reference the controller and includes the manifests including
  CRDs, permissions and the schedule of `Controllers`. It is also packaged and
  released as a container image
- The `Registry` contains references to all the operator version. It is 
  also built as a container image and deployed on the blaqkube infrastructure.

> Note: Releasing the operator requires access to the different components
> and the change to be correclty reviewed and tested. As a result only
> project owners can trigger a new release.

## Build and release the agent

To build the agent, you should execute `make build` in the `agent` directory.
Once the agent is built, you must change the reference in the
`operator/pkg/controller/instance/instance_controller.go` file. Search for
the `tag := "` in the file and replace the tag with the new one.

## Build and release Controllers

see `make controller` in the `mysql-operator` directory. This command builds 
the image with the controller and pushes to quay.io. Once done, it also change
the `operator.yaml` file with the new controller version.

## Build and release the Operator

see `make bundle` in the `mysql-operator` directory. This command updates
the CRDS and build an image of the operator that references the controller. It
also push the image to `quay.io`

## Build and release the Registry

see `make registry` in the `registry` directory. This commands update the
index.db file with the new operator. It creates a docker image to serve the
different versions. The publication of the registry is done outside of this
repository.

## Check the Operator get updated

Once the new registry published, you can check your subscription and control
the Operator has been updated. The subscription should show the version is
the last one:

```shell
kubectl get sub -n mysql -o yaml
```

> Note: there is a bug [1347](https://github.com/operator-framework/operator-lifecycle-manager/issues/1347)
> with OLM previous than 0.15.0 that prevents `clusterServiceVersionNames` to
> be set correctly and prevent the upgrade of the operator. Make sure you use
> OLM 0.15.0+

