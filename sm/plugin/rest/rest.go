package rest

import (
	"encoding/json"

	. "github.com/dotchev/sm-plugins/sm/plugin/json"
)

type Request struct {
	PathParams  map[string]string
	QueryParams map[string]string
	Body        JSON
}

func (r *Request) String() string {
	return stringify(r)
}

type Response struct {
	// StatusCode is the HTTP status code
	StatusCode int

	// Body is the response body parsed as JSON
	Body JSON
}

func (r *Response) String() string {
	return stringify(r)
}

func stringify(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

type Handler struct {
	ProcessRequest  func(*Request) error
	ProcessResponse func(*Request, *Response) error
}

type Plugin struct {
	GetCatalog Handler
	Provision  Handler
	Deproviosn Handler
	Bind       Handler
	Unbind     Handler
}
