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

func TestParentPart_ReceivesLogsFromChildProcess(t *testing.T) {
	// Record logs by means of a logs gatherer, so we can apply assertions afterwards
	gatherer := &dummyLogsGatherer{}
	logOutputSubject := logger.GetLogOutputSubject()
	logOutputSubject.AddObserver(gatherer, gatherer)

	part, err := NewParentPart("child-name", &marshal.JSONMarshalizer{})
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

	part.StartLoop(childStdout, arwenStderr)

	// TODO: Wait after a certain message
	time.Sleep(1 * time.Second)
	require.True(t, gatherer.ContainsLogLine("foo", logger.LogInfo, "foo-info"))
	require.True(t, gatherer.ContainsLogLine("bar", logger.LogInfo, "bar-info"))
	require.False(t, gatherer.ContainsText("foo-trace-no"))
	require.False(t, gatherer.ContainsText("bar-trace-no"))
	require.True(t, gatherer.ContainsLogLine("foo", logger.LogInfo, "foo-in-go"))
	require.True(t, gatherer.ContainsLogLine("bar", logger.LogInfo, "bar-in-go"))

	// Change logs profile
	logger.ToggleLoggerName(true)
	logger.SetLogLevel("*:TRACE")
	logger.NotifyProfileChange()

	// TODO: Wait after a certain message
	time.Sleep(2 * time.Second)
	require.True(t, gatherer.ContainsLogLine("foo", logger.LogTrace, "foo-trace-yes"))
	require.True(t, gatherer.ContainsLogLine("bar", logger.LogTrace, "bar-trace-yes"))
	require.True(t, gatherer.ContainsLogLine(textOutputSinkName, logger.LogTrace, "child-name"))
	require.True(t, gatherer.ContainsLogLine(textOutputSinkName, logger.LogError, "child-name"))
	require.True(t, gatherer.ContainsText("Here's some stderr"))
	require.True(t, gatherer.ContainsText("Here's some stdout"))
}
