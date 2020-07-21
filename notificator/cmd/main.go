package main

import (
	"github.com/alexmeli100/remit/notificator/cmd/service"
	notService "github.com/alexmeli100/remit/notificator/pkg/service"
	"log"
	"os"
	"strconv"
)

func main() {
	emailHost := os.Getenv("EMAIL_HOST")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailPort := os.Getenv("EMAIL_PORT")
	p, err := strconv.Atoi(emailPort)

	if err != nil {
		log.Fatalf("invalid port number: %s", emailPort)
	}

	emailServer := notService.NewSmtpServer(emailHost, emailPassword, p)
	templates := map[string]*notService.TemplateHandler{
		notService.Confirm:       {Filename: "confirm.html"},
		notService.PasswordReset: {Filename: "password-reset.html"},
		notService.Welcome:       {Filename: "welcome.html"},
	}

	svc := notService.NewNotificationService(
		notificatorWithEmail(emailServer),
		notificatorWithTemplates(templates))

	service.Run(svc)
}

func notificatorWithEmail(e notService.EmailSenser) func(*notService.NotificationService) {
	return func(svc *notService.NotificationService) {
		svc.EmailSender = e
	}
}

func notificatorWithTemplates(t map[string]*notService.TemplateHandler) func(*notService.NotificationService) {
	return func(svc *notService.NotificationService) {
		svc.EmailTemplates = t
	}
}
