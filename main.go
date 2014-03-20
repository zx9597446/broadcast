package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/codegangsta/martini"
	"github.com/zx9597446/wssf"
)

const (
	defaultHttpPort  = ":8887"
	defaultWsPort    = ":8888"
	defaultWsRoute   = "/ws"
	defaultHttpRoute = "/bd"
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

func rebuild(values url.Values) map[string]string {
	m := make(map[string]string, 0)
	for k, v := range values {
		m[k] = v[0]
	}
	return m
}

func broadcast(w http.ResponseWriter, r *http.Request) string {
	r.ParseForm()
	m := rebuild(r.Form)
	j, err := json.Marshal(m)
	if err != nil {
		log.Panicln(err)
	}
	wssf.BroadcastMsg(wssf.TextMessage, j)
	log.Printf("broadcasting message: [%s]\n", string(j))
	return "OK"
}

func serveHTTP() {
	m := martini.Classic()
	m.Get(defaultHttpRoute, broadcast)
	log.Printf("serving http %s on port %s\n", defaultHttpRoute, defaultHttpPort)
	log.Fatal(http.ListenAndServe(defaultHttpPort, m))
}

func serveWebsocket() {
	wssf.ServeWS(defaultWsRoute, "GET", "", NewHandler())
	log.Printf("serving websocket %s on port %s\n", defaultWsRoute, defaultWsPort)
	log.Fatal(http.ListenAndServe(defaultWsPort, nil))
}

func main() {
	go serveHTTP()
	serveWebsocket()
}
