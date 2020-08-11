:: Install proto3.
:: https://github.com/google/protobuf/releases
:: Update protoc Go bindings via
::  go get -u github.com/golang/protobuf/proto
::  go get -u github.com/golang/protobuf/protoc-gen-go
::
:: See also
::  https://github.com/grpc/grpc-go/tree/master/examples

protoc -I=. -I=%GOPATH%/src -I=%GOPATH%/src/github.com/gogo/protobuf/protobuf --gogoslick_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. payment.proto
protoc-go-inject-tag -input=payment.pb.go