package sm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dotchev/sm-plugins/sm/plugin/json"
	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

// SendJSON writes a JSON value and sets the specified HTTP Status code
func SendJSON(writer http.ResponseWriter, code int, body *json.JSON) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	_, err := writer.Write([]byte(body.Pretty()))
	return err
}

// ReadJSONBody parse request body
func ReadJSONBody(request *http.Request) (*json.JSON, error) {
	contentType := request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("Invalid media type provided: %s", contentType)
	}

	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	body, err := json.Parse(bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode request body: %s", err)
	}
	return body, nil
}

func readRequest(request *http.Request) (*rest.Request, error) {
	pathParams := mux.Vars(request)

	queryParams := map[string]string{}
	for k, v := range request.URL.Query() {
		queryParams[k] = v[0]
	}

	body := &json.JSON{}
	if request.Method == "PUT" || request.Method == "POST" {
		var err error
		body, err = ReadJSONBody(request)
		if err != nil {
			return nil, err
		}
	}

	return &rest.Request{
		PathParams:  pathParams,
		QueryParams: queryParams,
		Body:        body,
	}, nil
}
