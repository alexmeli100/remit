package main

import (
	"github.com/alexmeli100/remit/notificator/cmd/service"
	notService "github.com/alexmeli100/remit/notificator/pkg/service"
	"os"
)

func main() {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	svc := notService.NewNotificationService(
		notificatorWithMailer(notService.NewMailer(apiKey)))

	service.Run(svc)
}

func notificatorWithMailer(e *notService.Mailer) func(*notService.NotificationService) {
	return func(svc *notService.NotificationService) {
		svc.EmailSender = e
	}
}
