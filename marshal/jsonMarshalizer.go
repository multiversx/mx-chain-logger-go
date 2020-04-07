package marshal

import (
	"encoding/json"
)

// JSONMarshalizer -
type JSONMarshalizer struct {
}

// Marshal -
func (marshalizer *JSONMarshalizer) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// Unmarshal -
func (marshalizer *JSONMarshalizer) Unmarshal(obj interface{}, buff []byte) error {
	return json.Unmarshal(buff, obj)
}

// IsInterfaceNil -
func (marshalizer *JSONMarshalizer) IsInterfaceNil() bool {
	return marshalizer == nil
}
