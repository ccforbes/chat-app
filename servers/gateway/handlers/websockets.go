package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/sessions"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

// WSClient a
type WSClient struct {
	conn *websocket.Conn
	u    *users.User
}

// Notifier a
type Notifier struct {
	Connections map[int64][]*WSClient
	lock        sync.Mutex
}

// NewNotifier creates a new Notifier
func NewNotifier() *Notifier {
	return &Notifier{
		Connections: make(map[int64][]*WSClient),
		lock:        sync.Mutex{},
	}
}

// InsertConnection a
func (n *Notifier) InsertConnection(conn *websocket.Conn, u *users.User) int {
	n.lock.Lock()
	connID := len(n.Connections[u.ID])
	newWSClient := &WSClient{
		conn: conn,
		u:    u,
	}
	if len(n.Connections) == 0 {
		n.Connections = make(map[int64][]*WSClient)
	}
	n.Connections[u.ID] = append(n.Connections[u.ID], newWSClient)
	n.lock.Unlock()
	return connID
}

// RemoveConnection a
func (n *Notifier) RemoveConnection(connID int, u *users.User) {
	n.lock.Lock()
	n.Connections[u.ID] = append(n.Connections[u.ID][:connID], n.Connections[u.ID][connID+1:]...)
	n.lock.Unlock()
}

// WriteMessagesToSockets a
func (n *Notifier) WriteMessagesToSockets(messages <-chan amqp.Delivery) {
	for {
		for message := range messages {
			n.lock.Lock()
			data := message.Body
			msg := struct {
				Type    string  `json:"type,omitempty"`
				UserIDs []int64 `json:"userIDs,omitempty"`
			}{}
			json.Unmarshal(data, &msg)
			isPublic := len(msg.UserIDs) == 0

			if isPublic {
				for _, listOfWSClients := range n.Connections {
					for i, currWSClient := range listOfWSClients {
						err := currWSClient.conn.WriteMessage(1, data)
						if err != nil {
							listOfWSClients = append(listOfWSClients[:i], listOfWSClients[i+1:]...)
							currWSClient.conn.Close()
						}
					}
				}
			} else {
				for _, uid := range msg.UserIDs {
					for i, currWSClient := range n.Connections[uid] {
						err := currWSClient.conn.WriteMessage(1, data)
						if err != nil {
							n.Connections[uid] = append(n.Connections[uid][:i], n.Connections[uid][i+1:]...)
							currWSClient.conn.Close()
						}
					}
				}
			}
			n.lock.Unlock()
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//return strings.Contains(r.Header.Get("Origin"), "bopboyz222.xyz")
		return true
	},
}

// WebSocketConnectionHandler a
func (ctx *HandlerCtx) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// check for authentication
	var sessionState SessionState
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &sessionState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", http.StatusUnauthorized)
		return
	}

	connID := ctx.Notifier.InsertConnection(conn, sessionState.User)
	go HandleConnection(conn, connID, sessionState.User, ctx)

}

// HandleConnection asdf
func HandleConnection(conn *websocket.Conn, connID int, u *users.User, ctx *HandlerCtx) {
	defer conn.Close()
	defer ctx.Notifier.RemoveConnection(connID, u)
	for {
		messageType, _, err := conn.ReadMessage()
		if messageType == websocket.CloseMessage || err != nil {
			break
		}
	}
}
