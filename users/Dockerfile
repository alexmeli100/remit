FROM golang As builder

ADD ./cmd/main.go /src/main.go
ADD go.mod /src
ADD go.sum /src
ADD ./cmd/service /src/service

WORKDIR /src

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app/user_service -a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /src/app/user_service /app/user_service
EXPOSE 8081

ENTRYPOINT ["/app/user_service"]


