protoc -I=. -I=%GOPATH%/src -I=%GOPATH%/src/github.com/gogo/protobuf/protobuf -I=C:/Users/alexm/Documents/code/Go/remit/events/vendor --gogoslick_out=plugins=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. events.proto
protoc-go-inject-tag -input=events.pb.go