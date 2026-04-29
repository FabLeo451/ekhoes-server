package websocket

import (
	"ekhoes-server/utils"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketConnection struct {
	Conn         *websocket.Conn `json:"conn"`
	ConnectionId string          `json:"connectionId"`
	SessionId    string          `json:"sessionId"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	Created      time.Time       `json:"created"`
}

var (
	connections []WebsocketConnection
	mu          sync.Mutex
)

func GetConnections() []WebsocketConnection {
	mu.Lock()
	defer mu.Unlock()

	// copy to avoid extern updates
	result := make([]WebsocketConnection, len(connections))
	copy(result, connections)
	//fmt.Println(result)
	return result
}

func GetConnectionsCount() int32 {
	mu.Lock()
	defer mu.Unlock()
	return int32(len(connections))
}

func AddConnection(wsConn WebsocketConnection) {
	mu.Lock()
	defer mu.Unlock()

	wsConn.ConnectionId = utils.ULID()
	wsConn.Created = time.Now().UTC()

	connections = append(connections, wsConn)

	//fmt.Println(connections)
}

func RemoveConnection(sessionId string) {
	mu.Lock()
	defer mu.Unlock()

	for i, c := range connections {
		if c.SessionId == sessionId {
			connections = append(connections[:i], connections[i+1:]...)
			return
		}
	}
}

/*
	func UpdateConnection(sessionId string, activity string) {
		mu.Lock()
		defer mu.Unlock()

		for i := range connections {
			if connections[i].SessionId == sessionId {
				connections[i].LastActivity = activity
				connections[i].LastActivityTime = time.Now().UTC()
				return
			}
		}
	}
*/

func GetWebsocketConnection(sessionId string) *WebsocketConnection {
	mu.Lock()
	defer mu.Unlock()

	for i := range connections {
		if connections[i].SessionId == sessionId {
			return &connections[i]
		}
	}

	return nil
}

func CloseConnection(conn *websocket.Conn, code int, reason string) {
	_ = conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(
			code,
			reason,
		),
	)
	conn.Close()
}
