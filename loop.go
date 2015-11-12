package ice

import (
	"bufio"
	"log"
	"net/http"
	"strings"
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
	conn := &SocketConn{ws, (*UserBase)(nil), ""}

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
		_, msg := makeFromFactory("", string(cmd[0:len(cmd)-1]))
		if msg == nil {
			log.Println("not found factory")
			continue
		}
		err = ParseJSON(rd, &msg)
		if err != nil {
			continue
		}
		conn.cmd = string(cmd[0 : len(cmd)-1])
		msg.Handle(conn)
	}
}

func handleAPI(cmd string, req RequestHandler, w http.ResponseWriter, r *http.Request, routeData map[string]string) {
	var err error
	method := strings.ToLower(r.Method)

	if r.Header.Get("content-type") == "application/json" {
		if method != "get" {
			err = ParseJSON(r.Body, &req)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	if fv, ok := req.(FormValuesSetter); ok {
		r.ParseForm()
		for k, v := range routeData {
			r.Form[k] = []string{v}
		}
		fv.SetFormValues(r.Form)
		fv.ParseForm(req)
	}

	conn := &APIConn{w, r, (*UserBase)(nil), cmd}
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
