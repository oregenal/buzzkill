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

	res, err := tgClient.Update()
	if err != nil {
		panic(err)
	}
	// fmt.Println(res)

	// fmt.Println(time.Unix(int64(res[0].Time)), 0)
	for _, msg := range res {
		fmt.Println(
			time.Unix(msg.Message.Time, 0),
			msg.Message.Text,
			msg.Message.ID,
			"User:", msg.Message.User,
			"Chat:", msg.Message.Chat,
		)
	}
	if err := tgClient.SendMessage("Hello from GoTeleBot!"); err != nil {
		panic(err)
	}

	fmt.Println("In progress...")
}
