FROM golang As builder

ADD ./cmd /users/cmd
ADD ./pkg /users/pkg
ADD go.mod /users
ADD go.sum /users


WORKDIR /users/cmd

RUN go get -u ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user_service -a -tags netgo -ldflags '-s -w'

FROM scratch
MAINTAINER Alex Meli <alexmeli100@gmail.com>
COPY --from=builder /app/user_service /app/user_service
EXPOSE 8081

ENTRYPOINT ["/app/user_service"]


