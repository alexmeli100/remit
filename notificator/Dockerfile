FROM golang As builder

ADD ./cmd /notificator/cmd
ADD ./pkg /notificator/pkg
ADD go.mod /notificator
ADD go.sum /notificator

WORKDIR /notificator/cmd

RUN go get -u ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notificator_service -a -tags netgo -ldflags '-s -w'

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /app/notificator_service /app/notificator_service
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8081

ENTRYPOINT ["/app/notificator_service"]