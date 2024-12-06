package logger

import (
	"encoding/json"
	"fmt"
)

type JsonFormatter struct {
}

// Output converts the provided LogLineHandler into a slice of bytes ready for output
func (pf *JsonFormatter) Output(line LogLineHandler) []byte {
	if line == nil {
		return nil
	}

	output, err := json.Marshal(line)
	if err != nil {
		return []byte(fmt.Sprintf("error marshalling log line: %s", err.Error()))
	}

	return output
}

// IsInterfaceNil returns true if there is no value under the interface
func (pf *JsonFormatter) IsInterfaceNil() bool {
	return pf == nil
}
