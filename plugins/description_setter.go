package plugins

import (
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
		services := res.Body.Get("services")
		arr, _ := services.Array()
		for i, _ := range arr {
			v := services.GetIndex(i)
			name, _ := v.Get("name").String()
			id, _ := v.Get("id").String()
			v.Set("description", name+"-"+id)
		}
	}
	return res, err
}
