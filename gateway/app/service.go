package app

import (
	"context"
	"time"
)

// Initialize inirialize the app from the options
func (a *App) Initialize(opts ...func(*App) error) error {
	for _, opt := range opts {
		err := opt(a)

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	errc := make(chan error, 1)

	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			errc <- err
		}
	}()

	select {
	case <-ctx.Done():
		ctx, _ := context.WithTimeout(ctx, time.Second*15)
		return a.Server.Shutdown(ctx)
	case err := <-errc:
		return err
	}
}
