# Resources

`blaqkube/mysql-operator` comes with a number of resources to manage your
mysql database. Resources available are:

- [`Instance`](manifest/instance.md) defines a MySQL instance, its
  attributes and when useful, the backup to use as a source to create the
  instance. 
- [`Store`](manifest/store.md) defines backup stores,
- [`Backup`](manifest/backup.md) defines a backup request,
- [`Database`](manifest/database.md) defines a database that is part of a
  MySQL instance,
- [`User`](manifest/user.md) defines a user part of an instance as well as
  the databases the user can access,
- [`Grant`](manifest/grant.md) defines grant for user on a database,
