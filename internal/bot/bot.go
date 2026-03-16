package bot

import (
	"log"
	"fmt"
	"time"

	"buzzkill/clients/telegram"
)

type Updater interface {
	Update() ([]telegram.Update, error)
}

type Sender interface {
	SendMessage(telegram.Chat, string) error
}

type ContextCheker interface {
	Ctx() error
}

type Processor interface {
	Updater
	Sender
	ContextCheker
}

func Start(p Processor) {
	for {
		updates, err := p.Update()
		if p.Ctx() != nil {
			return
		}
		if err != nil {
			log.Println(err)
		}

		for _, upd := range updates {
			fmt.Println(
				time.Unix(upd.Message.Time, 0),
				upd.Message.User.UserName,
				upd.Message.Chat.Type,
				upd.Message.Text,
				// upd.Message.ID,
			)
			if err := p.SendMessage(
				upd.Message.Chat,
				upd.Message.Text,
			); err != nil {
				log.Println(err)
			}
		}
	}
}
