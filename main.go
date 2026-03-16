package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"buzzkill/clients/telegram"
	"buzzkill/internal/bot"
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
		log.Fatal(err)
	}

	bot.Start(tgClient)

	log.Println("Bot stopped...")
}
