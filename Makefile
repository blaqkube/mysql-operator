
.DEFAULT_GOAL=build
SHELL := /bin/bash
VERSION := $(shell git log -1 --format='%h')
OPVERSION := $(shell cat VERSION)

.PHONY: build
build:
	operator-sdk build quay.io/blaqkube/mysql-controller:$(VERSION)
	docker push quay.io/blaqkube/mysql-controller:$(VERSION)
	sed -i 's|REPLACE_IMAGE|quay.io/blaqkube/mysql-controller:'$(VERSION)'|g' deploy/operator.yaml
	sed -i 's|quay.io/blaqkube/mysql-controller.*|quay.io/blaqkube/mysql-controller:'$(VERSION)'|g' deploy/operator.yaml

prepare:
	operator-sdk generate csv --csv-version $(OPVERSION) --update-crds
