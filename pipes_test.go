package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ChildToParentThroughPipes(t *testing.T) {
	// Parent sets up the pipes
	readLogsFromChildFile, writeLogsToParentFile, err := os.Pipe()
	require.Nil(t, err)

	// Child setup
	pipeObserver := NewPipeObserver(writeLogsToParentFile)
}
