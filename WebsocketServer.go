package main

import (
	"golang.org/x/net/websocket"
)

type WebsocketServer struct {
	addWebsocket    chan *Client
	removeWebsocket chan *Client
	broadcast       chan interface{}
}

type Message struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}

func newWebsocketServer() *WebsocketServer {
	websocketServer := WebsocketServer{
		addWebsocket:make(chan *Client),
		broadcast:make(chan interface{}),
		removeWebsocket:make(chan *Client),
	}
	websocketServer.start()
	return &websocketServer
}

func (websocketServer *WebsocketServer) start() {
	go func() {
		websockets := make(map[*Client]bool)

		for {
			select {
			case newWebsocket := <-websocketServer.addWebsocket:
				websockets[newWebsocket] = true
			case ws := <-websocketServer.removeWebsocket:
				delete(websockets, ws)
			case message := <-websocketServer.broadcast:
				for ws := range websockets {
					ws.sendMessage(message)
				}
			}
		}
	}()
}

func (websocketServer *WebsocketServer) Broadcast(message interface{}) {
	websocketServer.broadcast <- message
}

func (websocketServer *WebsocketServer) RemoveWebsocket(client *Client) {
	websocketServer.removeWebsocket <- client
}

func (websocketServer *WebsocketServer) Handle(ws *websocket.Conn) {
	client := newClient(ws, websocketServer)
	websocketServer.addWebsocket <- client
	client.WritePump()
}
