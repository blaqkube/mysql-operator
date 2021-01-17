# Instance

Instances are used to create a stateful set with the `mysql:8.0.22` container
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
- `restore` is used to define the 2 mandatory parameter to base the instance
  on a previous backup:
  - `store` names the Store the backup is located in
  - `location` defines the key for the file. It should start with a `/`
- `backupSchedule` is used to define automatic backups. It should include 2
  parameters:
  - `store` names the store the backup are stored in
  - `schedule` defines the cron-like scheduled expression. Pay attention to the fact the timezone is UTC
  