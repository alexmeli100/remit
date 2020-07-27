FROM golang As builder

ADD ./cmd /cmd

WORKDIR /cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user_service -a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /app/user_service /app/user_service
EXPOSE 8081

ENTRYPOINT ["/app/user_service"]


