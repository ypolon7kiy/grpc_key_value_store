Usage:

* Server folder
  * `go run main.go`

* Client folder
  * `go run main.go`



* Install
```
go mod download
go mod tidy
```

* Install protobuf compiler
* `sudo apt install protobuf-compiler`

`protoc --go_out=pb_generated --go_opt=paths=source_relative --go-grpc_out=pb_generated --go-grpc_opt=paths=source_relative contracts/kvstore.proto`

Tests

# `go test -v ./...`