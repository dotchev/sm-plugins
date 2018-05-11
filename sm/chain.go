package sm

import (
	"reflect"

	"github.com/dotchev/sm-plugins/sm/plugin/osb"
)

func chainHandlers(plugins []osb.Plugin, handlerInterface reflect.Type, defaultHandler osb.Handler) osb.Handler {
	methodName := handlerInterface.Method(0).Name
	var chain func([]osb.Plugin) osb.Handler // yes, go sucks
	chain = func(plugins []osb.Plugin) osb.Handler {
		if len(plugins) == 0 {
			return defaultHandler
		}
		next := chain(plugins[1:])
		pluginType := reflect.TypeOf(plugins[0])
		if !pluginType.Implements(handlerInterface) {
			return next
		}
		method, _ := pluginType.MethodByName(methodName)
		return func(req *osb.Request) (res *osb.Response, err error) {
			ret := method.Func.Call([]reflect.Value{reflect.ValueOf(plugins[0]), reflect.ValueOf(req), reflect.ValueOf(next)})
			res, _ = ret[0].Interface().(*osb.Response)
			err, _ = ret[1].Interface().(error)
			return
		}
	}
	return chain(plugins)
}

// No reflection, but should be duplicated for each OSB operation
// func chainCatalogHandlers(plugins []osb.Plugin, defaultHandler osb.Handler) osb.Handler {
// 	var chain func([]osb.Plugin) osb.Handler // yes, go sucks
// 	chain = func(plugins []osb.Plugin) osb.Handler {
// 		if len(plugins) == 0 {
// 			return defaultHandler
// 		}
// 		next := chain(plugins[1:])
// 		fetcher, ok := plugins[0].(osb.CatalogFetcher)
// 		if !ok {
// 			return next
// 		}
// 		f := fetcher.FetchCatalog
// 		return func(req *osb.Request) (*osb.Response, error) {
// 			return f(req, next)
// 		}
// 	}
// 	return chain(plugins)
// }
