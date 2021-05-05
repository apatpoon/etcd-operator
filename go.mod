module github.com/Simonpoon93/etcd-operator

go 1.15

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.2.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/minio/minio-go/v7 v7.0.10
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/net v0.0.0-20210505024714-0287a6fb4125 // indirect
	golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6 // indirect
	google.golang.org/genproto v0.0.0-20210504143626-3b2ad6ccc450 // indirect
	google.golang.org/grpc v1.37.0 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.2

)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
