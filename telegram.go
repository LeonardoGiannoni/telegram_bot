package main

import (
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func createBot() {
	b, err := tb.NewBot(tb.Settings{
		Token: telegramSecret,
		//URL:   "http://195.129.111.17:8012",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	b.Handle(tb.OnChannelPost, func(m *tb.Message) {
		if m.FromChannel() {
			fmt.Println(m.Chat.Type)
			fmt.Println(m.Chat.Recipient())
			b.Send(m.Chat, "test")
		}
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Start()
}
