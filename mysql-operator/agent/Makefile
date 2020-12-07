SHELL := bash

.DEFAULT_GOAL=api

.PHONY: api
api:
	npx @openapitools/openapi-generator-cli generate \
	  -i ../../agent/mysql-agent.yaml -g go -o . \
	  --git-user-id blaqkube \
	  --git-repo-id mysql-operator/mysql-operator/agent
	go fmt *.go
	sed -i 's/package\sopenapi/package agent/' *.go
