package conn

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Event struct {
	Name string      `json:"e"`
	Data interface{} `json:"d,omitempty"`
}

type Connection struct {
	ID      string
	Session *Session
	Authed  bool
	Ws      *websocket.Conn
}

type ConnEvent struct {
	Conn *Connection
	*Event
}

func NewConnection(session *Session, ws *websocket.Conn) *Connection {
	return &Connection{Session: session, Ws: ws, Authed: false} // default unauthed
}

func (c *Connection) Handle() error {
	s := c.Session

	//err := c.Send(&Event{"doc", map[string]interface{}{
	//	"document": s.Document,
	//	"revision": len(c.Session.Operations),
	//	"clients":  s.Clients,
	//}})
	//if err != nil {
	//	return err
	//}

	s.RegisterConnection(c)

	for {
		e, err := c.ReadEvent()
		if err != nil {
			break
		}

		s.EventChan <- ConnEvent{c, e}
	}

	s.UnRegisterConnection(c)
	// tell all clients some guy is quiting(close the browser)
	c.Broadcast(&Event{"quit", c.ID})

	return nil
}

func (c *Connection) ReadEvent() (*Event, error) {
	_, msg, err := c.Ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	m := &Event{}
	if err = json.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Connection) Send(msg *Event) error {
	if !c.Authed {
		return nil
	}
	j, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err = c.Ws.WriteMessage(websocket.TextMessage, j); err != nil {
		return err
	}
	return nil
}

func (c *Connection) Broadcast(msg *Event) {
	for conn := range c.Session.Connections {
		if conn != c && conn.Authed { //only authed connection can receive broatcast
			conn.Send(msg)
		}
	}
}
