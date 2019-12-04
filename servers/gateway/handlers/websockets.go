package handlers

import (
	"UW-Marketplace/servers/gateway/sessions"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket
// Control messages for websocket
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

type SocketStore struct {
	Connections map[int64]*websocket.Conn
	lock        sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// This function's purpose is to reject websocket upgrade requests if the
		// origin of the websockete handshake request is coming from unknown domains.
		// This prevents some random domain from opening up a socket with your server.
		// TODO: make sure you modify this for your HW to check if r.Origin is your host
		return true
	},
}

func (s *SocketStore) InsertConnection(userID int64, conn *websocket.Conn) {
	s.lock.Lock()
	// connId := len(s.Connections)
	// insert socket connection
	// s.Connections = append(s.Connections, conn)
	if s.Connections == nil {
		connMap := make(map[int64]*websocket.Conn)
		connMap[userID] = conn
		s.Connections = connMap
	} else {
		s.Connections[userID] = conn
	}
	s.lock.Unlock()
	// return connId
}

func (s *SocketStore) RemoveConnection(connId int64) {
	s.lock.Lock()
	// insert socket connection
	delete(s.Connections, connId)
	// s.Connections = append(s.Connections[:connId], s.Connections[connId+1:]...)
	s.lock.Unlock()
}

func (s *SocketStore) WriteToAllConnections(messageType int, data []byte) error {
	var writeError error

	for _, conn := range s.Connections {
		writeError = conn.WriteMessage(messageType, data)
		if writeError != nil {
			return writeError
		}
	}

	return nil
}

func (ctx *Context) SocketHandler(w http.ResponseWriter, r *http.Request) {
	var sessionState SessionState
	sessionSid, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &sessionState)
	if err != nil {
		// w.WriteHeader(http.StatusUnauthorized)
		http.Error(w, "Websocket connection refused", http.StatusUnauthorized)
		return
	} else {
		_, err := sessions.ValidateID(sessionSid.String(), ctx.SigningKey)
		if err == nil {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				http.Error(w, "Failed upgrade", http.StatusBadRequest)
				return
			}
			ctx.SockStore.InsertConnection(sessionState.User.ID, conn)
			log.Println("connections", ctx.SockStore.Connections)
			go (func(conn *websocket.Conn) {
				defer conn.Close()
				defer ctx.SockStore.RemoveConnection(sessionState.User.ID)
				for {
					messageType, p, err := conn.ReadMessage()
					log.Println("message type", messageType)
					log.Println("message error", err)
					// 	if err != nil {
					// 		log.Println("Error reading message")
					// 		break
					// 	} else if messageType == 8 {
					// 		log.Println("Close message")
					// 		break
					// 	} else {
					// 		log.Println(p)
					// 	}
					// }
					if messageType == TextMessage || messageType == BinaryMessage {
						fmt.Printf("Client says %v\n", p)
						fmt.Printf("Writing %s to all sockets\n", string(p))
						ctx.SockStore.WriteToAllConnections(TextMessage, append([]byte("Hello from server: "), p...))
					} else if messageType == CloseMessage {
						fmt.Println("Close message received.")
						break
					} else if err != nil {
						fmt.Println("Error reading message.")
						break
					}
				}
			})(conn)
		}
	}
}
