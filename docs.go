package ice

type DocsRequest struct {
	FormValues
	Name string
	Age  int64
	Msg  []string
}

func (r *DocsRequest) Route() string { return "get /docs" }

func (r *DocsRequest) Handle(conn Conn) {
	println(r.Age)
	conn.SendRaw("", []byte("Hello docs! I'm "+r.Name))
}
