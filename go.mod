module github.com/s3nt3/sqlvine

go 1.17

replace google.golang.org/grpc v1.43.0 => google.golang.org/grpc v1.26.0

replace github.com/pingcap/parser v3.1.2+incompatible => github.com/pingcap/tidb/parser v0.0.0-20211230072350-57ef62b8f9e7

require (
	github.com/pingcap/tidb v1.1.0-beta.0.20211230064550-e52f49c4aec2
	github.com/pingcap/tidb/parser v0.0.0-20211230080750-24e7e5d755a9
	go.uber.org/zap v1.19.1
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/coocood/freecache v1.1.1 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/cznic/sortutil v0.0.0-20181122101858-f5f958428db8 // indirect
	github.com/danjacques/gofslock v0.0.0-20191023191349-0a45f885bc37 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgraph-io/ristretto v0.0.1 // indirect
	github.com/dgryski/go-farm v0.0.0-20190423205320-6a90982ecee2 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ngaut/pools v0.0.0-20180318154953-b7bc8c42aac7 // indirect
	github.com/ngaut/sync2 v0.0.0-20141008032647-7a24ed77b2ef // indirect
	github.com/opentracing/basictracer-go v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/pingcap/errors v0.11.5-0.20211009033009-93128226aaa3 // indirect
	github.com/pingcap/failpoint v0.0.0-20210316064728-7acb0f0a3dfd // indirect
	github.com/pingcap/kvproto v0.0.0-20211122024046-03abd340988f // indirect
	github.com/pingcap/log v0.0.0-20211215031037-e024ba4eb0ee // indirect
	github.com/pingcap/parser v3.1.2+incompatible // indirect
	github.com/pingcap/tidb-tools v5.2.2-0.20211019062242-37a8bef2fa17+incompatible // indirect
	github.com/pingcap/tipb v0.0.0-20211227115224-a06a85f9d2a5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.5.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.9.1 // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/shirou/gopsutil v3.21.3+incompatible // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tikv/client-go/v2 v2.0.0-rc.0.20211223062159-300275dee63e // indirect
	github.com/tikv/pd v1.1.0-beta.0.20211118054146-02848d2660ee // indirect
	github.com/tklauser/go-sysconf v0.3.4 // indirect
	github.com/tklauser/numcpus v0.2.1 // indirect
	github.com/twmb/murmur3 v1.1.3 // indirect
	github.com/uber/jaeger-client-go v2.22.1+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/wangjohn/quickselect v0.0.0-20161129230411-ed8402a42d5f // indirect
	go.etcd.io/etcd v0.5.0-alpha.5.0.20210512015243-d19fbe541bf9 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.8 // indirect
	google.golang.org/genproto v0.0.0-20210825212027-de86158e7fda // indirect
	google.golang.org/grpc v1.43.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
