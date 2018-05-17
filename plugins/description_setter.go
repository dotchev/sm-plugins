package plugins

import (
	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

// DescriptionSetter is a plugin that sets the description of each service
type DescriptionSetter struct{}

func (d DescriptionSetter) Middleware(route string) rest.Middleware {
	switch route {
	case "osb/catalog":
		return d.catalog
	default:
		return nil
	}
}

func (DescriptionSetter) catalog(req *rest.Request, next rest.Handler) (*rest.Response, error) {

	// call next middleware
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
