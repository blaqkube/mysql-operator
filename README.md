`blaqkube/mysql-operator` is a Kubernetes 0perator for MySQL Community Server.

[![mysql-operator](https://circleci.com/gh/blaqkube/mysql-operator.svg?style=svg)](https://circleci.com/gh/blaqkube/mysql-operator)

## **Important**

This project has been fun and we have learned a lot from it. Nevertheless, 💔
we have decided to stop it 🖖 and move our MySQL databases to a managed 🌦
service. If you are interested to understand our motivations or react to it,
check [#159](https://github.com/blaqkube/mysql-operator/issues/159). You can
obviously hand it over if you need/want.

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
- Create a new instance from a backup ❤, i.e. clone an instance
- Plug Prometheus and 🧐 Grafana
- Send events to 🤖 Slack

## Getting started

Get ready 🚀, check the [user documentation](./docs/users) or best the
[developer documentation](./docs/developers)
