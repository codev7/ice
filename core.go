package ice

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func init() {
	LoadConfig()
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

func Start(host string, notfound http.HandlerFunc) error {
	http.HandleFunc("/connect", socketLoop)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m, msg := makeFromFactory(strings.ToLower(r.Method), r.URL.Path)
		if msg == nil {
			if notfound != nil {
				notfound(w, r)
			} else {
				ServeAsset(w, r)
			}
			return
		}
		handleAPI(r.URL.Path, msg.(RequestHandler), w, r, m)
	})
	log.Println("Listening at " + host)
	return http.ListenAndServe(host, nil)
}

func ServeAsset(w http.ResponseWriter, r *http.Request) {
	if LoadAsset == nil {
		http.NotFound(w, r)
		return
	}
	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}
	data, err := LoadAsset(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, r.URL.Path, time.Now(), bytes.NewReader(data))
}
