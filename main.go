package main

import (
	"fmt"
	"net/http"

	"gopkg.in/macaron.v1"
)

func main() {
	b := createBot()
	go start(b)
	//http.HandleFunc("/", handler)
	m := macaron.Classic() // create a manager of web work application
	fmt.Println("server is online")
	createHandleServer(m)
	createHandleDataFromPersistenceManager(m, b)

	http.ListenAndServe("0.0.0.0:8080", m)
	m.Run()
}
