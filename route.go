package ice

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type RequestHandler interface {
	Handle(conn Conn)
}

type Routable interface {
	Route() string
}

type Authorizable interface {
	Authorize(conn Conn) bool
}

type JsonValuesSetter interface {
	ParseJSON(tg interface{}, r *http.Request)
}

type JsonValues struct {
}

func (jv *JsonValues) ParseJSON(tg interface{}, r *http.Request) {
}

type FormValuesSetter interface {
	SetFormValues(values url.Values)
	ParseForm(tf interface{})
}

type FormValues struct {
	url.Values `json:"-"`
}

func (r *FormValues) SetFormValues(values url.Values) {
	r.Values = values
}

func (r *FormValues) ParseForm(tg interface{}) {
	m := r.Values
	v := reflect.Indirect(reflect.ValueOf(tg))
	for k, _ := range m {
		f := v.FieldByNameFunc(func(s string) bool {
			return strings.ToLower(s) == strings.ToLower(k)
		})
		if !f.IsValid() || !f.CanSet() {
			continue
		}
		if f.Kind() == reflect.String {
			f.Set(reflect.ValueOf(m.Get(k)))
		}

		if f.Kind() == reflect.Int {
			i, err := strconv.Atoi(m.Get(k))
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Int64 {
			i, err := strconv.ParseInt(m.Get(k), 10, 64)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Int32 {
			i, err := strconv.ParseInt(m.Get(k), 10, 32)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Int16 {
			i, err := strconv.ParseInt(m.Get(k), 10, 16)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Int8 {
			i, err := strconv.ParseInt(m.Get(k), 10, 8)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		}

		if f.Kind() == reflect.Uint {
			i, err := strconv.ParseUint(m.Get(k), 10, 0)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Uint64 {
			i, err := strconv.ParseUint(m.Get(k), 10, 64)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Uint32 {
			i, err := strconv.ParseUint(m.Get(k), 10, 32)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Uint16 {
			i, err := strconv.ParseUint(m.Get(k), 10, 16)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Uint8 {
			i, err := strconv.ParseUint(m.Get(k), 10, 8)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		}

		if f.Kind() == reflect.Float32 {
			i, err := strconv.ParseFloat(m.Get(k), 32)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		} else if f.Kind() == reflect.Float64 {
			i, err := strconv.ParseFloat(m.Get(k), 64)
			if err == nil {
				f.Set(reflect.ValueOf(i))
			}
		}
		if f.Kind() == reflect.Slice && f.Type().Elem().Kind() == reflect.String {
			f.Set(reflect.ValueOf(m[k]))
		}

	}
}

type Route struct {
	Pattern string
	Regexp  *regexp.Regexp
	Method  string
	Type    *reflect.Type
}

func (r *Route) Match(method string, url string) (map[string]string, bool) {
	if r.Method != "" && r.Method != method {
		return nil, false
	}
	m := make(map[string]string)
	names := r.Regexp.SubexpNames()
	if len(names) < 2 {
		return nil, url == r.Pattern
	}
	matches := r.Regexp.FindAllStringSubmatch(url, 1)
	if len(matches) == 0 {
		return nil, false
	}
	for i, sub := range matches {
		m[names[i+1]] = sub[1]
	}
	return m, true
}

func (r *Route) NewRequest() RequestHandler {
	return reflect.New(*r.Type).Interface().(RequestHandler)
}

var factories []*Route

func Requests(prefix string, reqs ...RequestHandler) {
	if prefix != "" {
		prefix = "/" + strings.Trim(prefix, "/")
	}
	for _, r := range reqs {
		v := reflect.ValueOf(r)
		t := reflect.Indirect(v).Type()
		routable, ok := r.(Routable)
		if ok {
			parts := strings.Split(routable.Route(), " ")
			if len(parts) != 2 {
				panic(routable.Route() + " must specify the method and url")
			}
			factories = append(factories, &Route{
				Pattern: prefix + parts[1],
				Method:  parts[0],
				Regexp:  preparePattern(prefix + parts[1]),
				Type:    &t,
			})
		} else {
			panic(t.Name() + " must implement Routable interface")
		}
	}
}

func preparePattern(pattern string) *regexp.Regexp {
	exp := `{([a-z]+)(:([^}]+))?}`
	r := regexp.MustCompile(exp)
	re := r.ReplaceAllString(pattern, "(?P<$1>$2)")
	re = "^" + strings.Replace(re, ">)", ">\\w+)", -1) + "$"
	return regexp.MustCompile(re)
}

func makeFromFactory(method string, url string) (map[string]string, RequestHandler) {
	for _, r := range factories {
		if m, ok := r.Match(method, url); ok {
			return m, r.NewRequest()
		}
	}
	return nil, nil
}
