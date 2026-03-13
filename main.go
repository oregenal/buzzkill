package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"gotelebot/clients/telegram"
)

func main() {
	tgClient, err := telegram.New()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal)
		signal.Notify(exit, os.Kill, os.Interrupt)
		<-exit
		cancel()
	}()

	tgClient.Start(ctx)

	log.Println("Shutdown...")
}
