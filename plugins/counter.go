package plugins

import (
	"fmt"

	. "github.com/dotchev/sm-plugins/sm/plugin/json"
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
		for _, v := range res.Body.(Object)["services"].(Array) {
			v := v.(Object)
			c.counter++
			v["id"] = fmt.Sprintf("%s.%d", v["id"], c.counter)
		}
	}
	return res, err
}

func (c *Counter) provision(req *rest.Request, next rest.Handler) (*rest.Response, error) {
	c.counter++

	// modify request
	b := req.Body.(Object)
	b["service_id"] = fmt.Sprintf("%s.%d", b["service_id"], c.counter)
	b["plan_id"] = fmt.Sprintf("%s.%d", b["plan_id"], c.counter)

	// call next middleware
	res, err := next(req)

	// modify response
	if err == nil {
		b = res.Body.(Object)
		b["operation"] = fmt.Sprintf("counter.%d", c.counter)
	}
	return res, err
}
