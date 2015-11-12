package ice

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"
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
	if LoadAsset == nil {
		return nil, fmt.Errorf("Asset Loader (ice.LoadAsset) not set")
	}
	data, err := LoadAsset(name)
	if err != nil {
		return nil, err
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

func (v *view) Bytes(data interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := v.Template.Execute(b, data)
	return b.Bytes(), err
}
