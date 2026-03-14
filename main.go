package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"buzzkill/clients/telegram"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
		<-exit
		cancel()
	}()

	tgClient, err := telegram.New(ctx)
	if err != nil {
		panic(err)
	}

	tgClient.Start()

	log.Println("Bot stopped...")
}
