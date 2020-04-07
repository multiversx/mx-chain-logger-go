package pipes

import (
	"encoding/binary"
	"io"
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

// Receiver intermediates communication (message receiving) via pipes
type Receiver struct {
	reader      *os.File
	marshalizer logger.Marshalizer
}

// NewReceiver creates a new receiver
func NewReceiver(reader *os.File, marshalizer logger.Marshalizer) *Receiver {
	return &Receiver{
		reader:      reader,
		marshalizer: marshalizer,
	}
}

// Receive receives a message, reads it from the pipe
func (receiver *Receiver) Receive(message interface{}) error {
	length, err := receiver.receiveMessageLength()
	if err != nil {
		return err
	}

	err = receiver.readMessage(message, length)
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Receiver) receiveMessageLength() (int, error) {
	buffer := make([]byte, sizeOfUint32)
	_, err := io.ReadFull(receiver.reader, buffer)
	if err != nil {
		return 0, err
	}

	length := binary.LittleEndian.Uint32(buffer)
	return int(length), nil
}

func (receiver *Receiver) readMessage(message interface{}, length int) error {
	buffer := make([]byte, length)
	_, err := io.ReadFull(receiver.reader, buffer)
	if err != nil {
		return err
	}

	err = receiver.marshalizer.Unmarshal(message, buffer)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown closes the pipe
func (receiver *Receiver) Shutdown() error {
	err := receiver.reader.Close()
	return err
}
