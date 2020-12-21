# Modify mysql-agent

The `mysql-agent` is an OpenAPI 3.0 API that is used by the MySQL operator
to interact with the MySQL database. It is attached to each instances as a
sidecar and has access to the latest mysql tools. The 

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
make api
```

> **Note**: the command actually run a `make api` in the
> `mysql-operator/agent` directory to implement the agent client too.
