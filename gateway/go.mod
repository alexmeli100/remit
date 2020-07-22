module github.com/alexmeli100/remit/gateway

go 1.14

require (
	cloud.google.com/go v0.61.0 // indirect
	cloud.google.com/go/firestore v1.2.0 // indirect
	firebase.google.com/go/v4 v4.0.0
	github.com/alexmeli100/remit/events v0.0.0-20200714195037-ce2de15d6246
	github.com/alexmeli100/remit/notificator v0.0.0-20200714190917-a6f7237e9fa0
	github.com/alexmeli100/remit/users v0.0.0-20200722005611-50449cece51f
	github.com/go-kit/kit v0.10.0
	github.com/gorilla/mux v1.7.4
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	golang.org/x/tools v0.0.0-20200713235242-6acd2ab80ede // indirect
	google.golang.org/api v0.29.0
	google.golang.org/grpc v1.30.0
)
