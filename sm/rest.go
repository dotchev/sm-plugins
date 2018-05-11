package sm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dotchev/sm-plugins/sm/plugin/osb"
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

func readOSBRequest(request *http.Request) (*osb.Request, error) {
	params := mux.Vars(request)             // get path parameters
	for k, v := range request.URL.Query() { // get query parameters
		params[k] = v[0]
	}

	var body interface{}
	if request.Method == "PUT" || request.Method == "POST" {
		if err := ReadJSONBody(request, &body); err != nil {
			return nil, err
		}
	}
	r := &osb.Request{Params: params, Body: body}
	return r, nil
}
