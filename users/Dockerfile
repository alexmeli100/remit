FROM golang

RUN mkdir -p /go/src/remit

ADD . /go/src/remit

WORKDIR /go/src/remit/users
RUN ls "/go/src/remit/users/cmd"

RUN go mod download

WORKDIR cmd
RUN go install

CMD cmd


