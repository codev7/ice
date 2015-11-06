package ice

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Conn interface {
	Send(cmd string, data interface{})
	SendErrors(cmd string, errors interface{})
	SetUser(user Identity)
	User() Identity
}

type APIConn struct {
	writer http.ResponseWriter
	req    *http.Request
	user   Identity
	cmd    string
}

func (c *APIConn) Send(cmd string, msg interface{}) {
	if cmd == "" {
		cmd = c.cmd
	}
	EncodeJSON(c.writer, map[string]interface{}{
		"type": cmd,
		"data": msg,
	})
}

func (c *APIConn) SendErrors(cmd string, errors interface{}) {
	if cmd == "" {
		cmd = c.cmd
	}
	EncodeJSON(c.writer, map[string]interface{}{
		"type":   cmd,
		"errors": errors,
	})
}

func (c *APIConn) SetUser(user Identity) {
	c.user = user
}

func (c *APIConn) User() Identity {
	return c.user
}

type SocketConn struct {
	*websocket.Conn
	user Identity
	cmd  string
}

func (c *SocketConn) Send(cmd string, msg interface{}) {
	if cmd == "" {
		cmd = c.cmd
	}

	c.WriteJSON(map[string]interface{}{
		"type": cmd,
		"data": msg,
	})
}

func (c *SocketConn) SendErrors(cmd string, errors interface{}) {
	if cmd == "" {
		cmd = c.cmd
	}

	c.WriteJSON(map[string]interface{}{
		"type":   cmd,
		"errors": errors,
	})
}

func (c *SocketConn) SetUser(user Identity) {
	c.user = user
}

func (c *SocketConn) User() Identity {
	return c.user
}
