SHELL := /bin/bash

.DEFAULT_GOAL=controller

SHA := $(shell git log -1 --format='%h')
VERSION := $(shell cat VERSION)

.PHONY: controller
controller:
	operator-sdk build quay.io/blaqkube/mysql-controller:$(SHA)
	docker push quay.io/blaqkube/mysql-controller:$(SHA)
	sed -i 's|REPLACE_IMAGE|quay.io/blaqkube/mysql-controller:'$(SHA)'|g' deploy/operator.yaml
	sed -i 's|quay.io/blaqkube/mysql-controller.*|quay.io/blaqkube/mysql-controller:'$(SHA)'|g' deploy/operator.yaml

.PHONY: operator
operator:
	operator-sdk generate csv --csv-version=$(VERSION) --update-crds
	operator-sdk bundle create quay.io/blaqkube/mysql-operator:$(VERSION) --package mysql-operator --channels alpha --default-channel alpha
	docker push quay.io/blaqkube/mysql-operator:$(VERSION)
