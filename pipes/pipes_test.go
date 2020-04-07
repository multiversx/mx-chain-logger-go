package pipes

import (
	"fmt"
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
	readProfileFromParentFile, writeProfileToChildFile, err := os.Pipe()
	require.Nil(t, err)
	require.NotNil(t, readLogsFromChildFile)
	require.NotNil(t, writeLogsToParentFile)

	// Parent setup
	// parentOutputSubject := logger.NewLogOutputSubject()
	// parentOutputSubject.ClearObservers()
	// parentOutputSubject.AddObserver(os.Stdout, &logger.ConsoleFormatter{})
	genericLoggerSink := logger.GetOrCreate("generic")
	parentForwarder := NewPipeObserverForwarder(readLogsFromChildFile, &jsonMarshalizer{}, genericLoggerSink)
	parentForwarder.StartFowarding()

	// Child setup
	childOutputSubject := logger.NewLogOutputSubject()
	childLogger := logger.NewLogger("child/foo", logger.LogTrace, childOutputSubject)
	childPart := NewChildPartWithLogOutputSubject(childOutputSubject, readProfileFromParentFile, writeLogsToParentFile)

	fmt.Println(childPart)
	fmt.Println(writeProfileToChildFile)

	// Child writes logs
	childLogger.Trace("test")
	childLogger.Trace("foobar")
	childLogger.Trace("foo", "answer", 42)

	time.Sleep(1 * time.Second)
}
