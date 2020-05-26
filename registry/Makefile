SHELL := /bin/bash

.DEFAULT_GOAL=registry

SHA := $(shell git log -1 --format='%h' .)
VERSION := $(shell head -n 1 ../VERSION)

.PHONY: registry
registry:
	opm registry add --bundle-images quay.io/blaqkube/mysql-operator:$(VERSION) -d index.db
	docker build -t quay.io/blaqkube/operators-catalog:$(SHA) .
	docker push quay.io/blaqkube/operators-catalog:$(SHA)
