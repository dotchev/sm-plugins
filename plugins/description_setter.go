package plugins

import (
	"fmt"

	. "github.com/dotchev/sm-plugins/sm/plugin/rest"
)

// DescriptionSetter is a plugin that sets the description of each service
type DescriptionSetter struct{}

func (d DescriptionSetter) Middleware(route string) Middleware {
	switch route {
	case "osb/catalog":
		return d.catalog
	default:
		return nil
	}
}

func (DescriptionSetter) catalog(req *Request, next Handler) (*Response, error) {

	// call next middleware
	res, err := next(req)

	// modify response
	if err == nil {
		for i, v := range res.Body.Get("services").Array() {
			res.Body.Set(fmt.Sprintf("services.%d.description", i),
				v.Get("name").String()+"-"+v.Get("id").String())
		}
	}
	return res, err
}
