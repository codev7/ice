package ice

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Conn interface {
	Send(cmd string, data interface{})
	SendErrors(cmd string, errors interface{})
	SendRaw(cmd string, data []byte)
	SendView(name string, data interface{}) error
	SetUser(user Identity)
	User() Identity
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
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

func (c *APIConn) SendRaw(cmd string, data []byte) {
	if cmd == "" {
		cmd = c.cmd
	}
	c.writer.Write(data)
}

func (c *APIConn) SendView(name string, data interface{}) error {
	v := NewView(name)
	if v == nil {
		return fmt.Errorf("Could not load view %s", name)
	}
	b, err := v.Render(data)
	if err != nil {
		return err
	}
	http.ServeContent(c.writer, c.req, c.req.URL.Path, time.Now(), b)
	return nil
}

func (c *APIConn) SetUser(user Identity) {
	c.user = user
}

func (c *APIConn) User() Identity {
	return c.user
}

func (c *APIConn) Request() *http.Request              { return c.req }
func (c *APIConn) ResponseWriter() http.ResponseWriter { return c.writer }

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

func (c *SocketConn) SendRaw(cmd string, data []byte) {
	if cmd == "" {
		cmd = c.cmd
	}
	c.WriteMessage(websocket.TextMessage, data)
}

func (c *SocketConn) SendView(name string, data interface{}) error {
	panic("Not suported on web socket transport")
	return nil
}

func (c *SocketConn) SetUser(user Identity) {
	c.user = user
}

func (c *SocketConn) User() Identity {
	return c.user
}

func (c *SocketConn) Request() *http.Request              { panic("not supported") }
func (c *SocketConn) ResponseWriter() http.ResponseWriter { panic("Not supported") }
