package logger

import (
	"io"
	"os"
)

var _ io.Writer = (*pipeObserver)(nil)

type pipeObserver struct {
	writePipe *os.File
}

// NewPipeObserver -
func NewPipeObserver(writePipe *os.File) *pipeObserver {
	return &pipeObserver{
		writePipe: writePipe,
	}
}

// Write -
func (observer *pipeObserver) Write(p []byte) (n int, err error) {
	return observer.writePipe.Write(p)
}
