package ice

import (
	"bufio"
	"log"
	"net/http"
	"time"

	validator "github.com/asaskevich/govalidator"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true }}

func socketLoop(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn := &SocketConn{ws, nil}

	for {
		_, reader, err := ws.NextReader()
		if err != nil {
			continue
		}
		rd := bufio.NewReaderSize(reader, 255)

		cmd, err := rd.ReadBytes([]byte(":")[0])
		if err != nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		factory, ok := handlers[string(cmd[0:len(cmd)-1])]
		if !ok {
			log.Println("not found factory")
			continue
		}
		msg := factory()
		err = ParseJSON(rd, &msg)
		if err != nil {
			continue
		}
		msg.Handle(conn)
	}
}

func handleAPI(cmd string, req Request, w http.ResponseWriter, r *http.Request) {
	err := ParseJSON(r.Body, &req)
	if err != nil {
		return
	}

	conn := &APIConn{w, r, nil}
	if token := r.Header.Get("token"); token != "" && AuthenticateUser != nil {
		AuthenticateUser(token, conn)
	}

	var authorizable Authorizable
	authorizable, _ = req.(Authorizable)
	if authorizable != nil && !authorizable.Authorize(conn) {
		http.Error(w, "Forbidden", 403)
		return
	}

	_, err = validator.ValidateStruct(req)
	if err != nil {
		conn.SendErrors("validation-failed", validator.ErrorsByField(err))
		return
	}

	req.Handle(conn)
}
