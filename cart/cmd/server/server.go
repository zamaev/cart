package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"route256/cart/internal/app/server"
	"syscall"
)

func main() {
	app := server.NewApp()

	go func() {
		if err := app.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("app shutting down")
	if err := app.Shutdown(context.Background()); err != nil {
		log.Printf("server shutdown error: %v\n", err)
	}
}
