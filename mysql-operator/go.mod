module github.com/blaqkube/mysql-operator/mysql-operator

go 1.15

require (
	github.com/antihax/optional v1.0.0
	github.com/blaqkube/mysql-operator/agent v0.0.0-20201219151856-6983e53ab2f7
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.2.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/operator-framework/operator-lib v0.3.0
	github.com/prometheus/common v0.10.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.15.0
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v0.19.4
	sigs.k8s.io/controller-runtime v0.7.0
)
