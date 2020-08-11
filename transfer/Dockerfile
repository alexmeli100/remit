FROM golang AS builder

ADD ./cmd /transfer/cmd
ADD ./pkg /transfer/pkg
ADD go.mod /transfer
ADD go.sum /transfer

RUN go get -u ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user_service -a -tags netgo -ldflags '-s -w'

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /app/transfer_service /app/transfer_service
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8085

ENTRYPOINT ["/app/transfer_service"]
