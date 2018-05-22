package plugins

import (
	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

// DescriptionSetter is a plugin that sets the description of each service
func DescriptionSetter() *rest.Plugin {
	var plugin rest.Plugin
	plugin.GetCatalog.ProcessResponse = setDescription
	return &plugin
}

func setDescription(req *rest.Request, res *rest.Response) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = p.(error)
		}
	}()

	// modify response
	for _, v := range res.Body.(Object)["services"].(Array) {
		v := v.(Object)
		v["description"] = v["name"].(string) + "-" + v["id"].(string)
	}
	return nil
}
