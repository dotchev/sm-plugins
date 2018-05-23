package json

import (
	"errors"

	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"

	"github.com/tidwall/gjson"
)

type JSON struct {
	gjson.Result
}

func Parse(data []byte) (*JSON, error) {
	if !gjson.ValidBytes(data) {
		return nil, errors.New("Invalid JSON")
	}
	return &JSON{gjson.ParseBytes(data)}, nil
}

func (json *JSON) Delete(path string) error {
	s, err := sjson.Delete(json.Raw, path)
	if err == nil {
		json.Result = gjson.Parse(s)
	}
	return err
}

func (json *JSON) Set(path string, value interface{}) error {
	s, err := sjson.Set(json.Raw, path, value)
	if err == nil {
		json.Result = gjson.Parse(s)
	}
	return err
}

func (json *JSON) Pretty() string {
	return string(pretty.Pretty([]byte(json.Raw)))
}
