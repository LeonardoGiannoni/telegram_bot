package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	tb "gopkg.in/tucnak/telebot.v2"
)

type JSONPost struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Time        string `json:"time"`
	ValueMin    string `json:"value_min"`
	ValueMax    string `json:"value_max"`
	ValueReal   string `json:"value_real"`
	Description string `json:"description"`
}

type chatTarget struct {
	payload string
}

func (c *chatTarget) Recipient() string {
	return c.payload
}

func parseJSON(s string) map[string]interface{} {
	var res map[string]interface{}
	json.Unmarshal([]byte(s), &res)

	return res
}

func createHandleServer(srv *macaron.Macaron) {
	srv.Get("/", func(ctx *macaron.Context) string {
		return "Working\n"
	})
}

//SendDataToPersistenceManager sends json data to a persistence manager
func SendDataToPersistenceManager(m *tb.Message) {

	postObj := JSONPost{Key: m.Payload, Type: m.Payload, Time: "test", ValueMin: "t", ValueReal: "t", ValueMax: "r", Description: "w"}
	requestBody, _ := json.Marshal(postObj)
	resp, err := http.Post("URL_PYTHON:8080/test", "application/json", bytes.NewBuffer(requestBody)) //write real URL of pers_manager
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

}

func createHandleDataFromPersistenceManager(srv *macaron.Macaron, b *tb.Bot) {
	/*
		{
			"type": (auto|manual),
			"description": "stringa",
			"key": 1234
			"time": 2020-03-17
			"value": [min, actualValue, max]
		}
	*/
	srv.Post("/", binding.Json(JSONPost{}), func(jp JSONPost) string {
		//s, _ := ctx.Req.Body().String()
		/*j = parseJSON(s)*/
		fmt.Println(jp.Description)
		chatsToWriteTo := r.SMembers("alarm_" + jp.Key).Val()
		msg := ""
		if jp.ValueReal > jp.ValueMax {
			msg = "Overflow alarm\n\n" + jp.Description + "\ntime: " + jp.Time
		} else {
			if jp.ValueReal < jp.ValueMin {
				msg = "Underflow alarm\n\n" + jp.Description + "\ntime: " + jp.Time
			}
		}
		for _, chat := range chatsToWriteTo {
			var target tb.Chat
			val, _ := strconv.ParseInt(chat, 10, 64)
			target.ID = val
			b.Send(&target, msg)
		}

		return ""
	})
}
