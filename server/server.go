package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type hit_msg struct {
	Hostname   string `json:"hostname"`
	Smooth_num string `json:"smooth_num"`
	Smoothness string `json:"smoothness"`
}

type Client struct {
	Hostname   string `json:"hostname"`
	Init       string `json:"init"`
	Inc        string `json:"inc"`
	Last_itr   string `json:"last_itr"`
	Smooth_num string `json:"smooth_num"`
	Smoothness string `json:"smoothness"`
}

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var smooth_clients map[string]*Client

func genIniMessage() []byte {
	msg := WSMessage{
		Type:    "ini",
		Payload: smooth_clients,
	}
	wsbytes, err := json.Marshal(&msg)
	if err != nil {
		log.Fatal(err)
	}
	return wsbytes
}

func itrHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("itrHandler called")
	var c Client
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(msg, &c)
	if err != nil {
		log.Fatal(err)
	}
	smooth_clients[c.Hostname] = &c
	wsmsg := WSMessage{
		Type:    "itr",
		Payload: c,
	}
	wsbytes, err := json.Marshal(&wsmsg)
	if err != nil {
		log.Fatal(err)
	}
	h.broadcast <- wsbytes
	w.WriteHeader(http.StatusOK)

	// dbg print clients
	fmt.Println("Clients\n-------")
	for h, _ := range smooth_clients {
		fmt.Println(h)
	}
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	var hit hit_msg
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(msg, &hit)
	if err != nil {
		log.Fatal(err)
	}
	if c, exists := smooth_clients[hit.Hostname]; exists {
		c.Smoothness = hit.Smoothness
		c.Smooth_num = hit.Smooth_num
	}
	// Notify web clients
	wsmsg := WSMessage{
		Type:    "hit",
		Payload: hit,
	}
	wsbytes, err := json.Marshal(&wsmsg)
	if err != nil {
		log.Fatal(err)
	}
	h.broadcast <- wsbytes
	w.WriteHeader(http.StatusOK)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	c.send <- genIniMessage()
	go c.writePump()
	c.readPump()
}

func main() {
	smooth_clients = make(map[string]*Client)
	router := mux.NewRouter()
	router.HandleFunc("/itr", itrHandler).Methods("POST")
	router.HandleFunc("/hit", hitHandler).Methods("POST")
	router.HandleFunc("/ws", wsHandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
	go h.run()
	log.Fatal(http.ListenAndServe(":3000", router))
}
