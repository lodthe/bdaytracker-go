package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/internal/application"
	"github.com/lodthe/bdaytracker-go/internal/tgcallback"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	setupLogging()
	tgcallback.Init()

	app, err := application.NewApplication(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create application")
	}

	go app.Run()

	<-stop

	cancel()
	app.Shutdown()
}

func setupLogging() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
