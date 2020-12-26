# User

User are created with a MySQL Instance. This is an example of a manifest:

```yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: User
metadata:
  name: myuser
spec:
  instance: blue
  username: myuser
  password: changeme
```

The properties are the following:

_ `instance` defines the instance the user is created in
- `username` defines the user name
- `password` set the user password
