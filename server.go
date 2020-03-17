package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/macaron.v1"
	tb "gopkg.in/tucnak/telebot.v2"
)

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
	var j map[string]interface{}
	srv.Post("/", func(ctx *macaron.Context) string {
		s, _ := ctx.Req.Body().String()
		j = parseJSON(s)
		fmt.Println(j["value"])

		return "OK"
	})
}

func 
