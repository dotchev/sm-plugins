package plugins

import (
	"fmt"

	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

// Counter is a plugin that appends a counter to service and plan ids
type Counter struct {
	counter int
}

func (c *Counter) Middleware(route string) rest.Middleware {
	switch route {
	case "osb/catalog":
		return c.catalog
	case "osb/provision":
		return c.provision
	default:
		return nil
	}
}

func (c *Counter) catalog(req *rest.Request, next rest.Handler) (*rest.Response, error) {
	// call next middleware
	res, err := next(req)

	// modify response
	if err == nil {
		services := res.Body.Get("services")
		arr, _ := services.Array()
		for i, _ := range arr {
			service := services.GetIndex(i)
			// same as service.GetPath("metadata", "provider", "name").String()
			provider, _ := service.Get("metadata").Get("provider").Get("name").String()
			if provider != "SAP" {
				c.counter++
				id, _ := service.Get("id").String()
				service.Set("id", fmt.Sprintf("%s.%d", id, c.counter))
			}
		}
	}
	return res, err
}

func (c *Counter) provision(req *rest.Request, next rest.Handler) (*rest.Response, error) {
	c.counter++

	// modify request
	b := req.Body
	serviceid, _ := b.Get("service_id").String()
	planid, _ := b.Get("plan_id").String()
	b.Set("service_id", fmt.Sprintf("%s.%d", serviceid, c.counter))
	b.Set("plan_id", fmt.Sprintf("%s.%d", planid, c.counter))

	// call next middleware
	res, err := next(req)

	// modify response
	if err == nil {
		b := res.Body
		b.Set("operation", fmt.Sprintf("counter.%d", c.counter))
	}
	return res, err
}
