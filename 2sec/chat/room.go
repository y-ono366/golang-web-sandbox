package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"dev/trace"
)

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネルです。
	forward chan []byte
	// joinはチャットルームに参加しようとしているクライアントのためのチャネルです。
	join chan *client
	// leaveはチャットルームから退室しようとしているクライアントのためのチャネルです。
	leave chan *client
	// clientsには在室しているすべてのクライアントが保持されます。
	clients map[*client]bool
	// tracerはチャットルーム上で行われた操作のログを受け取ります。
	tracer trace.Tracer
}

func (r *room) run() {
	for {
		//共有されているメモリに対して同期化や変更がいくつか必要な任意の箇所で、このselect文を利用
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
			r.tracer.Trace("[success]client参加")
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("[success]client退室")
		case msg := <-r.forward:
			r.tracer.Trace("[success]受信成功")

			// すべてのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
				// メッセージを送信
				r.tracer.Trace("[success]送信成功")
				default:
				// 送信に失敗
				delete(r.clients, client)
				close(client.send)
				r.tracer.Trace("[faild]送信失敗")
				}
			}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		tracer: trace.Off(),
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
