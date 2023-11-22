package cqhttp

import (
	"chatgpt-qq/config"
	"chatgpt-qq/log"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func RunWsServer() {
	logger := log.NewLogger()
	logger.SetOutput(os.Stderr)
	logger.SetLevel("info")

	http.Handle("/ws", new(wsHandler))
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		fmt.Fprint(writer, "ok")
	})
	cnf := config.GetConf()
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", cnf.CqHttp.WsServerHost, cnf.CqHttp.WsServerPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type wsHandler struct{}

func (*wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accessToken := r.FormValue("access_token")
	cnf := config.GetConf()
	if accessToken != cnf.CqHttp.AccessToken {
		log.Error("权限校验失败")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var upgrader = websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()
	bot := NewBot()
	bot.Conn = conn
	bot.Read()
}
