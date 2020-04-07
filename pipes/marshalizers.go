package pipes

import (
	"encoding/json"
	"fmt"
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

type noopMarshalizer struct {
}

func (marshalizer *noopMarshalizer) Marshal(obj interface{}) ([]byte, error) {
	bytes, ok := obj.([]byte)
	if !ok {
		return nil, fmt.Errorf("obj is not []byte")
	}

	return bytes, nil
}

func (marshalizer *noopMarshalizer) Unmarshal(obj interface{}, buff []byte) error {
	return fmt.Errorf("Unmarshal not implemented")
}

func (marshalizer *noopMarshalizer) IsInterfaceNil() bool {
	return marshalizer == nil
}
