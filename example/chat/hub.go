package main

//集线器，管理服务端的所有客户端实例,管理所有连接

type Hub struct {
	clients map[*Client]bool

	broadcast chan []byte //用于存放需要转发的消息

	register chan *Client

	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for { //对于channel的阻塞式for循环，会套一个select
		select { //http连接升级成websocket连接后，
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister: //取消注册，需要判断是否真的有
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message: //尝试给client.send发送消息,可以发送则是正常的活跃链接
				default:
					close(client.send) //如果失败就关闭隧道
					delete(h.clients, client)
				}
			}
		}

	}
}
