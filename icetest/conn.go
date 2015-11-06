package icetest

import (
	"github.com/nirandas/ice"
)

type TestConn struct {
	Cmd    string
	Data   interface{}
	Errors interface{}
	user   ice.Identity
}

func (c *TestConn) Send(cmd string, data interface{}) {
	c.Cmd = cmd
	c.Data = data
}

func (c *TestConn) SendErrors(cmd string, errors interface{}) {
	c.Cmd = cmd
	c.Errors = errors
}

func (c *TestConn) SetUser(user ice.Identity) {
	c.user = user
}

func (c *TestConn) User() ice.Identity {
	return c.user
}

func NewTestConn(user ice.Identity) *TestConn {
	return &TestConn{
		user: user,
	}
}
