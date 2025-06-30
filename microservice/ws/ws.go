package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsHub struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

var hub = wsHub{
	clients: make(map[*websocket.Conn]bool),
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	hub.lock.Lock()
	hub.clients[conn] = true
	hub.lock.Unlock()
	defer func() {
		hub.lock.Lock()
		delete(hub.clients, conn)
		hub.lock.Unlock()
		conn.Close()
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func broadcastWS(message []byte) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	for conn := range hub.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("WebSocket write error:", err)
			conn.Close()
			delete(hub.clients, conn)
		}
	}
} 