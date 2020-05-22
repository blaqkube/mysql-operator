SHELL := /bin/bash

.DEFAULT_GOAL=build

SHA := $(shell git log -1 --format='%h' .)

.PHONY: build
build:
	docker build --network=host -t quay.io/blaqkube/mysql-agent:$(SHA) .
	docker push quay.io/blaqkube/mysql-agent:$(SHA)

.PHONY: api
api:
	npx openapi-generator generate -i mysql-agent.yaml -g go-server -o . \
	   --git-user-id blaqkube --git-repo-id mysql-operator/agent
