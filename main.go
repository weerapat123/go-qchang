package main

import (
	"context"
	"flag"
	"fmt"
	"go-qchang/router"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
}

func main() {
	port := flag.String("port", "8080", "service port number")
	flag.Parse()

	e := router.New()

	go func() {
		e.Logger.Info("starting the server")
		if err := e.Start(fmt.Sprint(":", *port)); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
