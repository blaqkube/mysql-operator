6. Regenerate the CSV and check the CR has been integrated

operator-sdk generate csv --csv-version=0.0.3 --update-crds    
vi deploy/olm-catalog/mysql-operator/manifests/mysql-operator.clusterserviceversion.yaml

7. Build the controller and move it to quay.io

see `make build`

8. Modify the version to match the --csv-version and create the bundle

see `make bundle`

9. Create a docker image with the catalog

cd hack
make

11. Check the version of the is available

kubectl get sub -n mysql -o yaml
The subscription should show the version is the last one


14. Vérifier la configuration du plan

Il semble que le problème vienne de   clusterServiceVersionNames: [] comme indiqué dans https://github.com/operator-framework/operator-lifecycle-manager/issues/1347

