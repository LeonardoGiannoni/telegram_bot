package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-macaron/binding"
	"io/ioutil"
	"gopkg.in/macaron.v1"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"errors"
)

type JSONAll []struct {
	Description string `json:"description"`
	ID			string `json:"id"`
	Time        string `json:"time"`
	ValueMin    string `json:"value_min"`
	ValueMax    string `json:"value_max"`
	ValueReal   string `json:"value_real"`
	Key         string `json:"key"`
}



type JSONPost struct {
	Description string `json:"description"`
	ID			string `json:"id"`
	Time        string `json:"time"`
	ValueMin    string `json:"value_min"`
	ValueMax    string `json:"value_max"`
	ValueReal   string `json:"value_real"`
	Key         string `json:"key"`
}

var jg JSONPost
var jAll JSONAll

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

//SendPostToPersistenceManager sends json data to a persistence manager
func SendPostToPersistenceManager(jp JSONPost) bool {
	fmt.Println(jp.Description,jp.ValueMin,jp.ValueMax)
	requestBody, _ := json.Marshal(jp)
	resp, err := http.Post("http://172.28.157.234:8000/alarm", "application/json", bytes.NewBuffer(requestBody)) //write real URL of pers_manager
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode == 201{
		return true
	}else{
		return false
	}
}
//"SendDataToPersistenceManager" sends a param to show the value of the "id" allarm
func SendGetToPersistenceManager(id string) (JSONPost, error)  {
	var url strings.Builder
	var jNull JSONPost
	url.WriteString("http://172.28.157.234:8000/alarm?id=")
	url.WriteString(id)
	resp, err:= http.Get(url.String())
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode == 200{
		bodyBytes, _:=ioutil.ReadAll(resp.Body)
		
		err2 := json.Unmarshal(bodyBytes, &jg)
		if err2!= nil{
			fmt.Println(err)
		}
		return jg,nil

	}else if resp.StatusCode==404{
		return jNull,errors.New("")
	}
	return jNull,errors.New("")
}

func SendGetALLToPersistenceManager() ( JSONAll, error)  {
	var url strings.Builder
	var jNull JSONAll
	url.WriteString("http://172.28.157.234:8000/alarm")
	resp, err:= http.Get(url.String())
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode == 200{
		bodyBytes, _:=ioutil.ReadAll(resp.Body)
		err2 := json.Unmarshal(bodyBytes, &jAll)
		if err2!= nil{
			fmt.Println(err)
		}
		fmt.Println("Json ritornato")
		return jAll,nil

	}else if resp.StatusCode==404{
		return jNull,errors.New("")
	}
	return jNull,errors.New("")
}

//"SendDataToPersistenceManager" sends a param to delete the allarm that has the same "id"
func SendDeleteToPersistenceManager(id string) bool {
	client := &http.Client{}
	/*create a url*/
	var url strings.Builder
	url.WriteString("http://172.28.157.234:8000/alarm?id=")
	url.WriteString(id)
	/**************/
	/*create a request set as a HTTP DELETE*/
	req, err := http.NewRequest("DELETE",url.String(),nil)
	resp, err:=client.Do(req)
	/*************/
	if err!=nil{
		log.Fatalln(err)
	}
	if resp.StatusCode ==200{
		return true
	}else {
		return false
	}
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
