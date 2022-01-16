# Go Data Collection

This service is intented to collect data using grpc using Go lang backend and cassandra DB as storage.

## Dev Setup

```bash
make recreate_test_env
make test
```

### Test Client
run test server
```bash
go run main.go --debug --config internal/config/testdir/config.yaml server
```

in another terminal

```bash
go run main.go --debug --config internal/config/testdir/config.yaml client
```


## TODO:

* convert cli client test  to intigration test.
* enable non timeseries data to be collected.
* setup build pipeline.


## Dependencies caveats

* **protoc-gen-go:** go install google.golang.org/protobuf/cmd/protoc-gen-go
* **protoc-gen-go-grpc** go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
