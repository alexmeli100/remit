FROM golang As builder

ADD ./cmd /gateway/cmd
ADD ./app /gateway/app
ADD go.mod /gateway
ADD go.sum /gateway

WORKDIR /gateway/cmd

RUN go get -u ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gateway -a -tags netgo -ldflags '-s -w'

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /app/gateway /app/gateway
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8083

ENTRYPOINT ["/app/gateway"]
