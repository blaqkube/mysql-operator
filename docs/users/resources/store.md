# Store

Stores are used to keep backups. For now, only S3 stores are supported. This is
an example of a Store manifest:

```yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Store
metadata:
  name: store-sample
spec:
  backend: s3
  bucket: logs.blaqkube.io
  prefix: /backup/black    
  envs:
  - name: AWS_ACCESS_KEY_ID
    valueFrom:
      secretKeyRef:
        name: store-sample
        key: AWS_ACCESS_KEY_ID
  - name: AWS_SECRET_ACCESS_KEY
    valueFrom:
      secretKeyRef:
        name: store-sample
        key: AWS_SECRET_ACCESS_KEY
  - name: AWS_REGION
    valueFrom:
      secretKeyRef:
        name: store-sample
        key: AWS_REGION
```

The properties are the following:

- `backend` defines the backend. It supports S3 and GCP storage
- `bucket` defines the bucket to store backups
- `prefix` defines the prefix used to prefix backups. It should start with
  `/` and ended without any.
- `envs` contains a set of environment variables that can be used to connect to
  the bucket. It can reference a `name`/`value` pair or a `name`/`valueFrom` 
  pair with a `secretKeyRef` definition.

## Amazon S3

The default storage backend for backups is Amazon S3. In order to use it:

- Set the `backend` property to `s3`
- Rely on environment variables or use a role with the container. Mind that
  both the operator AND the statefulset/pod should have the role profile set 

## GCP storage

It is possible to use GCP storage as a backup store. Below is an example of
a configuration:

```yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Store
metadata:
  name: store-sample
spec:
  backend: gcp
  bucket: logs.blaqkube.io
  prefix: /backup/black    
  envs:
  - name: GOOGLE_APPLICATION_CREDENTIALS
    valueFrom:
      secretKeyRef:
        name: store-sample
        key: GOOGLE_APPLICATION_CREDENTIALS
```

- Set the `backend` property to `gcp`
- You can rely on the `GOOGLE_APPLICATION_CREDENTIALS` environment variable
  and include the token in it, as it is the case above with a secret that
  should look like the one below:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: store-sample
type: Opaque
stringData:
  GOOGLE_APPLICATION_CREDENTIALS: '{ "type": "service_account", ...}'
```

> Note: Using the token inside the `GOOGLE_APPLICATION_CREDENTIALS`
> environment variable has been added in the backend to ease setup.
