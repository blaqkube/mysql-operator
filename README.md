`blaqkube/mysql-operator` is a Kubernetes 0perator for MySQL Community Server.

[![mysql-operator](https://circleci.com/gh/blaqkube/mysql-operator.svg?style=svg)](https://circleci.com/gh/blaqkube/mysql-operator)

## Features

`blaqkube/mysql-operator` supports
[MySQL Community Edition](https://www.mysql.com/products/community/). It is
built with [operator-sdk](https://sdk.operatorframework.io/) and
[kubebuilder](https://book.kubebuilder.io/).

From a simple manifest, you can:

- Create a MySQL instance 👌
- Add databases to the newly created instance 🏋
- Add users 🎅 to a MySQL instance
- Grant access 🕳 to databases for a user
- Create a backup store 💯 with S3 and GCP
- Generate a backup in the store 💥
- Create a new instance from a backup ❤
- Plug Prometheus and Grafana

## Getting started

Get ready 🚀, check our [online documentation](https://docs.blaqkube.io)

## Contribute

Contributions are welcomed: open issues; ask for help; comment and
request enhancements. Do not hesitate to open PR to correct the documentation
or the code.

If you want to know more about how to modify, build or test the code, check the 
[online developers guide](https://docs.blaqkube.io/developers/welcome) or the
[developers section](docs/developers) of this project.
