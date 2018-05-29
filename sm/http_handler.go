package sm

import (
	"log"
	"net/http"

	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

type HTTPHandler struct {
	restHandler rest.Handler
}

func NewHTTPHandler(plugins []rest.Plugin, route string,
	defaultHandler rest.Handler) http.Handler {

	return HTTPHandler{
		chain(plugins, route, defaultHandler),
	}
}

func (hh HTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	restReq, err := readRequest(req)
	if err != nil {
		SendJSON(res, 400, Object{"description": err.Error()})
		return
	}

	restRes, err := hh.restHandler(restReq)
	if err != nil {
		log.Println(err)
		SendJSON(res, 500, Object{"description": err.Error()})
		return
	}

	// copy response headers
	for k, v := range restRes.Header {
		if k != "Content-Length" {
			res.Header()[k] = v
		}
	}

	code := restRes.StatusCode
	if code == 0 {
		code = 200
	}
	if err := SendJSON(res, code, restRes.Body); err != nil {
		log.Println(err)
	}
}

func chain(plugins []rest.Plugin, route string, defaultHandler rest.Handler) rest.Handler {
	if len(plugins) == 0 {
		return defaultHandler
	}
	next := chain(plugins[1:], route, defaultHandler)
	m := plugins[0].Middleware(route)
	if m == nil {
		return next
	}
	return func(req *rest.Request) (*rest.Response, error) {
		return m(req, next)
	}
}
