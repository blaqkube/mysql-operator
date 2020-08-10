module github.com/blaqkube/mysql-operator/mysql-operator

go 1.13

require (
	github.com/aws/aws-sdk-go v1.34.0
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.0
	github.com/johannesboyne/gofakes3 v0.0.0-20200716060623-6b2b4cb092cc
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/operator-framework/operator-lib v0.1.0
	github.com/prometheus/common v0.4.1
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/stretchr/testify v1.5.1
	go.uber.org/zap v1.10.0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
)
