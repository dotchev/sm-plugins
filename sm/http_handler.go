package sm

import (
	"log"
	"net/http"
	"reflect"

	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/osb"
)

type HTTPHandler struct {
	osbHandler osb.Handler
}

func NewHTTPHandler(plugins []osb.Plugin, handlerInterface interface{},
	defaultHandler osb.Handler) http.Handler {

	handlerType := reflect.TypeOf(handlerInterface).Elem()
	return HTTPHandler{
		chainHandlers(plugins, handlerType, defaultHandler),
	}
}

func (hh HTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	osbReq, err := readOSBRequest(req)
	if err != nil {
		SendJSON(res, 400, Object{"description": err.Error()})
		return
	}

	osbRes, err := hh.osbHandler(osbReq)
	if err != nil {
		log.Println(err)
		SendJSON(res, 500, Object{"description": err.Error()})
		return
	}

	code := osbRes.StatusCode
	if code == 0 {
		code = 200
	}
	if err := SendJSON(res, code, osbRes.Body); err != nil {
		log.Println(err)
	}
}
