FROM golang As builder

ADD ./cmd /gateway/cmd
ADD ./app /gateway/app

RUN go get ./...

WORKDIR /gateway/cmd

RUN CGO_ENABLED=0 GOOS=linux go build . -o /gateway_service -a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /gateway_service /app/gateway_service
EXPOSE 5000

ENTRYPOINT ["/app/gateway_service"]