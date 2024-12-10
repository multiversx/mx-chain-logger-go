package logger

import (
	"testing"

	"github.com/multiversx/mx-chain-logger-go/proto"
	"github.com/stretchr/testify/require"
)

func TestJSONFormatter_Output(t *testing.T) {
	formatter := &JSONFormatter{}

	t.Run("with nil line", func(t *testing.T) {
		output := formatter.Output(nil)
		require.Nil(t, output)
	})

	t.Run("with line (without correlation)", func(t *testing.T) {
		ToggleCorrelation(false)

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
		require.Equal(t, `{"t":1122334455,"l":2,"n":"foo","m":"bar","a":["a","42","b","43"]}`+"\n", string(output))
	})

	t.Run("with line (with correlation)", func(t *testing.T) {
		ToggleCorrelation(true)

		line := &LogLineWrapper{
			LogLineMessage: proto.LogLineMessage{
				LoggerName: "foo",
				Message:    "bar",
				LogLevel:   int32(LogInfo),
				Args:       []string{"a", "42", "b", "43"},
				Timestamp:  1122334455,
				Correlation: proto.LogCorrelationMessage{
					Shard:    "3",
					Epoch:    42,
					Round:    4343,
					SubRound: "end",
				},
			},
		}

		output := formatter.Output(line)
		require.NotNil(t, output)
		require.Equal(t, `{"t":1122334455,"l":2,"n":"foo","s":"3","e":42,"r":4343,"sr":"end","m":"bar","a":["a","42","b","43"]}`+"\n", string(output))
	})

	t.Run("with bad line", func(t *testing.T) {
		ToggleCorrelation(false)

		line := &badLogLine{}

		output := formatter.Output(line)
		require.NotNil(t, output)
		require.Equal(t, `{"t":0,"l":0}`+"\n", string(output))
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
