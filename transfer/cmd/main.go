package main

import (
	"context"
	"github.com/alexmeli100/remit/mtn"
	"github.com/alexmeli100/remit/transfer/cmd/service"
	transfer "github.com/alexmeli100/remit/transfer/pkg/service"
	"github.com/go-kit/kit/log"
	"os"
	"os/signal"
)

var logger log.Logger
var momoApiKey string
var momoUserId string
var momoSecret string

func main() {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	initFromEnv()

	ctx, cancel := context.WithCancel(context.Background())
	svc := transfer.NewMobileTransfer(withMtnMomo())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		oscall := <-c
		logger.Log("system call", oscall)
		cancel()
	}()

	service.Run(ctx, svc, logger)
}

func withMtnMomo() func(*transfer.MobileTransfer) {

	return func(t *transfer.MobileTransfer) {
		momo := mtn.CreateMomoApp(nil)
		r := momo.NewRemittance(&mtn.ProductConfig{
			ApiSecret:  momoSecret,
			PrimaryKey: momoApiKey,
			UserId:     momoUserId,
		})

		t.Services[transfer.MTN] = r
	}
}

func initFromEnv() {
	momoApiKey = os.Getenv("MOMO_API_KEY")
	momoUserId = os.Getenv("MOMO_USER_ID")
	momoSecret = os.Getenv("MOMO_USER_SECRET")
}
