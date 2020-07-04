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
	GmailHost = "smtp.gmail.com"
	GmailPort = 587
	MIME      = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// list of email services
	Confirm       = "confirm"
	PasswordReset = "password-reset"
	Welcome       = "welcome"
)

type emailRequest struct {
	from    string
	to      []string
	subject string
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) serve(r io.Writer, i interface{}) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(r, i)
}

type smtpServer struct {
	host string
	port int
}

func (s *smtpServer) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (s *smtpServer) Send(from string, to []string, message string) error {
	mess := []byte(message)
	pass := os.Getenv("EMAIL_PASSWORD")
	auth := smtp.PlainAuth("", from, pass, s.host)

	if err := smtp.SendMail(s.Address(), auth, from, to, mess); err != nil {
		return errors.Wrap(err, "error sending email")
	}

	return nil
}

type EmailService struct {
	server    *smtpServer
	templates map[string]*templateHandler
}

func NewEmailService() NotificatorService {
	server := &smtpServer{GmailHost, GmailPort}
	templates := map[string]*templateHandler{
		Confirm:       {filename: "confirm.html"},
		PasswordReset: {filename: "password-reset.html"},
		Welcome:       {filename: "welcome.html"},
	}

	return &EmailService{server, templates}
}

func (e *EmailService) SendConfirmEmail(ctx context.Context, name, addr, link string) error {
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

func (e *EmailService) SendPasswordResetEmail(ctx context.Context, addr, link string) error {
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

func (e *EmailService) SendWelcomeEmail(ctx context.Context, name, addr string) error {
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

func (e *EmailService) send(r *emailRequest, values interface{}, service string) error {
	var b *bytes.Buffer
	e.templates[service].serve(b, values)

	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + b.String()
	return e.server.Send(r.from, r.to, body)
}
