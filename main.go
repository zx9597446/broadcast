package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/codegangsta/martini"
	"github.com/zx9597446/wssf"
)

const (
	defaultHttpPort  = ":8887"
	defaultWsPort    = ":9503"
	defaultWsRoute   = "/ws"
	defaultHttpRoute = "/bd"
	defaultViewRoute = "/view"
)

type wattingMsg struct {
	Msg      string
	Count    int
	Interval int
}

var chWatting = make(chan wattingMsg, 10)

func addMsg(m wattingMsg) {
	chWatting <- m
}

func startSendMsgLoop() {
	for {
		m := <-chWatting
		wssf.BroadcastMsg(wssf.TextMessage, []byte(m.Msg))
		m.Count = m.Count - 1
		log.Printf("broadcasted msg [%s], left [%d]\n", m.Msg, m.Count)
		if m.Count <= 0 {
			continue
		}
		time.AfterFunc(time.Duration(m.Interval)*time.Minute, func() {
			addMsg(m)
		})
	}
}

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

func view() string {
	return htmlView
}

func addPost(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	count := r.FormValue("count")
	interval := r.FormValue("interval")
	icount, _ := strconv.Atoi(count)
	iinterval, _ := strconv.Atoi(interval)
	addMsg(wattingMsg{msg, icount, iinterval})
}

func serveHTTP() {
	m := martini.Classic()
	m.Get(defaultHttpRoute, broadcast)
	m.Get(defaultViewRoute, view)
	m.Post("/add", addPost)
	log.Printf("serving http %s on port %s\n", defaultHttpRoute, defaultHttpPort)
	log.Printf("ui view by: %s\n", defaultViewRoute)
	log.Fatal(http.ListenAndServe(defaultHttpPort, m))
}

func serveWebsocket() {
	wssf.ServeWS(defaultWsRoute, "GET", "", NewHandler())
	log.Printf("serving websocket %s on port %s\n", defaultWsRoute, defaultWsPort)
	log.Fatal(http.ListenAndServe(defaultWsPort, nil))
}

func main() {
	go serveHTTP()
	go startSendMsgLoop()
	serveWebsocket()
}

const htmlView = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">

<head>
<meta content="zh-cn" http-equiv="Content-Language" />
<meta content="text/html; charset=utf-8" http-equiv="Content-Type" />
<title></title>
</head>

<body>
	<form action="/add" method="POST">
	<label>要广播的消息:</label><input name="msg" type="text" size="100" /><br />
	<label>播放次数:</label><input name="count" type="text" /><br />
	<label>间隔时间(单位分钟):</label><input name="interval" type="text" /><br />
	<input name="Submit1" type="submit" value="提交" />
	</form>
	<br />
</body>

</html>

`
