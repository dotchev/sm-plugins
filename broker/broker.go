package broker

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"

	"github.com/dotchev/sm-plugins/sm"
	. "github.com/dotchev/sm-plugins/sm/plugin/json"
)

func Start() *httptest.Server {
	router := mux.NewRouter()
	router.Path("/v2/catalog").Methods("GET").HandlerFunc(getCatalog)
	router.Path("/v2/service_instances/{instance_id}").Methods("PUT").HandlerFunc(provision)
	return httptest.NewServer(router)
}

func getCatalog(res http.ResponseWriter, req *http.Request) {
	catalog := Object{
		"services": Array{
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
		},
	}
	sm.SendJSON(res, 200, catalog)
}

func provision(res http.ResponseWriter, req *http.Request) {
	body := Object{
		"dashboard_url": "http://service-dashboard",
	}
	sm.SendJSON(res, 201, body)
}
