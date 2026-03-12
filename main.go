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

	for {
		ubdates, err := tgClient.Update()
		if err != nil {
			panic(err)
		}
		// fmt.Println(ubdates)

		// fmt.Println(ubdates)
		for _, upd := range ubdates {
			fmt.Println(
				time.Unix(upd.Message.Time, 0),
				upd.Message.User.UserName,
				upd.Message.Chat.Type,
				upd.Message.Text,
				// upd.Message.ID,
			)
			if err := tgClient.SendMessage(
				upd.Message.Chat,
				upd.Message.Text,
			); err != nil {
				panic(err)
			}
		}
	}
	// if err := tgClient.SendMessage("Hello from GoTeleBot!"); err != nil {
	// 	panic(err)
	// }

	fmt.Println("In progress...")
}
