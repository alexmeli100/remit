module github.com/alexmeli100/remit

go 1.15

require (
	firebase.google.com/go/v4 v4.1.0
	github.com/alexmeli100/remit/events v0.0.10 // indirect
	github.com/alexmeli100/remit/users v0.0.9 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/uuid v1.1.3
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.9.0
	github.com/nats-io/stan.go v0.7.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/sendgrid/rest v2.6.2+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.7.2+incompatible
	github.com/stripe/stripe-go/v71 v71.48.0
	go.uber.org/zap v1.13.0
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.34.0
)
