FROM golang As builder

ADD ./cmd /gateway/cmd
ADD ./app /gateway/app
ADD go.mod /gateway
ADD go.sum /gateway

WORKDIR /gateway/cmd

RUN go get -u ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway/app/gateway -a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /gateway/app/gateway /app/gateway
RUN apk add --no-cache ca-certificates && update-ca-certificates
EXPOSE 5000

ENTRYPOINT ["/app/gateway"]
