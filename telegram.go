package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

//JSONPOST
type JSONPost struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Time        string `json:"time"`
	ValueMin    string `json:"value_min"`
	ValueMax    string `json:"value_max"`
	ValueReal   string `json:"value_real"`
	Description string `json:"description"`
	Id_val 		string `json:"id_val"`
}

func alarmTelegramUI(b *tb.Bot) {
	var jp JSONPost
	mapAttribute := map[string]int32{ // Map of db attribute
		"air_temperature":   0,
		"pressure":          1,
		"humidity_perc":     2,
		"earth_temperature": 3,
		"brightness":        4,
	}
	b.Handle("/set", func(m *tb.Message) {
		fmt.Println(m.Payload)
		res := strings.Split(m.Payload, " ")
		var temp [3]string
		var y int = 0

		for i := range res {
			if res[i] != "" {
				temp[y] = res[i]
			} else {
				b.Send(m.Chat, "error in some query")
			}
		}

		if len(temp) == 3 {
			if _, ok := mapAttribute[temp[0]]; ok {
				jp.Description = temp[0]
				for i := 1; i < len(temp); i++ {
					if value, err := strconv.Itoa(strconv.Atoi(temp[i])); err == nil { // conversione ad int per veificare se è un int
						if i == 2 {
							if temp[1] < temp[2] {
								jp.ValueMin = temp[1]
								jp.ValueMax = temp[2]
							} else {
								jp.ValueMin = temp[2]
								jp.ValueMax = temp[1]
							}
							SendDataToPersistenceManager(jp)
						}
					} else {
						b.Send(m.Chat, "To much element in the query!")
					}
				}
			} else {
				b.Send(m.Chat, "This attribute doesn't extist")
			}
		} else {
			b.Send(m.Chat, "Error in yhe query")
		}
	})

	b.Handle("/show", func(m *tb.Message) {
		res := strings.Split(m.Payload, " ")
		var temp string
		var y int = 0

		for i := range res {
			if temp {
			temp= res[i]
			} else {
				b.Send(m.Chat, "error in some query")
			}
		
		}
		jp.Id_val=temp
		SendGetToPersistenceManager(temp)
	} 
})

	b.Handle("/delete", func(m *tb.Message) {
			fmt.Println(m.Payload)
			res := strings.Split(m.Payload, " ")
			var temp string
			var y int = 0

			for i := range res {
				if temp {
				temp= res[i]
				} else {
					b.Send(m.Chat, "error in some query")
				}
			
			}
			SendGetToPersistenceManager(temp)
		} 
	})
}

func createBot() *tb.Bot {
	b, _ := tb.NewBot(tb.Settings{
		Token: telegramSecret,
		//URL:   "http://195.129.111.17:8012",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	b.Handle("/hello", func(m *tb.Message) {
		if m.FromChannel() {
			b.Send(m.Chat, "/hello va fatto solo nei canali")
			return
		}
		pipe := r.Pipeline()
		// tutti gli utenti che hanno comunicato al bot in una certa chat
		pipe.SAdd(m.Chat.Recipient(), m.Sender.Recipient())
		// tutte le chat in cui e' stato attivato il bot con /hello
		pipe.SAdd("chats", m.Chat.Recipient())
		// tutte le chat in cui un certo utente e' presente e ha comunicato con il bot
		pipe.SAdd(m.Sender.Recipient()+"_all", m.Chat.Recipient())
		if m.Sender.Username != "" {
			pipe.Set(m.Sender.Recipient(), m.Sender.Username, 0)
			pipe.Set(m.Sender.Username, m.Sender.Recipient(), 0)
			fmt.Println("Username: " + m.Sender.Username)
		}
		pipe.Exec()
		fmt.Println("Chat: " + m.Chat.Recipient())
		fmt.Println("User: " + m.Sender.Recipient())
		b.Send(m.Chat, "test")
	})

	b.Handle("/identify", func(m *tb.Message) {
		key := strings.TrimSpace(m.Text[9:])
		p := "alarm_"
		pipe := r.Pipeline()
		pipe.SAdd(p+"chats", m.Chat.Recipient())
		pipe.SAdd(p+key, m.Chat.Recipient())
		pipe.Exec()
		b.Send(m.Chat, "Registered!")
	})

	alarmTelegramUI(b)

	return b
}

func start(b *tb.Bot) {
	b.Start()
}
