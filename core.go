package ice

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func init() {
	LoadConfig()
}

type RequestHandler interface {
	Handle(conn Conn)
}

type Routable interface {
	Route() string
}

type Authorizable interface {
	Authorize(conn Conn) bool
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := makeFromFactory(r.URL.Path)
		if msg == nil {
			http.ServeFile(w, r, "public/")
			return
		}
		handleAPI(r.URL.Path, msg.(RequestHandler), w, r)
	})
	log.Println("Listening at " + host)
	return http.ListenAndServe(host, nil)
}

var factories map[string]*reflect.Type

func Requests(prefix string, reqs ...RequestHandler) {
	if factories == nil {
		factories = make(map[string]*reflect.Type)
	}
	prefix = "/" + strings.Trim(prefix, "/")
	for _, r := range reqs {
		v := reflect.ValueOf(r)
		t := reflect.Indirect(v).Type()
		routable, ok := r.(Routable)
		if ok {
			factories[prefix+routable.Route()] = &t
		} else {
			factories[prefix+t.Name()] = &t
		}
	}
}

func makeFromFactory(key string) RequestHandler {
	t, ok := factories[key]
	if ok == false {
		return nil
	}
	return reflect.New(*t).Interface().(RequestHandler)
}
