package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gotelebot/clients/telegram"
)

func main() {
	tgClient, err := telegram.New()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
		<-exit
		cancel()
	}()

	tgClient.Start(ctx)

	log.Println("Bot stopped...")
}
