package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/zx9597446/wssf"
)

const (
	defaultPort = ":8080"
)

type MyConnectionHandler struct {
}

func (h *MyConnectionHandler) OnConnected(conn *wssf.Connection) {
	//conn.AddTimer(60*time.Second, h.onTimer)
}

func (h *MyConnectionHandler) onTimer() {
	log.Println("onTimer")
}

func (h *MyConnectionHandler) OnDisconnected(conn *wssf.Connection) {
}

func (h *MyConnectionHandler) OnReceived(conn *wssf.Connection, mt int, data []byte) bool {
	return true
}

func (h *MyConnectionHandler) OnError(err error) {
	log.Println(err)
}

func (h *MyConnectionHandler) OnNotify(v interface{}) {
}

func NewHandler() *MyConnectionHandler {
	return &MyConnectionHandler{}
}

func broadcast(params martini.Params) (int, string) {
	msg := params["msg"]
	if msg != "" {
		wssf.BroadcastMsg(wssf.TextMessage, []byte(msg))
		return 200, "OK"
	}
	return 400, "bad request"
}

func main() {
	wssf.ServeWS("/ws", "GET", "", NewHandler())
	m := martini.Classic()
	m.Get("/broadcast", broadcast)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
