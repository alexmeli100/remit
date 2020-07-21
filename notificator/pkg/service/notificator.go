package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"net/smtp"
	"os"
	"path/filepath"
	"sync"
)

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// list of email services
	Confirm       = "confirm"
	PasswordReset = "password-reset"
	Welcome       = "welcome"
)

type EmailSenser interface {
	Send(from string, to []string, message string) error
}

type emailRequest struct {
	from    string
	to      []string
	subject string
}

type TemplateHandler struct {
	once     sync.Once
	Filename string
	templ    *template.Template
}

func (t *TemplateHandler) serve(r io.Writer, i interface{}) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("EmailTemplates", t.Filename)))
	})

	t.templ.Execute(r, i)
}

type SmtpServer struct {
	host     string
	port     int
	password string
}

func NewSmtpServer(host, password string, port int) *SmtpServer {
	return &SmtpServer{host, port, password}
}

func (s *SmtpServer) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)

}

func (s *SmtpServer) Send(from string, to []string, message string) error {
	mess := []byte(message)
	auth := smtp.PlainAuth("", from, s.password, s.host)

	if err := smtp.SendMail(s.Address(), auth, from, to, mess); err != nil {
		return errors.Wrap(err, "error sending email")
	}

	return nil
}

type NotificationService struct {
	EmailSender    EmailSenser
	EmailTemplates map[string]*TemplateHandler
}

func NewNotificationService(opts ...func(service *NotificationService)) NotificatorService {
	svc := &NotificationService{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (e *NotificationService) SendConfirmEmail(ctx context.Context, name, addr, link string) error {
	values := struct {
		Name string
		Link string
	}{name, link}

	r := &emailRequest{
		from:    os.Getenv("EMAIL_ADDRESS"),
		to:      []string{addr},
		subject: "Activate Account",
	}

	return e.send(r, values, Confirm)

}

func (e *NotificationService) SendPasswordResetEmail(ctx context.Context, addr, link string) error {
	values := struct {
		Link string
	}{link}

	r := &emailRequest{
		from:    os.Getenv("EMAIL_ADDRESS"),
		to:      []string{addr},
		subject: "Password Reset",
	}

	return e.send(r, values, PasswordReset)
}

func (e *NotificationService) SendWelcomeEmail(ctx context.Context, name, addr string) error {
	values := struct {
		Name string
	}{name}

	r := &emailRequest{
		from:    os.Getenv("EMAIL_ADDRESS"),
		to:      []string{addr},
		subject: "Welcome",
	}

	return e.send(r, values, Welcome)
}

func (e *NotificationService) send(r *emailRequest, values interface{}, service string) error {
	var b *bytes.Buffer
	e.EmailTemplates[service].serve(b, values)

	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + b.String()
	return e.EmailSender.Send(r.from, r.to, body)
}
