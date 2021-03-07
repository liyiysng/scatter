module github.com/liyiysng/scatter

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.2
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/consul/api v1.8.1
	github.com/olivere/elastic/v7 v7.0.22
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/viper v1.7.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
)

//replace google.golang.org/grpc => google.golang.org/grpc v1.28.0
