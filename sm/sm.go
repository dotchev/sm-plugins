package sm

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dotchev/sm-plugins/sm/plugin/osb"
	"github.com/gorilla/mux"
	"github.com/parnurzeal/gorequest"
)

var request = gorequest.New()

type serviceManager struct {
	*mux.Router
	options Options
}

type Options struct {
	OSBPlugins []osb.Plugin
	BrokerURL  string
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
		sm.options.OSBPlugins,
		(*osb.CatalogFetcher)(nil),
		sm.catalogHandler,
	))
	router.Path("/v2/service_instances/{instance_id}").Methods("PUT").Handler(NewHTTPHandler(
		sm.options.OSBPlugins,
		(*osb.Provisioner)(nil),
		sm.provisionHandler,
	))
}

func (sm *serviceManager) catalogHandler(req *osb.Request) (*osb.Response, error) {
	log.Println("Catalog request:", req)

	url := sm.options.BrokerURL + "/v2/catalog"
	log.Printf("Requesting broker at %s", url)
	resp, body, err := request.Get(url).End()
	if err != nil {
		log.Println(err)
	}
	var reply interface{}
	json.Unmarshal([]byte(body), &reply)
	res := &osb.Response{
		Body:       reply,
		StatusCode: resp.StatusCode,
	}

	log.Println("Catalog response:", res)
	return res, nil
}

func (sm *serviceManager) provisionHandler(req *osb.Request) (*osb.Response, error) {
	log.Println("Provision request:", req)

	url := fmt.Sprintf("%s/v2/service_instances/%s",
		sm.options.BrokerURL,
		req.Params["instance_id"])
	log.Printf("Requesting broker at %s", url)
	resp, body, err := request.Put(url).Send(req.Body).End()
	if err != nil {
		log.Println(err)
	}
	var reply interface{}
	json.Unmarshal([]byte(body), &reply)
	res := &osb.Response{
		Body:       reply,
		StatusCode: resp.StatusCode,
	}

	log.Println("Provision response:", res)
	return res, nil
}
