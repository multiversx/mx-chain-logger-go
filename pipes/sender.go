package pipes

import (
	"encoding/binary"
	"os"
)

// Sender intermediates communication (message sending) via pipes
type Sender struct {
	writer *os.File
}

// NewSender creates a new sender
func NewSender(writer *os.File) *Sender {
	return &Sender{
		writer: writer,
	}
}

// Send sends a message over the pipe
func (sender *Sender) Send(message []byte) (int, error) {
	length := len(message)
	err := sender.sendMessageLength(length)
	if err != nil {
		return 0, err
	}

	_, err = sender.writer.Write(message)
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (sender *Sender) sendMessageLength(length int) error {
	buffer := make([]byte, sizeOfUint32)
	binary.LittleEndian.PutUint32(buffer, uint32(length))
	_, err := sender.writer.Write(buffer)
	return err
}

// Shutdown closes the pipe
func (sender *Sender) Shutdown() error {
	err := sender.writer.Close()
	return err
}
