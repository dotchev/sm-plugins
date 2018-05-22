package plugins

import (
	"fmt"

	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

func CounterPlugin() *rest.Plugin {
	var cnt Counter
	var plugin rest.Plugin
	plugin.GetCatalog.ProcessResponse = cnt.catalog
	plugin.Provision.ProcessRequest = cnt.preProvision
	plugin.Provision.ProcessResponse = cnt.postProvision
	return &plugin
}

// Counter is a plugin that appends a counter to service and plan ids
type Counter struct {
	counter int
}

func (c *Counter) catalog(req *rest.Request, res *rest.Response) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = p.(error)
		}
	}()

	for _, v := range res.Body.(Object)["services"].(Array) {
		v := v.(Object)
		c.counter++
		v["id"] = fmt.Sprintf("%s.%d", v["id"], c.counter)
	}
	return nil
}

func (c *Counter) preProvision(req *rest.Request) (err error) {
	c.counter++

	// modify request
	b := req.Body.(Object)
	b["service_id"] = fmt.Sprintf("%s.%d", b["service_id"], c.counter)
	b["plan_id"] = fmt.Sprintf("%s.%d", b["plan_id"], c.counter)
	return nil
}

func (c *Counter) postProvision(req *rest.Request, res *rest.Response) (err error) {
	// modify response
	b := res.Body.(Object)
	b["operation"] = fmt.Sprintf("counter.%d", c.counter)
	return nil
}
