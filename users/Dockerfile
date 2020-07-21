FROM golang As builder

ADD ./cmd /users/cmd
ADD go.mod /users
ADD go.sum /users
ADD ./pkg /users/pkg

WORKDIR /users/cmd

RUN go get -u ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /users/app/user_service -a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /users/app/user_service /app/user_service
EXPOSE 8081

ENTRYPOINT ["/app/user_service"]


