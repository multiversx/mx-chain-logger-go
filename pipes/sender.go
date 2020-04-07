package pipes

import (
	"encoding/binary"
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

// Sender intermediates communication (message sending) via pipes
type Sender struct {
	writer      *os.File
	marshalizer logger.Marshalizer
}

// NewSender creates a new sender
func NewSender(writer *os.File, marshalizer logger.Marshalizer) *Sender {
	return &Sender{
		writer:      writer,
		marshalizer: marshalizer,
	}
}

// Send sends a message over the pipe
func (sender *Sender) Send(message interface{}) (int, error) {
	dataBytes, err := sender.marshalizer.Marshal(message)
	if err != nil {
		return 0, err
	}

	length := len(dataBytes)
	err = sender.sendMessageLength(length)
	if err != nil {
		return 0, err
	}

	_, err = sender.writer.Write(dataBytes)
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
