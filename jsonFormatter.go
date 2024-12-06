package logger

import (
	"encoding/json"
	"fmt"
)

type JSONFormatter struct {
}

// Output converts the provided LogLineHandler into a slice of bytes ready for output
func (pf *JSONFormatter) Output(line LogLineHandler) []byte {
	if line == nil {
		return nil
	}

	output, err := json.Marshal(line)
	if err != nil {
		return []byte(fmt.Sprintf("error marshalling log line: %s", err.Error()))
	}

	output = append(output, '\n')
	return output
}

// IsInterfaceNil returns true if there is no value under the interface
func (pf *JSONFormatter) IsInterfaceNil() bool {
	return pf == nil
}
