SHELL := bash

.DEFAULT_GOAL=api

.PHONY: api
api:
	AGTVER=$$(git log --format=%H -1 ../../agent/ |cut -c1-16) && \
	   echo "Upgrading agent to $$AGTVER" && \
	   sed -i "s/^.*DefaultAgentVersion\s=\s.*$$/	DefaultAgentVersion = \"$$AGTVER\"/" ../main.go
	npx @openapitools/openapi-generator-cli generate \
	  -i ../../agent/mysql-agent.yaml -g go -o . \
	  --git-user-id blaqkube \
	  --git-repo-id mysql-operator/mysql-operator/agent
	go fmt *.go
	sed -i 's/package\sopenapi/package agent/' *.go
