package main

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func CreateBot() {
	b, _ := tb.NewBot(tb.Settings{
		Token:  TelegramSecret,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
}
