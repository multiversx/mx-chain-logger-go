package logger

import (
	"testing"

	"github.com/multiversx/mx-chain-logger-go/proto"
	"github.com/stretchr/testify/require"
)

func TestJSONFormatter_Output(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{}

	t.Run("with nil line", func(t *testing.T) {
		t.Parallel()

		output := formatter.Output(nil)
		require.Nil(t, output)
	})

	t.Run("with line", func(t *testing.T) {
		t.Parallel()

		line := &LogLineWrapper{
			LogLineMessage: proto.LogLineMessage{
				LoggerName: "foo",
				Message:    "bar",
				LogLevel:   int32(LogInfo),
				Args:       []string{"a", "42", "b", "43"},
				Timestamp:  1122334455,
			},
		}

		output := formatter.Output(line)
		require.NotNil(t, output)
		require.Equal(t, `{"Message":"bar","LogLevel":2,"Args":["a","42","b","43"],"Timestamp":1122334455,"LoggerName":"foo","Correlation":{}}`+"\n", string(output))
	})

	t.Run("with bad line", func(t *testing.T) {
		t.Parallel()

		line := &badLogLine{}

		output := formatter.Output(line)
		require.NotNil(t, output)
		require.Equal(t, "error marshalling log line: json: unsupported type: func()\n", string(output))
	})
}

type badLogLine struct {
	proto.LogLineMessage
	NotMarshalizable func() `json:"notMarshalizable"`
}

// IsInterfaceNil -
func (line *badLogLine) IsInterfaceNil() bool {
	return false
}
