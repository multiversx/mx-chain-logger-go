package logger

import (
	"encoding/json"
	"fmt"
)

type JSONFormatter struct {
}

type jsonOptimizedLogLine struct {
	Timestamp  int64    `json:"t"`
	LogLevel   int32    `json:"l"`
	LoggerName string   `json:"n,omitempty"`
	Message    string   `json:"m,omitempty"`
	Args       []string `json:"a,omitempty"`
}

type jsonOptimizedLogLineWithCorrelation struct {
	Timestamp  int64    `json:"t"`
	LogLevel   int32    `json:"l"`
	LoggerName string   `json:"n,omitempty"`
	Shard      string   `json:"s"`
	Epoch      uint32   `json:"e"`
	Round      int64    `json:"r"`
	SubRound   string   `json:"sr"`
	Message    string   `json:"m,omitempty"`
	Args       []string `json:"a,omitempty"`
}

// Output converts the provided LogLineHandler into a slice of bytes ready for output
func (pf *JSONFormatter) Output(line LogLineHandler) []byte {
	if line == nil {
		return nil
	}

	optimizedLine := getOptimizedLogLine(line)
	output, err := json.Marshal(optimizedLine)
	if err != nil {
		// Generally speaking, this is dead code, since the structures are well-defined and can always be marshalled.
		errorMessage := fmt.Sprintf("error marshalling log line: %s, %s: %s\n", line.GetLoggerName(), line.GetMessage(), err.Error())
		return []byte(errorMessage)
	}

	output = append(output, '\n')
	return output
}

func getOptimizedLogLine(line LogLineHandler) interface{} {
	if IsEnabledCorrelation() {
		return &jsonOptimizedLogLineWithCorrelation{
			Timestamp:  line.GetTimestamp(),
			LogLevel:   line.GetLogLevel(),
			LoggerName: line.GetLoggerName(),
			Shard:      line.GetCorrelation().Shard,
			Epoch:      line.GetCorrelation().Epoch,
			Round:      line.GetCorrelation().Round,
			SubRound:   line.GetCorrelation().SubRound,
			Message:    line.GetMessage(),
			Args:       line.GetArgs(),
		}
	}

	return &jsonOptimizedLogLine{
		Timestamp:  line.GetTimestamp(),
		LogLevel:   line.GetLogLevel(),
		LoggerName: line.GetLoggerName(),
		Message:    line.GetMessage(),
		Args:       line.GetArgs(),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (pf *JSONFormatter) IsInterfaceNil() bool {
	return pf == nil
}
