package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dotchev/sm-plugins/sm"
	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/osb"
)

func main() {
	m := sm.NewServiceManager(&sm.Options{
		[]osb.Plugin{
			DescriptionSetter{},
			&Counter{},
		},
	})

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", m))
}

type DescriptionSetter struct{}

func (DescriptionSetter) FetchCatalog(req *osb.Request, next osb.Handler) (*osb.Response, error) {
	res, err := next(req)

	// modify response
	if err == nil {
		for _, v := range res.Body.(Object)["services"].(Array) {
			v := v.(Object)
			v["description"] = v["name"].(string) + "-" + v["id"].(string)
		}
	}
	return res, err
}

type Counter struct {
	counter int
}

func (a *Counter) FetchCatalog(req *osb.Request, next osb.Handler) (*osb.Response, error) {
	res, err := next(req)

	// modify response
	if err == nil {
		for _, v := range res.Body.(Object)["services"].(Array) {
			v := v.(Object)
			a.counter++
			v["id"] = fmt.Sprintf("%s.%d", v["id"], a.counter)
		}
	}
	return res, err
}

func (a *Counter) Provision(req *osb.Request, next osb.Handler) (*osb.Response, error) {
	a.counter++

	// modify request
	b := req.Body.(Object)
	b["service_id"] = fmt.Sprintf("%s.%d", b["service_id"], a.counter)
	b["plan_id"] = fmt.Sprintf("%s.%d", b["plan_id"], a.counter)

	res, err := next(req)

	// modify response
	if err == nil {
		b = res.Body.(Object)
		b["operation"] = fmt.Sprintf("counter.%d", a.counter)
	}
	return res, err
}
