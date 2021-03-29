package service

import (
	"context"
	"go.uber.org/zap"
)

// Middleware describes a service middleware.
type Middleware func(UsersService) UsersService

type loggingMiddleware struct {
	logger *zap.Logger
	next   UsersService
}

func (l loggingMiddleware) DeleteContact(ctx context.Context, contact *Contact) error {
	err := l.next.DeleteContact(ctx, contact)

	defer func() {
		l.logger.Info("DeleteContact", zap.Error(err))
	}()

	return err
}

func (l loggingMiddleware) UpdateUserProfile(ctx context.Context, user *User) (*User, error) {
	u, err := l.next.UpdateUserProfile(ctx, user)

	defer func() {
		l.logger.Info("UpdateUserProfile", zap.Error(err))
	}()

	return u, err
}

func (l loggingMiddleware) UpdateContact(ctx context.Context, contact *Contact) (*Contact, error) {
	c, err := l.next.UpdateContact(ctx, contact)

	defer func() {
		l.logger.Info("UpdateContact", zap.Error(err))
	}()

	return c, err
}

func (l loggingMiddleware) SetUserProfile(ctx context.Context, user *User) (*User, error) {
	u, err := l.next.SetUserProfile(ctx, user)

	defer func() {
		l.logger.Info("SetUserProfile", zap.Error(err))
	}()

	return u, err
}

func (l loggingMiddleware) CreateContact(ctx context.Context, contact *Contact) (*Contact, error) {
	c, err := l.next.CreateContact(ctx, contact)

	defer func() {
		l.logger.Info("CreateContact", zap.Error(err))
	}()

	return c, err
}

func (l loggingMiddleware) GetContacts(ctx context.Context, uid int64) ([]*Contact, error) {
	cs, err := l.next.GetContacts(ctx, uid)

	defer func() {
		l.logger.Info("GetContacts", zap.Error(err))
	}()

	return cs, err
}

func (l loggingMiddleware) CreateUser(ctx context.Context, user *User) (*User, error) {
	u, err := l.next.CreateUser(ctx, user)

	defer func() {
		l.logger.Info("CreateUser", zap.Error(err))
	}()

	return u, err
}

func (l loggingMiddleware) GetUserByID(ctx context.Context, id int64) (*User, error) {
	user, err := l.next.GetUserByID(ctx, id)

	defer func() {
		l.logger.Info("GetUserByID", zap.Error(err))
	}()

	return user, err
}

func (l loggingMiddleware) GetUserByUUID(ctx context.Context, uuid string) (*User, error) {
	user, err := l.next.GetUserByUUID(ctx, uuid)

	defer func() {
		l.logger.Info("GetUserByUUID", zap.Error(err))
	}()

	return user, err
}

func (l loggingMiddleware) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := l.next.GetUserByEmail(ctx, email)

	defer func() {
		l.logger.Info("GetUserByEmail", zap.Error(err))
	}()

	return user, err
}

func (l loggingMiddleware) UpdateEmail(ctx context.Context, user *User) error {
	err := l.next.UpdateEmail(ctx, user)

	defer func() {
		l.logger.Info("UpdateEmail", zap.Error(err))
	}()

	return err
}

// LoggingMiddleware takes a logger as a dependency
// and returns a UsersService Middleware.
func LoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next UsersService) UsersService {
		return &loggingMiddleware{logger, next}
	}
}
