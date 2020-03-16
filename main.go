package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type postData struct {
	Description string `json: "description"`
	Type        string `json: "type"`
	Time        string `json: "time"`
}

func handler(ctx *macaron.Context) string {
	return "Questa è una: " + ctx.Req.RequestURI
}

func main() {
	go createBot()
	//http.HandleFunc("/", handler)
	m := macaron.Classic() // create a manager of web work application
	fmt.Println("server is online")
	m.Get("/prova", handler) //registering GET method

	m.Post("/", binding.Json(postData{}), func(jp postData) string {
		//return fmt.Sprintf("Auto: %s\nType: %s\nTime: %s\n",jp.Description, jp.Type, jp.Time)
		//Type:modalità di allarme automatica o pre-impostata
		//time orario trigger
		//desription descrive cosa successo {user-text: "ciao", threshold_values: [minVal, Measured_Val, maxVal,]}
		var jsonData []byte
		jsonData, err := json.Marshal(jp)
		if err != nil {
			log.Println(err)
			return ""
		}
		return string(jsonData)

	})
	log.Println(http.ListenAndServe("0.0.0.0:8080", m))
	m.Run()
}
