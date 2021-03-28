package app

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/alexmeli100/remit/events"
	eventpb "github.com/alexmeli100/remit/events/pb"
	"github.com/pkg/errors"
)

func (a *App) OnUserCreated(ctx context.Context, data *eventpb.EventData) error {
	u := data.GetUser()

	if u == nil {
		return &events.ErrorShouldAck{Err: "error: got nil user"}
	}

	client, err := a.FireApp.Auth(ctx)

	if err != nil {
		return err
	}

	url, err := client.EmailVerificationLink(ctx, u.Email)

	if auth.IsUserNotFound(err) {
		message := fmt.Sprintf("gateway onUserCreated: user not found %s", u.Email)
		return &events.ErrorShouldAck{Err: message}
	} else if err != nil {
		return errors.Wrap(err, "error getting email confirmation link")
	}

	if err = a.Notificator.SendConfirmEmail(ctx, u.FirstName, u.Email, url); err != nil {
		return errors.Wrap(err, "error sending confirmation email")
	}

	if err = a.Notificator.SendWelcomeEmail(ctx, u.FirstName, u.Email); err != nil {
		return errors.Wrap(err, "error sending welcome email")
	}

	return nil
}

func (a *App) OnPasswordReset(ctx context.Context, data *eventpb.EventData) error {
	u := data.GetUser()
	client, err := a.FireApp.Auth(ctx)

	if err != nil {
		return err
	}

	url, err := client.PasswordResetLink(ctx, u.Email)

	if auth.IsUserNotFound(err) {
		message := fmt.Sprintf("gateway onPasswordReset: user not found %s", u.Email)
		return &events.ErrorShouldAck{Err: message}
	} else if err != nil {
		return errors.Wrap(err, "error getting password reset link")
	}

	if err = a.Notificator.SendPasswordResetEmail(ctx, u.Email, url); err != nil {
		return errors.Wrap(err, "error sending password reset link")
	}

	return nil
}
