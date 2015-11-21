package ice

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"sync"
	"time"
)

var LoadAsset func(name string) ([]byte, error)
var views map[string]*template.Template
var viewsLock sync.Mutex

func getTemplate(name string) (*template.Template, error) {
	viewsLock.Lock()
	defer viewsLock.Unlock()
	if views == nil {
		views = make(map[string]*template.Template)
	}

	t, ok := views[name]
	if ok {
		return t, nil
	}

	var data []byte
	var err error
	if LoadAsset != nil {
		data, err = LoadAsset(name)
	} else {
		data, err = Asset(name)
	}
	if err != nil {
		data, err = Asset(name)
		if err != nil {
			return nil, err
		}
	}
	t, err = template.New(name).Parse(string(data))
	if err != nil {
		return nil, err
	}
	if Config.Production {
		views[name] = t
	}
	return t, nil
}

type view struct {
	Name     string
	Template *template.Template
}

func NewView(name string) *view {
	v := &view{}
	v.Name = name
	t, err := getTemplate(name)
	if err != nil {
		return nil
	}
	v.Template = t
	return v
}

func (v *view) Render(data interface{}) (io.ReadSeeker, error) {
	var b bytes.Buffer
	err := v.Template.Execute(&b, data)
	return bytes.NewReader(b.Bytes()), err
}

func View(name string, data interface{}) Response {
	v := NewView(name)
	return &viewResponse{v, data}
}

type viewResponse struct {
	*view
	data interface{}
}

func (vr *viewResponse) Execute(conn Conn) {
	b, err := vr.Render(vr.data)
	if err != nil {
		panic(err)
	}
	http.ServeContent(conn.ResponseWriter(), conn.Request(), conn.Request().URL.Path, time.Now(), b)
}
