package pipes

// Messenger intermediates communication (message exchange) via pipes
type Messenger struct {
	receiver *Receiver
	sender   *Sender
}

// NewMessenger creates a new messenger
func NewMessenger(receiver *Receiver, sender *Sender) *Messenger {
	return &Messenger{
		receiver: receiver,
		sender:   sender,
	}
}

// Send sends a message over the pipe
func (messenger *Messenger) Send(message []byte) (int, error) {
	return messenger.sender.Send(message)
}

// Receive receives a message, reads it from the pipe
func (messenger *Messenger) Receive() ([]byte, error) {
	return messenger.receiver.Receive()
}

// Shutdown closes the pipes
func (messenger *Messenger) Shutdown() {
	_ = messenger.receiver.Shutdown()
	_ = messenger.sender.Shutdown()
}
