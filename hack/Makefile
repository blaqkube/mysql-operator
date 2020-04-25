SHELL := /bin/bash

.DEFAULT_GOAL=registry

SHA := $(shell git log -1 --format='%h' .)

.PHONY: registry
registry:
	docker build -t quay.io/blaqkube/operators-catalog:$(SHA) .
	docker push quay.io/blaqkube/operators-catalog:$(SHA)
