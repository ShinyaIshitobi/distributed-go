module github.com/ShinyaIshitobi/distributed-go

go 1.13

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	github.com/hashicorp/raft-boltdb => github.com/travisjeffery/raft-boltdb v1.0.0
)

require (
	github.com/casbin/casbin/v2 v2.77.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/hashicorp/serf v0.10.1
	github.com/stretchr/testify v1.8.3
	github.com/travisjeffery/go-dynaport v1.0.0
	github.com/tysonmote/gommap v0.0.2
	go.opencensus.io v0.24.0
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.52.0
	google.golang.org/protobuf v1.31.0
)

require (
	cloud.google.com/go v0.109.0 // indirect
	cloud.google.com/go/compute v1.18.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-sockaddr v1.0.5 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/raft v1.5.0
	github.com/hashicorp/raft-boltdb v0.0.0-00010101000000-000000000000
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/miekg/dns v1.1.56 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/soheilhy/cmux v0.1.5
	github.com/tidwall/gjson v1.17.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.9.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20230127162408-596548ed4efa
)
