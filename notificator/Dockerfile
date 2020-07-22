FROM golang As builder

ADD ./cmd /notificator/cmd
ADD ./pkg /notificator/pkg
ADD go.mod /notificator
ADD go.sum /notificator

WORKDIR /notificator/cmd

RUN go get -u ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /notificator/app/notificator_service -a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /notificator/app/notificator_service /app/notificator_service
EXPOSE 8081

ENTRYPOINT ["/app/notificator_service"]