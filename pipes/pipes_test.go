package pipes

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/stretchr/testify/require"
)

func Test_ChildToParentThroughPipes(t *testing.T) {
	logger.ToggleLoggerName(true)

	// Parent sets up the pipes
	readLogsFromChildFile, writeLogsToParentFile, err := os.Pipe()
	require.Nil(t, err)
	require.NotNil(t, readLogsFromChildFile)
	require.NotNil(t, writeLogsToParentFile)

	// Parent setup
	parentOutputSubject := logger.NewLogOutputSubject()
	parentOutputSubject.ClearObservers()
	parentOutputSubject.AddObserver(os.Stdout, &logger.ConsoleFormatter{})
	genericLoggerSink := logger.GetOrCreate("generic")
	parentForwarder := NewPipeObserverForwarder(readLogsFromChildFile, &jsonMarshalizer{}, genericLoggerSink)
	parentForwarder.StartFowarding()

	// Child setup
	pipeObserver := NewPipeObserver(writeLogsToParentFile)
	childOutputSubject := logger.NewLogOutputSubject()
	childOutputSubject.ClearObservers()
	logLineFormatter, _ := logger.NewLogLineWrapperFormatter(&jsonMarshalizer{})
	childOutputSubject.AddObserver(pipeObserver, logLineFormatter)
	childLogger := logger.NewLogger("child/foo", logger.LogTrace, childOutputSubject)

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
