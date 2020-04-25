SHELL := /bin/bash

.DEFAULT_GOAL=build

SHA := $(shell git log -1 --format='%h')
VERSION := $(shell cat VERSION)

.PHONY: build
build:
	operator-sdk build quay.io/blaqkube/mysql-controller:$(SHA)
	docker push quay.io/blaqkube/mysql-controller:$(SHA)
	sed -i 's|REPLACE_IMAGE|quay.io/blaqkube/mysql-controller:'$(SHA)'|g' deploy/operator.yaml
	sed -i 's|quay.io/blaqkube/mysql-controller.*|quay.io/blaqkube/mysql-controller:'$(SHA)'|g' deploy/operator.yaml

.PHONY: bundle
bundle:
	operator-sdk generate csv --csv-version=$(VERSION) --update-crds
	operator-sdk bundle create quay.io/blaqkube/mysql-operator:v$(VERSION) --package mysql-operator --channels alpha --default-channel alpha
	docker push quay.io/blaqkube/mysql-operator:v$(VERSION)
