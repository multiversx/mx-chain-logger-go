package pipes

import (
	"os"
	"os/exec"
	"testing"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/marshal"
	"github.com/stretchr/testify/require"
)

func Test_ChildPartLogsToParentPart(t *testing.T) {
	logger.ToggleLoggerName(true)

	logLineMarshalizer := &marshal.JSONMarshalizer{}

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

func Test_ChilProcessLogsToParentProcess(t *testing.T) {
	logLineMarshalizer := &marshal.JSONMarshalizer{}

	part, err := NewParentPart(logLineMarshalizer)
	require.Nil(t, err)
	profileReader, logsWriter := part.GetChildPipes()

	command := exec.Command("./testchild")
	command.ExtraFiles = []*os.File{profileReader, logsWriter}

	childStdout, err := command.StdoutPipe()
	require.Nil(t, err)
	arwenStderr, err := command.StderrPipe()
	require.Nil(t, err)

	err = command.Start()
	require.Nil(t, err)

	part.StartLoop()
	part.ContinuouslyReadTextualOutput(childStdout, arwenStderr, "child-tag")

	time.Sleep(1 * time.Second)

	logger.ToggleLoggerName(true)
	logger.SetLogLevel("*:TRACE")

	logger.NotifyProfileChange()

	time.Sleep(1 * time.Second)
}
