# Backup

Backups are used to trigger backups. Below is an example of a Backup manifest:

```yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Backup
metadata:
  name: blue-backup
spec:
  store: docs
  instance: blue
```

The properties are the following:

- `store` defines the store used to perform the backup
- `instance` defines the instance to backup.

