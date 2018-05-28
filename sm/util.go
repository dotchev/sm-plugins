package sm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"

	"github.com/gorilla/mux"

	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

// SendJSON writes a JSON value and sets the specified HTTP Status code
func SendJSON(writer http.ResponseWriter, code int, value interface{}) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

// ReadJSONBody parse request body
func ReadJSONBody(request *http.Request) (*simplejson.Json, error) {
	contentType := request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("Invalid media type provided: %s", contentType)
	}
	return simplejson.NewFromReader(request.Body)
}

func readRequest(request *http.Request) (*rest.Request, error) {
	pathParams := mux.Vars(request)

	queryParams := map[string]string{}
	for k, v := range request.URL.Query() {
		queryParams[k] = v[0]
	}

	var body *simplejson.Json
	if request.Method == "PUT" || request.Method == "POST" {
		var err error
		if body, err = ReadJSONBody(request); err != nil {
			return nil, err
		}
	}

	return &rest.Request{
		PathParams:  pathParams,
		QueryParams: queryParams,
		Body:        body,
	}, nil
}
