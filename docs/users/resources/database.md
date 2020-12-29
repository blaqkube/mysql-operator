# Database

Databases are created in the MySQL Instance. This is an example of a manifest:

```yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Database
metadata:
  name: red
spec:
  instance: blue
  name: red
```

The properties are the following:

- `name` defines the database name
_ `instance` defines the instance the database is created with
