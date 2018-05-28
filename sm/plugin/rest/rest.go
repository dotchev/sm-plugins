package rest

import (
	"encoding/json"

	"github.com/bitly/go-simplejson"
)

type Request struct {
	PathParams  map[string]string
	QueryParams map[string]string
	Body        *simplejson.Json
}

func (r *Request) String() string {
	return stringify(r)
}

type Response struct {
	// StatusCode is the HTTP status code
	StatusCode int

	// Body is the response body parsed as JSON
	Body *simplejson.Json
}

func (r *Response) String() string {
	return stringify(r)
}

func stringify(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

type Handler func(*Request) (*Response, error)
type Middleware func(req *Request, next Handler) (*Response, error)

// Plugin handles OSB operations by implementing some of the interfaces below
type Plugin interface {
	// Middleware returns a middleware for given route
	// returns nil, if plugin has no middleware for this route
	Middleware(route string) Middleware
}
