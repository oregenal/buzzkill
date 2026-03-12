package main

import (
	"fmt"

	"gotelebot/clients/telegram"
)

func main() {
	tgClient, err := telegram.New()
	if err != nil {
		panic(err)
	}
	// fmt.Println(tgClient)

	tgClient.Start()

	fmt.Println("In progress...")
}
