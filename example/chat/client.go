package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub *Hub //向集线器写入内容的时候，拿到指针（引用）

	conn *websocket.Conn

	send chan []byte
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) //升级的具体过程没实现，但持有这个被劫持的connection，就具备了收发消息的能力
	if err != nil {                          //go一般都会判断是否有错误
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
}
