# API Client

To regenerate the API client, you should run the following set of
commands:

```shell
cd pkg/client-agent
npx openapi-generator generate -i ../../agent/mysql-agent.yaml -g go -o .
```
