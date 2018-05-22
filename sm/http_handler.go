package sm

import (
	"log"
	"net/http"

	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

type finalHandler func(req *rest.Request) (*rest.Response, error)

type HTTPHandler struct {
	handler finalHandler
}

func NewHTTPHandler(plugins []*rest.Handler, h finalHandler) http.Handler {
	return HTTPHandler{
		chain(plugins, h),
	}
}

func (hh HTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	restReq, err := readRequest(req)
	if err != nil {
		SendJSON(res, 400, Object{"description": err.Error()})
		return
	}

	restRes, err := hh.handler(restReq)
	if err != nil {
		log.Println(err)
		SendJSON(res, 500, Object{"description": err.Error()})
		return
	}

	code := restRes.StatusCode
	if code == 0 {
		code = 200
	}
	if err := SendJSON(res, code, restRes.Body); err != nil {
		log.Println(err)
	}
}

func chain(plugins []*rest.Handler, h finalHandler) finalHandler {
	return func(req *rest.Request) (*rest.Response, error) {
		// preprocessing
		for i := 0; i < len(plugins); i++ {
			if plugins[i].ProcessRequest != nil {
				if err := plugins[i].ProcessRequest(req); err != nil {
					return nil, err
				}
			}
		}

		// default handler
		res, err := h(req)
		if err != nil {
			return nil, err
		}

		// postprocessing
		for i := len(plugins) - 1; i >= 0; i-- {
			if plugins[i].ProcessResponse != nil {
				if err := plugins[i].ProcessResponse(req, res); err != nil {
					return nil, err
				}
			}
		}

		return res, nil
	}
}
