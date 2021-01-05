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
- A password that can be made either from:
  - `password` set the user password in plain text (do not do that)
  - `passwordFrom` that allow to reference a `secretKeyRef` like for the 
    environment variables of a pod
