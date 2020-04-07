package pipes

import (
	"encoding/json"
	"testing"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/stretchr/testify/require"
)

func Test_ChildToParentThroughPipes(t *testing.T) {
	logger.ToggleLoggerName(true)

	logLineMarshalizer := &jsonMarshalizer{}

	// Parent setup
	parentPart, err := NewParentPart(logLineMarshalizer)
	require.Nil(t, err)
	parentPart.StartLoop()

	// Child setup
	profileReader, logsWriter := parentPart.GetChildPipes()
	childOutputSubject := logger.NewLogOutputSubject()
	childLogger := logger.NewLogger("child/foo", logger.LogTrace, childOutputSubject)
	childPart := NewChildPartWithLogOutputSubject(childOutputSubject, profileReader, logsWriter, logLineMarshalizer)
	childPart.StartLoop()

	// Child writes logs
	childLogger.Trace("test")
	childLogger.Trace("foobar")
	childLogger.Trace("foo", "answer", 42)

	time.Sleep(1 * time.Second)
}

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
