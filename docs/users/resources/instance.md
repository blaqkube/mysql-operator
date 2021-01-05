# Instance

Instances are used to create a stateful set with the `mysql:8.0.20` container
as well as the associated sidecars. This is an example of an Instance manifest:

```yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Instance
metadata:
  name: blue
spec:
  database: blue
  restore:
    store: docs
    location: "/location/backup-20200526110950.sql"
```

The properties are the following:

- `database` calls a database to create at startup
_ `restore` is used to define the 2 mandatory parameter to base the instance
  on a previous backup:
  - `store` name the Store the backup is located in
  - `filePath` defines the key for the file. It should start with a `/`
- `maintenance` is used to define maintenance parameters for the instance:
