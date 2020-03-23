package main

import (
	"fmt"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

mapAttribute := map[string]int32{           // Map of db attribute 
    "air_temperature":  0,
	"pressure": 1,
	"humidity_perc": 2,
	"earth_temperature": 3,
	"brightness": 4,

}

type JSONPost struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Time        string `json:"time"`
	ValueMin    string `json:"value_min"`
	ValueMax    string `json:"value_max"`
	ValueReal   string `json:"value_real"`
	Description string `json:"description"`
}

func alarmTelegramUI(b *tb.Bot) {
var jp JSONPost
	b.Handle("/command", func(m *tb.Message) {
		fmt.Println(m.Payload)
		res := strings.Split(m.Payload, " ")
		for i := 0; i < len(res); i++ {
			if res[i]!= nil{
				if _, ok := mapAttribute[res[i]]; ok {
					jp.Description=mapAttribute[res[i]]
					starvariable:=i
					break
				}
				else{
					b.Send(m.Chat,"error in some query")
				}		
			}
		}
		finish:=0
		for y := starvariable; y < len(res); y++ {
			if res[y]!= nil || finish==2{
				if value, err := strconv.Itoa(strconv.Atoi(res[y])); err == nil {// conversione ad int per veificare se è un int 
					finish++
					if finish == 1
						temp1:=value
					if finish == 2
						temp2:=value //nuova conversione per inserirla in un JSON
				}
			}
		}
		SendDataToPersistenceManager(jp)
	})
	//estrapolo min max e l'attributo e posso fare una post anche dal Handle della funzione alarmTelegramUI
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
