package ice

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Conn interface {
	Send(cmd string, data interface{})
	SendErrors(cmd string, errors interface{})
	SetUser(user *User)
	User() *User
}

type APIConn struct {
	writer http.ResponseWriter
	req    *http.Request
	user   *User
}

func (c *APIConn) Send(cmd string, msg interface{}) {
	EncodeJSON(c.writer, map[string]interface{}{
		"type": cmd,
		"data": msg,
	})
}

func (c *APIConn) SendErrors(cmd string, errors interface{}) {
	EncodeJSON(c.writer, map[string]interface{}{
		"type":   cmd,
		"errors": errors,
	})
}

func (c *APIConn) SetUser(user *User) {
	c.user = user
}

func (c *APIConn) User() *User {
	return c.user
}

type SocketConn struct {
	*websocket.Conn
	user *User
}

func (c *SocketConn) Send(cmd string, msg interface{}) {
	c.WriteJSON(map[string]interface{}{
		"type": cmd,
		"data": msg,
	})
}

func (c *SocketConn) SendErrors(cmd string, errors interface{}) {
	c.WriteJSON(map[string]interface{}{
		"type":   cmd,
		"errors": errors,
	})
}

func (c *SocketConn) SetUser(user *User) {
	c.user = user
}

func (c *SocketConn) User() *User {
	return c.user
}
