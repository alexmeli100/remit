package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/alexmeli100/remit/gateway/app"
	"github.com/alexmeli100/remit/mtn"
	notificator "github.com/alexmeli100/remit/notificator/pkg/service"
	transfer "github.com/alexmeli100/remit/transfer/pkg/service"
	"github.com/alexmeli100/remit/users/pkg/service"
	user "github.com/alexmeli100/remit/users/pkg/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"time"
)

// server options
func serverWithWriteTimeout(t time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.WriteTimeout = t
	}
}

func serverWithReadTimeout(t time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.ReadTimeout = t
	}
}

func serverWithAddress(addr string) func(*http.Server) {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

func serverWithHandler(r *mux.Router) func(*http.Server) {
	return func(s *http.Server) {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"X-Requested-With", "Authorization", "Content-Type"},
			AllowCredentials: true,
		})

		s.Handler = c.Handler(r)
	}
}

func serverWithIdleTimeout(t time.Duration) func(server *http.Server) {
	return func(s *http.Server) {
		s.IdleTimeout = t
	}
}

// application options
func appWithServer(opts ...func(*http.Server)) func(*app.App) error {
	return func(app *app.App) error {
		s := &http.Server{}

		for _, opt := range opts {
			opt(s)
		}

		app.Server = s
		return nil
	}
}

func appWithUserPostgService(_ context.Context, db *sqlx.DB, middleware ...user.Middleware) func(*app.App) error {
	return func(a *app.App) error {
		svc := service.NewPostgService(db)

		for _, m := range middleware {
			svc = m(svc)
		}

		a.UsersService = svc

		return nil
	}
}

func notificatorWithMailer(e *notificator.Mailer) func(*notificator.NotificationService) {
	return func(svc *notificator.NotificationService) {
		svc.EmailSender = e
	}
}

func appWithNotificatorService(_ context.Context, sendGrid string, middleware ...notificator.Middleware) func(*app.App) error {
	return func(a *app.App) error {
		mailSvc := notificator.NewMailer(sendGrid)
		svc := notificator.NewNotificationService(
			notificatorWithMailer(mailSvc))

		for _, m := range middleware {
			svc = m(svc)
		}

		a.Notificator = svc
		return nil
	}
}

func withMtnMomo() func(*transfer.MobileTransfer) {
	return func(t *transfer.MobileTransfer) {
		momoApiKey := os.Getenv("MOMO_API_KEY")
		momoUserId := os.Getenv("MOMO_USER_ID")
		momoSecret := os.Getenv("MOMO_USER_SECRET")

		momo := mtn.CreateMomoApp(nil)
		r := momo.NewRemittance(&mtn.ProductConfig{
			ApiSecret:  momoSecret,
			PrimaryKey: momoApiKey,
			UserId:     momoUserId,
		})

		t.Services[transfer.MTN] = r
	}
}

func appWithTransferService(_ context.Context, middleware ...transfer.Middleware) func(*app.App) error {
	return func(a *app.App) error {
		svc := transfer.NewMobileTransfer(withMtnMomo())

		for _, m := range middleware {
			svc = m(svc)
		}

		a.TransferService = svc
		return nil
	}
}

func appWithFirebase(ctx context.Context, service string) func(*app.App) error {
	return func(a *app.App) error {
		opt := option.WithCredentialsFile(service)

		fireApp, err := firebase.NewApp(ctx, nil, opt)

		if err != nil {
			return err
		}

		a.FireApp = fireApp
		return nil
	}
}
func appWithLogger(logger *zap.Logger) func(*app.App) error {
	return func(a *app.App) error {
		a.Logger = logger
		return nil
	}
}
