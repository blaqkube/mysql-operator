# Resources

`blaqkube/mysql-operator` comes with a number of resources to manage MySQL.
Resources available are:

- [`Instance`](resources/instance.md) defines a MySQL instance, its
  attributes and when useful, the backup to use as a source to create the
  instance. 
- [`Store`](resources/store.md) defines backup stores,
- [`Backup`](resources/backup.md) defines a backup requests,
- [`Database`](resources/database.md) defines a database that is part of a
  MySQL instance,
- [`User`](resources/user.md) defines a user part of an instance as well as
  the databases the user can access,
- [`Grant`](resources/grant.md) defines grant for user on a database.
