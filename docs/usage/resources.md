---
title: Resources
---

`blaqkube/mysql-operator` comes with a number of resources to manage your
mysql instances. Resources available are:

- `Instance` defines a MySQL instance, its attributes and when useful, the
  backup to use as a source to create the instance,
- `Database` defines a Database that is part of a MySQL instance,
- `User` defines a User part of an Instance as well as the databases the
  user can access,
- `Store` defines backup stores,
- `Backup` defines a backup request.

## Instance

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
    filePath: "/location/backup-20200526110950.sql"
  maintenance:
    backup: true
    backupStore: docs
    windowStart: 02:30
```

The properties are the following:

- `database` calls a database to create at startup
_ `restore` is used to define the 2 mandatory parameter to base the instance
  on a previous backup:
  - `store` name the Store the backup is located in
  - `filePath` defines the key for the file. It should start with a `/`
- `maintenance` is used to define maintenance parameters for the instance:
  - `backup` is a boolean value that you should set to `true` so that backups
    are automatically scheduled by the operator.
  - `windowStart` defines the start time for the maintenance window. Note that
    this time is defined in UTC
  - `backupStore` defines the store used for automatic backups. The associated
    store should have been previously created for backups to work.

## Database

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

## User

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
  grants:
    - database: red
      accessMode: readwrite
```

The properties are the following:

_ `instance` defines the instance the user is created in
- `username` defines the user name
- `password` set the user password
- `grants` contains a list of `database` and `accessMode` peers:
  - `database` defines the databases the user can access
  - `accessMode` defines the privileges associated with the user. For now
    only `readwrite` is supported.

## Store

Stores are used to keep backups. For now, only S3 stores are supported. This is
an example of a Store manifest:

```yaml
apiVersion: mysql.blaqkube.io/v1alpha1
kind: Store
metadata:
  name: mybackupstore
spec:
  backend: s3
  s3access:
    awsConfig:
      aws_access_key_id: AKIAXXXXXXXXXXXXXXX
      aws_secret_access_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
      region: us-east-1
    bucket: mybucket
    path: /location
```

The properties are the following:

- `backend` defines the backend, for now only s3 is supported
- `s3access` defines the properties for s3
  - `awsConfig` defines some AWS configuration properties. Supported
    properties are:
    - `aws_access_key_id` defines an AWS access key
    - `aws_secret_access_key` defines the associated AWS secret access key
    - `region` defines the region for the s3 bucket
  - `bucket` defines the bucket to store backups
  - `path` defines the prefix used to prefix backups. It should start with
    `/` and ended without any.

## Backup

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

