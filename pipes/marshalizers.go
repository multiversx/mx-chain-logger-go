package pipes

import (
	"encoding/json"
)

type jsonMarshalizer struct {
}

func (marshalizer *jsonMarshalizer) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (marshalizer *jsonMarshalizer) Unmarshal(obj interface{}, buff []byte) error {
	return json.Unmarshal(buff, obj)
}

func (marshalizer *jsonMarshalizer) IsInterfaceNil() bool {
	return marshalizer == nil
}
