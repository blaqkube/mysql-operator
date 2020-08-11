# Go API client for the MySQL Agent

To update the client, run the command below:

```shell
npx openapi-generator generate \
  -i ../../agent/api/openapi.yaml -g go -o . \
  --git-user-id blaqkube --git-repo-id mysql-operator/mysql-operator/agent

sed -i 's/package\sopenapi/package agent/' *.go
```

