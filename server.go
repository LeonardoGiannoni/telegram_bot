package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	tb "gopkg.in/tucnak/telebot.v2"
)

type JSONPost struct {
	Key         string   `json:"key"`
	Type        string   `json:"type"`
	Time        string   `json:"time"`
	Value       []string `json:"value"`
	description string   `json:"description"`
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
	var j map[string]interface{}
	srv.Post("/", binding.Json(Json_Post{}), func(jp JSONPost) string {
		//s, _ := ctx.Req.Body().String()
		/*j = parseJSON(s)*/
		var jsonData JSONPost
		json.Unmarshal([]byte(jp), &json_Data)
		if err != nil {
			log.Println(err)
			return ""
		}
		if json_Data.Value[1]>json_Data.Value[2]{
			b.Send(chat, "Overflow alarm:\n\n"+json_Data.description+"at :"+json_Data.Time)
		}
		else if json_Data.Value[1]<json_Data.Value[0]{
			b.Send(chat, "Underflow alarm:\n\n"+json_Data.description+"at :"+json_Data.Time)
		}
			/*if j["value"][1] > j["value"][2] {
			// over max
			for _, chat := range r.SMembers("alarm_" + string(j["key"])).Result() {
				b.Send(chat, "Overflow alarm:\n\n"+j["description"])
			}
		}*/

		return ""
	})
}
