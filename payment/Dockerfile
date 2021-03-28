FROM golang AS builder

ADD ./cmd /payment/cmd
ADD ./pkg /payment/pkg
ADD go.mod /payment
ADD go.sum /payment

WORKDIR /payment/cmd

RUN go get -u ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/payment_service -a -tags netgo -ldflags '-s -w'

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /app/payment_service /app/payment_service
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8086

ENTRYPOINT ["/app/payment_service"]
