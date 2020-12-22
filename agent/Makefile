SHELL := bash

.DEFAULT_GOAL=build

SHA := $(shell git log -1 --format='%H' . | cut -c1-16)

.PHONY: sha
sha:
	@echo $(SHA)

.PHONY: build
build:
	go build -o mysql-agent .

.PHONY: docker
docker:
	docker build --network=host -t quay.io/blaqkube/mysql-agent:$(SHA) .
	docker push quay.io/blaqkube/mysql-agent:$(SHA)
	sed -i "s/tag\s:=\s\"[^\"]*\"/tag := \"$(SHA)\"/g" ../mysql-operator/pkg/controller/instance/instance_controller.go

.PHONY: generate
generate:
	npx @openapitools/openapi-generator-cli generate \
	  -i mysql-agent.yaml -g go-server -o . \
	  --git-user-id blaqkube \
	  --git-repo-id mysql-operator/agent
	go fmt go/*.go
	cd $$(git rev-parse --show-toplevel) && cd mysql-operator/agent && make api

.PHONY: test
test:
	go test ./... -cover
