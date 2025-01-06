package main

import (
	"flag"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address") //配置启动命令，控制监听端口

func serveHome(w http.ResponseWriter, r *http.Request) { //前端请求页面，可以获得一个html的页面
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./example/chat/home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run() //启动服务端

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { //触发该方法时，将http协议升级为websocket协议
		serveWs(hub, w, r)
	})
}
