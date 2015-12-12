package ice

import (
	"bufio"
	"log"
	"net/http"
	//	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true }}

func socketLoop(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn := &SocketConn{ws, nil, ""}

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
		log.Println(string(cmd))
		_, msg := makeFromFactory("message", string(cmd[0:len(cmd)-1]))
		if msg == nil {
			log.Println("not found factory")
			continue
		}
		err = ParseJSON(rd, &msg)
		if err != nil {
			continue
		}
		conn.cmd = string(cmd[0 : len(cmd)-1])

		if validator, ok := msg.(RequestValidator); ok {
			validator.Validate(msg)
		}

		data := (msg.(HttpRequestHandler)).Handle(conn)

		if data != nil {
			if resp, ok := data.(Response); ok {
				resp.Execute(conn)
			} else {
				conn.Send("", data)
			}
		}

	}
}

func handleAPI(cmd string, req interface{}, w http.ResponseWriter, r *http.Request, routeData map[string]string) {
	if jv, ok := req.(JsonValuesSetter); ok {
		jv.ParseJSON(req, r)
	}

	if fv, ok := req.(FormValuesSetter); ok {
		r.ParseForm()
		for k, v := range routeData {
			r.Form[k] = []string{v}
		}
		fv.SetFormValues(r.Form)
		fv.ParseForm(req)
	}

	conn := &HttpRequest{w, r, nil, cmd}
	if token := r.Header.Get("token"); token != "" && AuthenticateUser != nil {
		AuthenticateUser(token, conn)
	}
	defer func() {
		resp := recover()
		processHttpResponse(resp, conn)
	}()
	data := HandleHttpRequest(req, conn)
	processHttpResponse(data, conn)
}

func processHttpResponse(data interface{}, conn Conn) {
	if data != nil {
		if resp, ok := data.(Response); ok {
			resp.Execute(conn)
		} else {
			conn.Send("", data)
		}
	}
}

func HandleHttpRequest(req interface{}, conn HttpConn) interface{} {
	/*
		auth, authorizable := req.(Authorizable)
		if authorizable {
			if authorized, forbiddenMessage := auth.Authorize(conn); !authorized {
				return ForbiddenError(forbiddenMessage)
			}
		}
	*/

	if validator, ok := req.(RequestValidator); ok {
		validator.Validate(req)
	}

	if mwp, ok := req.(MiddlewareProvider); ok {
		middlewares := mwp.Middlewares()
		var next func() interface{}
		next = func() interface{} {
			n := middlewares[0]
			middlewares = middlewares[1:]
			return n(req, conn, next)
		}

		middlewares = append(middlewares, func(req interface{}, conn HttpConn, next func() interface{}) interface{} {
			return (req.(HttpRequestHandler)).Handle(conn)
		})

		return next()
	} else {
		return (req.(HttpRequestHandler)).Handle(conn)
	}
}
