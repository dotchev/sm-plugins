package sm

import (
	"fmt"
	"log"

	"github.com/bitly/go-simplejson"

	"github.com/dotchev/sm-plugins/sm/plugin/rest"
	"github.com/gorilla/mux"
	"github.com/parnurzeal/gorequest"
)

var request = gorequest.New()

type serviceManager struct {
	*mux.Router
	options Options
}

type Options struct {
	Plugins   []rest.Plugin
	BrokerURL string
}

func NewServiceManager(options *Options) *serviceManager {
	router := mux.NewRouter()
	sm := &serviceManager{
		router,
		*options,
	}
	sm.mountOSB(router.PathPrefix("/osb/").Subrouter())
	return sm
}

func (sm *serviceManager) mountOSB(router *mux.Router) {
	router.Path("/v2/catalog").Methods("GET").Handler(NewHTTPHandler(
		sm.options.Plugins,
		"osb/catalog",
		sm.catalogHandler,
	))
	router.Path("/v2/service_instances/{instance_id}").Methods("PUT").Handler(NewHTTPHandler(
		sm.options.Plugins,
		"osb/provision",
		sm.provisionHandler,
	))
}

func (sm *serviceManager) catalogHandler(req *rest.Request) (*rest.Response, error) {
	log.Println("Catalog request:", req)

	url := sm.options.BrokerURL + "/v2/catalog"
	log.Printf("Requesting broker at %s", url)
	resp, body, errors := request.Get(url).End()
	if errors != nil {
		log.Println(errors)
		return nil, errors[0]
	}
	j, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return nil, err
	}
	res := &rest.Response{
		Body:       j,
		StatusCode: resp.StatusCode,
	}

	log.Println("Catalog response:", res)
	return res, nil
}

func (sm *serviceManager) provisionHandler(req *rest.Request) (*rest.Response, error) {
	log.Println("Provision request:", req)

	url := fmt.Sprintf("%s/v2/service_instances/%s",
		sm.options.BrokerURL,
		req.PathParams["instance_id"])
	log.Printf("Requesting broker at %s", url)
	resp, body, errors := request.Put(url).Send(req.Body).End()
	if errors != nil {
		log.Println(errors)
		return nil, errors[0]
	}
	j, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return nil, err
	}
	res := &rest.Response{
		Body:       j,
		StatusCode: resp.StatusCode,
	}

	log.Println("Provision response:", res)
	return res, nil
}
