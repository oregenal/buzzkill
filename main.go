package main

import (
	"fmt"
	"time"

	"gotelebot/clients/telegram"
)

func main() {
	tgClient, err := telegram.New()
	if err != nil {
		panic(err)
	}
	// fmt.Println(tgClient)

	ubdates, err := tgClient.Update()
	if err != nil {
		panic(err)
	}
	// fmt.Println(ubdates)

	// fmt.Println(ubdates)
	for _, upd := range ubdates {
		fmt.Println(
			time.Unix(upd.Message.Time, 0),
			upd.Message.Text,
			upd.Message.ID,
			"User:", upd.Message.User,
			"Chat:", upd.Message.Chat,
		)
	}
	if err := tgClient.SendMessage("Hello from GoTeleBot!"); err != nil {
		panic(err)
	}

	fmt.Println("In progress...")
}
