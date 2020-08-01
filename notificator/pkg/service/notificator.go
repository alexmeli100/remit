package service

import (
	"context"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
)

type Mailer struct {
	apiKey string
}

func NewMailer(apiKey string) *Mailer {
	return &Mailer{apiKey}
}

func (s *Mailer) SendMail(m *mail.SGMailV3) error {
	req := sendgrid.GetRequest(s.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	req.Method = "POST"
	req.Body = mail.GetRequestBody(m)

	_, err := sendgrid.API(req)

	return err
}

type NotificationService struct {
	EmailSender *Mailer
}

func NewNotificationService(opts ...func(service *NotificationService)) NotificatorService {
	svc := &NotificationService{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (e *NotificationService) SendConfirmEmail(ctx context.Context, name, addr, link string) error {
	templateId := os.Getenv("CONFIRM_TEMPLATE_ID")
	tos := []string{addr}
	data := map[string]interface{}{
		"name": name,
		"link": link,
	}

	return e.sendMail(templateId, tos, data)

}

func (e *NotificationService) SendPasswordResetEmail(ctx context.Context, addr, link string) error {
	templateId := os.Getenv("PASSWORD_RESET_TEMPLATE_ID")
	tos := []string{addr}
	data := map[string]interface{}{
		"link": link,
	}

	return e.sendMail(templateId, tos, data)
}

func (e *NotificationService) SendWelcomeEmail(ctx context.Context, name, addr string) error {
	templateId := os.Getenv("WELCOME_TEMPLATE_ID")
	tos := []string{addr}
	data := map[string]interface{}{
		"name": name,
	}

	return e.sendMail(templateId, tos, data)
}

func (e *NotificationService) sendMail(templateId string, tos []string, data map[string]interface{}) error {
	m := mail.NewV3Mail()
	eaddr := os.Getenv("EMAIL_ADDRESS")
	ename := os.Getenv("EMAIL_USERNAME")
	n := mail.NewEmail(ename, eaddr)
	m.SetFrom(n)
	m.SetTemplateID(templateId)
	p := mail.NewPersonalization()

	var t []*mail.Email

	for _, to := range tos {
		t = append(t, mail.NewEmail("", to))
	}

	p.AddTos(t...)

	for k, v := range data {
		p.SetDynamicTemplateData(k, v)
	}

	m.AddPersonalizations(p)
	return e.EmailSender.SendMail(m)
}
