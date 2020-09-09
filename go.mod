module github.com/liangjfblue/gpusher

go 1.13

require (
	github.com/Shopify/sarama v1.26.4
	github.com/chasex/redis-go-cluster v1.0.0
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.5 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/liangjfblue/cheetah v0.0.0-20200430060218-ab065fd3ec64
	github.com/prometheus/client_golang v1.6.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v2 v2.3.0
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

replace github.com/chasex/redis-go-cluster => github.com/chasex/redis-go-cluster v1.0.1-0.20161207023922-222d81891f1d
