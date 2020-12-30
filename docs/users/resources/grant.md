# Grant

Grant grants access on a database to a user. Below is an example of of the
manifest:

```yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Grant
metadata:
  name: instance-user-database
spec:
  user: instance-user
  database: instance-database
  accessMode: readWrite
```

The properties are the following:

- `user` defines the resource that references the user
- `database` defines the resource that references the database
- `accessMode` should be set to `readOnly` or `readWrite` 
