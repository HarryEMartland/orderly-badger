package main

import (
	"golang.org/x/net/websocket"
)

type Client struct {
	sendChan        chan interface{}
	websocket       *websocket.Conn
	websocketServer *WebsocketServer
}

func (client *Client)WritePump() {

	for {
		select {
		case message := <-client.sendChan:
			err := websocket.Message.Send(client.websocket, message)
			if (err != nil) {
				client.websocketServer.RemoveWebsocket(client)
				return
			}
		}
	}
}

func (client *Client)sendMessage(message interface{}) {
	client.sendChan <- message
}

func newClient(websocket *websocket.Conn, websocketServer *WebsocketServer) *Client {
	client := Client{
		sendChan: make(chan interface{}, 10),
		websocket:websocket,
		websocketServer:websocketServer,
	}
	return &client
}
