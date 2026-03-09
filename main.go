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

	res, err := tgClient.Update()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	fmt.Println("In progress...")
}
