package sm

import (
	"log"

	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/osb"
	"github.com/gorilla/mux"
)

type serviceManager struct {
	*mux.Router
	options Options
}

type Options struct {
	OSBPlugins []osb.Plugin
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
		catalogHandler,
	))
	router.Path("/v2/service_instances/{instance_id}").Methods("PUT").Handler(NewHTTPHandler(
		sm.options.OSBPlugins,
		(*osb.Provisioner)(nil),
		provisionHandler,
	))
}

func catalogHandler(req *osb.Request) (*osb.Response, error) {
	log.Println("Catalog request:", req)
	res := &osb.Response{
		Body: Object{"services": Array{
			Object{
				"name": "dummy",
				"id":   "123",
				"plans": Array{
					Object{
						"name": "default",
						"id":   "789",
					},
				},
			},
		}},
	}
	log.Println("Catalog response:", res)
	return res, nil
}

func provisionHandler(req *osb.Request) (*osb.Response, error) {
	log.Println("Provision request:", req)
	res := &osb.Response{
		Body: Object{
			"dashboard_url": "http://service-dashboard",
		},
	}
	log.Println("Provision response:", res)
	return res, nil
}
