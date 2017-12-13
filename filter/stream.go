package filter

import (
	"fmt"

	"github.com/json-iterator/go"
)

// Stream is handed over to every filter implementation
type Stream struct {
	// JSONContent is the original log entry in json
	JSONContent string
	// MapContent is a map representation of the JsonContent
	MapContent map[string]interface{}
	// IndexName the name if the ElasticSearch Index to write to
	IndexName string
}

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Unmarshal the json to a map
func (fs *Stream) Unmarshal() error {
	data := make(map[string]interface{})
	err := json.UnmarshalFromString(fs.JSONContent, &data)
	if err != nil {
		return fmt.Errorf("cannot decode data:%v", err)
	}
	fs.MapContent = data
	return nil
}

// Marshal the map back to json
func (fs *Stream) Marshal() error {
	result, err := json.MarshalToString(&fs.MapContent)
	if err != nil {
		return fmt.Errorf("cannot encode data:%v", err)
	}
	fs.JSONContent = result
	return nil
}
