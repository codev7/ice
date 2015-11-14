//go:generate go-bindata -pkg ice -prefix "asset" asset
package ice

import (
//"reflect"
)

type APIDescription interface {
	RequestDescription() string
}

type DocsRequest struct {
}

func (r *DocsRequest) Route() string { return "get /docs" }

type apidoc struct {
	Route       string
	Description string
	Request     interface{}
	Fields      []string
	Formats     []string
}

func (r *DocsRequest) Handle(conn Conn) {
	var docs []apidoc
	for _, f := range factories {
		doc := apidoc{
			Route:   f.Method + " " + f.Pattern,
			Request: f.NewRequest(),
		}

		t := *f.Type
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Tag.Get("json") == "-" {
				continue
			}
			fd := field.Name + " " + field.Type.Name() + " " + field.Tag.Get("valid") + " " + field.Tag.Get("description")
			doc.Fields = append(doc.Fields, fd)
		}

		if _, ok := doc.Request.(FormValuesSetter); ok {
			doc.Formats = append(doc.Formats, "Accept input via path, query string and form post")
		}

		if _, ok := doc.Request.(JsonValuesSetter); ok {
			doc.Formats = append(doc.Formats, "Accept input via body in JSON format")
		}

		if d, ok := doc.Request.(APIDescription); ok {
			doc.Description = d.RequestDescription()
		}
		docs = append(docs, doc)
	}
	conn.SendView("docs.html", map[string]interface{}{
		"title": Config.Name,
		"docs":  docs,
	})
}
