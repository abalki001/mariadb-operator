module github.com/persistentsys/mariadb-operator

go 1.13

require (
	github.com/coreos/prometheus-operator v0.38.0
	github.com/integr8ly/grafana-operator v2.0.0+incompatible
	github.com/operator-framework/operator-sdk v0.16.1-0.20200402040752-a0a5778d9957
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.5.2
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.17.4 // Required by prometheus-operator
)
