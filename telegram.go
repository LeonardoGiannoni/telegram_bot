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
		res := strings.Split(m.Payload, " ") //slipt del payload del messagio di testo per eliminare gli spazi 
		if len(res) == 3 {// verificare che la query abbia il numero esatto di argomenti ... "/set val_min val_max "
			if _, ok := mapAttribute[res[0]]; ok {
				jp.Description = res[0] //inserire nel JSON l'attributo su cui lavorare
				for i := 1; i < 3; i++ {
					if _, err := strconv.Atoi(res[i]); err == nil { // conversione ad int per veificare se è un numero 
						if i == 2 {
							//inserire nel JSON il valore minimo e massimo specificato dall'utente  
							if res[1] < res[2] {
								jp.ValueMin = res[1]
								jp.ValueMax = res[2]
							} else {
								jp.ValueMin = res[2]
								jp.ValueMax = res[1]
							}
							SendPostToPersistenceManager(jp)//POST JSON al persistence manager
							b.Send(m.Chat, "New allarm set")
						}
					} else {
						b.Send(m.Chat, "Error: not a number in the query!")
					}
				}
			} else {
				b.Send(m.Chat, "Error: attribute missing!")
			}
		} else {
			b.Send(m.Chat, "Error in the query")
		}
	})

	b.Handle("/show", func(m *tb.Message) {
		res := strings.Split(m.Payload, " ")
		if len(res)==1{// verificare che la query abbia il numero esatto di argomenti ... "/show id_allarmeDaMostrare "
			if _, err := strconv.Atoi(res[0]); err == nil { // conversione ad int per veificare se è un numero
				SendGetToPersistenceManager(res[0])//GET JSON al persistence manager che ritornerà i valori associati ad id_allarmeDaMostrare
				b.Send(m.Chat, "Command show executed")
			} else {
				b.Send(m.Chat, "Error: the id must be a number!")
			}
			
		}else{
			b.Send(m.Chat, "Error in the query, please write only the allarm id")
		}
		
	})

	b.Handle("/delete", func(m *tb.Message) {
		res := strings.Split(m.Payload, " ")
		if len(res)==1{// verificare che la query abbia il numero esatto di argomenti ... "/show id_allarmeDaCancellare "
			if _, err := strconv.Atoi(res[0]); err == nil { // conversione ad int per veificare se è un numero
				SendDeleteToPersistenceManager(res[0])//HTTP DELETE al persistence manager 
		
				b.Send(m.Chat, "Command delete executed")
			} else {
				b.Send(m.Chat, "Error: the id must be a number!")
			}
			
		}else{
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
