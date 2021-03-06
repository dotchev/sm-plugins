package sm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

// SendJSON writes a JSON value and sets the specified HTTP Status code
func SendJSON(writer http.ResponseWriter, code int, value interface{}) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

// ReadJSONBody parse request body
func ReadJSONBody(request *http.Request, value interface{}) error {
	contentType := request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return fmt.Errorf("Invalid media type provided: %s", contentType)
	}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("Failed to decode request body: %s", err)
	}
	return nil
}

func readRequest(request *http.Request) (*rest.Request, error) {
	pathParams := mux.Vars(request)

	queryParams := map[string]string{}
	for k, v := range request.URL.Query() {
		queryParams[k] = v[0]
	}

	var body interface{}
	if request.Method == "PUT" || request.Method == "POST" {
		if err := ReadJSONBody(request, &body); err != nil {
			return nil, err
		}
	}

	return &rest.Request{
		Request:     request,
		PathParams:  pathParams,
		QueryParams: queryParams,
		Body:        body,
	}, nil
}
