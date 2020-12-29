# Modify mysql-agent

The `mysql-agent` is an OpenAPI 3.0 API that is used by the MySQL operator
to interact with the MySQL database. It is attached to each instances as a
sidecar and has access to the latest MySQL tools.

## How to modify the mysql-agent API?

The `mysql-agent.yaml` located in the `mysql-agent` directory is the OpenAPI
interface for the agent. To modify the agent, modify the API and regenerate  
the Go code with [openapi-generator](https://openapi-generator.tech/). Once
done, you can modify the go code accordingly to your needs.

A typical use of openapi-generator, assuming you've installed `npx` looks
like the command below:

```shell
cd $(git rev-parse --show-toplevel)
cd agent
make generate
```

A typical use of openapi-generator, assuming you've installed it with `npm`
look like the command below:

```shell
cd mysql-agent
npx openapi-generator generate \
  -i mysql-agent.yaml -g go-server -o . \
  --git-user-id blaqkube \
  --git-repo-id mysql-operator/agent
```
