package plugins

import (
	"fmt"

	. "github.com/dotchev/sm-plugins/sm/plugin/rest"
)

// Counter is a plugin that appends a counter to service and plan ids
type Counter struct {
	counter int
}

func (c *Counter) Middleware(route string) Middleware {
	switch route {
	case "osb/catalog":
		return c.catalog
	case "osb/provision":
		return c.provision
	default:
		return nil
	}
}

func (c *Counter) catalog(req *Request, next Handler) (*Response, error) {
	// call next middleware
	res, err := next(req)

	// modify response
	if err == nil {
		for i, v := range res.Body.Get("services").Array() {
			c.counter++
			res.Body.Set(fmt.Sprintf("services.%d.id", i),
				fmt.Sprintf("%s.%d", v.Get("id"), c.counter))
		}
	}
	return res, err
}

func (c *Counter) provision(req *Request, next Handler) (*Response, error) {
	c.counter++

	// modify request
	req.Body.Set("service_id",
		fmt.Sprintf("%s.%d", req.Body.Get("service_id"), c.counter))
	req.Body.Set("plan_id",
		fmt.Sprintf("%s.%d", req.Body.Get("plan_id"), c.counter))

	// call next middleware
	res, err := next(req)

	// modify response
	if err == nil {
		res.Body.Set("operation", fmt.Sprintf("counter.%d", c.counter))
	}
	return res, err
}
