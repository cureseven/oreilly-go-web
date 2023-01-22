package main

import (
	"log"
	"net/http"

	"github.com/cureseven/oreilly-go-web/trace"
	"github.com/gorilla/websocket"
)

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
	// チャットるーむに参加しようとしているクライアントのためのチャネル
	join chan *client
	// チャットルームから退出しようとしているクライアントのためのチャネル
	leave chan *client
	// 在室しているすべてのクライアント
	// 入退室を繰り返したとき無駄にメモリを使わないためにスライスにしてない
	clients map[*client]bool
	// ログ受け取り
	tracer trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for { // 無限ループ.goroutineなのでおk.1個ずつ処理するため整合性気にしなくていい
		select { // switch的な
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室しました")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました: ", string(msg))
			// すべてのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを送信
					r.tracer.Trace(" -- クライアントに送信されました")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンアップします")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// webSocketコネクションを取得
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	// 入室
	r.join <- client
	// 退室
	defer func() { r.leave <- client }()
	// ゴルーチンで書き込み
	go client.write()
	// 読み取り
	client.read()
}
