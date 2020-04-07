package pipes

import (
	"fmt"
	"testing"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/stretchr/testify/require"
)

func Test_ChildToParentThroughPipes(t *testing.T) {
	logger.ToggleLoggerName(true)

	// Parent setup
	parentPart, err := NewParentPart()
	require.Nil(t, err)
	parentPart.StartLoop() // StartLoop

	// Child setup
	profileReader, logsWriter := parentPart.GetChildPipes()
	childOutputSubject := logger.NewLogOutputSubject()
	childLogger := logger.NewLogger("child/foo", logger.LogTrace, childOutputSubject)
	childPart := NewChildPartWithLogOutputSubject(childOutputSubject, profileReader, logsWriter)

	//childPart.StartLoop()
	fmt.Println(childPart)

	// Child writes logs
	childLogger.Trace("test")
	childLogger.Trace("foobar")
	childLogger.Trace("foo", "answer", 42)

	time.Sleep(1 * time.Second)
}
