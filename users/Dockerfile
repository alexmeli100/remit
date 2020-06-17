FROM golang

RUN mkdir -p /go/src/remit

ADD . /go/src/remit

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run remit/users/cmd  -watch remit/users
