package osb

import (
	"encoding/json"

	. "github.com/dotchev/sm-plugins/sm/plugin/json"
)

type Request struct {
	// Params contains path and query parameters
	// all keys are lower case
	// if changed, the request sent to the real broker is changed too
	// since we define the path parameters, all keys that are not path parameters
	// are serialized back to query parameters
	Params map[string]string

	// Body is the request body parsed as JSON
	Body JSON
}

func (r *Request) String() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

type Response struct {
	// StatusCode is the HTTP status code
	StatusCode int

	// Body is the response body parsed as JSON
	Body JSON
}

func (r *Response) String() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

// Handler handles an OSB operation
type Handler func(*Request) (*Response, error)

// Plugin handles OSB operations by implementing some of the interfaces below
type Plugin interface{}

type CatalogFetcher interface {
	FetchCatalog(req *Request, next Handler) (*Response, error)
}

type Provisioner interface {
	Provision(req *Request, next Handler) (*Response, error)
}

type Deprovisioner interface {
	Deprovision(req *Request, next Handler) (*Response, error)
}

type ServiceUpdater interface {
	UpdateService(req *Request, next Handler) (*Response, error)
}

type ServiceFetcher interface {
	FetchService(req *Request, next Handler) (*Response, error)
}

type Binder interface {
	Bind(req *Request, next Handler) (*Response, error)
}

type Unbinder interface {
	Unbind(req *Request, next Handler) (*Response, error)
}

type BindingFetcher interface {
	FetchBinding(req *Request, next Handler) (*Response, error)
}
