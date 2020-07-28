FROM golang As builder

ADD ./cmd /cmd

WORKDIR /cmd

#RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notificator_service -a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /app/notificator_service /app/notificator_service
EXPOSE 8081

ENTRYPOINT ["/app/notificator_service"]