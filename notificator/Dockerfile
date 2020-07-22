FROM golang As builder

ADD ./cmd /notificator/cmd
ADD go.mod /notificator
ADD go.sum /notificator
ADD ./pkg /notificator/pkg

RUN go get ./...

WORKDIR /notificator/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /notificator_service - a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /notificator_service /app/notificator_service
EXPOSE 8081

ENTRYPOINT ["/app/notificator_service"]