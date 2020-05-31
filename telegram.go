package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func alarmTelegramUI(b *tb.Bot) {
	var jp JSONPost
	var done bool
	mapAttribute := map[string]int32{ // Map of db attribute
		"air_temperature":   0,
		"pressure":          1,
		"humidity_perc":     2,
		"earth_temperature": 3,
		"brightness":        4,
	}
	b.Handle("/set", func(m *tb.Message) {
		logged := r.SIsMember("alarm_chats", m.Chat.Recipient()).Val()
		if !logged {
			b.Send(m.Chat, "You aren't logged!")
			return
		}
		resTemp := strings.Split(m.Payload, " ") //slipt del payload del messagio di testo per eliminare gli spazi
		var res []string
		for _, str := range resTemp {
			if str != "" {
				res = append(res, str)
			}
		}
		if len(res) != 3 {
			b.Send(m.Chat, "Error in the query")
			return
		}
		if _, ok := mapAttribute[res[0]]; !ok {
			b.Send(m.Chat, "Error in the query")
			return
		}
		jp.Description = res[0]
		for i := 1; i < 3; i++ {
			if _, err := strconv.Atoi(res[i]); err == nil {
				if i == 2 {
					intRes1, _ := strconv.Atoi(res[1])
					intRes2, _ := strconv.Atoi(res[2])
					if intRes1 < intRes2 {
						jp.ValueMin = res[1]
						jp.ValueMax = res[2]
					} else {
						jp.ValueMin = res[2]
						jp.ValueMax = res[1]
					}
					done = SendPostToPersistenceManager(jp) //POST JSON al persistence manager
					if done {
						b.Send(m.Chat, "New allarm set")
					} else {
						b.Send(m.Chat, "Error in the request:JSON not send")
					}
				}
			}
		}
	})

	b.Handle("/show", func(m *tb.Message) {
		//logged := r.SIsMember("alarm_chats", m.Chat.Recipient()).Val()
		if false { //salta il login
			b.Send(m.Chat, "You aren't logged!")
			return
		}
		resTemp := strings.Split(m.Payload, " ")
		var res []string
		for _, str := range resTemp {
			if str != "" {
				res = append(res, str)
			}
		}
		if len(res) == 0 { // gestione caso /show... policy:get all alarms in DB
			jp, err := SendGetALLToPersistenceManager()
			if err != nil {
				b.Send(m.Chat, "Err!=nill")
				return
			}

			for _, value := range jp {
				b.Send(m.Chat, "ID_alarm: "+value.ID+" Description: "+value.Description+"\nValueMin: "+value.ValueMin+" ValueMax: "+value.ValueMax)
			}
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if len(res) != 1 { // verificare che la query abbia il numero esatto di argomenti ... "/show id_allarmeDaMostrare "
			b.Send(m.Chat, "Error in the query")
			return
		}
		if _, err := strconv.Atoi(res[0]); err == nil { // conversione ad int per veificare se è un numero
			jp, err2 := SendGetToPersistenceManager(res[0]) //GET al persistence manager che ritornerà i valori associati ad id_allarmeDaMostrare
			if err2 != nil {
				b.Send(m.Chat, "ID doesn't exist")
				return
			}
			b.Send(m.Chat, "ID_alarm: "+jp.ID+" Description: \n "+jp.Description+"ValueMin: "+jp.ValueMin+" ValueMax: "+jp.ValueMax)
		} else {
			b.Send(m.Chat, "Error: the id must be a number!")
		}
	})

	b.Handle("/delete", func(m *tb.Message) {
		logged := r.SIsMember("alarm_chats", m.Chat.Recipient()).Val() //verifica che l'id scritto nella chat dall'utente sia registrato su Redis
		fmt.Println(logged)
		if !logged {
			b.Send(m.Chat, "You aren't logged!")
			fmt.Println("Sei qua effettivamente ")
			return
		}
		resTemp := strings.Split(m.Payload, " ") //slipt del payload del messagio di testo per eliminare gli spazi
		var res []string
		for _, str := range resTemp {
			if str != "" {
				res = append(res, str)
			}
		}

		if len(res) != 1 {
			b.Send(m.Chat, "Error in the query, please write only the allarm id")
		}

		if len(res) == 1 { // verificare che la query abbia il numero esatto di argomenti ... "/show id_allarmeDaCancellare "
			if _, err := strconv.Atoi(res[0]); err == nil { // conversione ad int per veificare se è un numero
				SendDeleteToPersistenceManager(res[0]) //HTTP DELETE al persistence manager

				b.Send(m.Chat, "Command delete executed")
			} else {
				b.Send(m.Chat, "Error: the id must be a number!")
			}

		} else {
			b.Send(m.Chat, "Error in the query, please write only the allarm id")
		}
	})
}

func createBot() *tb.Bot {
	b, _ := tb.NewBot(tb.Settings{
		Token: "944274536:AAGfuzoP3yKOungM3LIng47SI4ZyLZwBCow",
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
		pipe.SAdd(p+"chats", m.Chat.Recipient()) // r.Sismember(alarm_chats, m.Chat.Recipient());
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
