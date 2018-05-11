package plugins

import (
	. "github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/osb"
)

// DescriptionSetter is a plugin that sets the description of each service
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
