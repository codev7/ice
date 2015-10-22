package ice

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func init() {
	loadConfig()
}

type Request interface {
	Handle(conn Conn)
}

type Authorizable interface {
	Authorize(conn Conn) bool
}

type RequestFactory func() Request

var handlers map[string]RequestFactory

func Register(m map[string]RequestFactory) {
	if handlers == nil {
		handlers = make(map[string]RequestFactory)
	}
	for cmd, rf := range m {
		handlers[cmd] = rf
	}
}

//parse json from reader
func ParseJSON(reader io.Reader, out interface{}) error {
	d := json.NewDecoder(reader)
	err := d.Decode(&out)
	if err != nil {
		return err
	}
	return nil
}

//encode json
func EncodeJSON(writer io.Writer, data interface{}) error {
	e := json.NewEncoder(writer)
	return e.Encode(data)
}

func Start(host string) error {
	http.HandleFunc("/connect", socketLoop)

	for cmd, rf := range handlers {
		log.Printf("Setting up handler: ", cmd)
		cmd := cmd
		rf := rf
		http.HandleFunc(cmd, func(w http.ResponseWriter, r *http.Request) {
			handleAPI(cmd, rf(), w, r)
		})
	}

	log.Println("Listening at " + host)
	return http.ListenAndServe(host, nil)
}
