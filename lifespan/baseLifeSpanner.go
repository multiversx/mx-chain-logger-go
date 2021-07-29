package lifespan

type baseLifeSpanner struct {
	currentFile     string
	closeChannel    chan struct{}
	lifeSpanChannel chan string
}

func newBaseLifeSpanner() *baseLifeSpanner {
	return &baseLifeSpanner{
		closeChannel:    make(chan struct{}),
		lifeSpanChannel: make(chan string),
	}
}

// GetNotification - gets the channel associated with a log recreate event
func (bls *baseLifeSpanner) GetNotification() <-chan string {
	return bls.lifeSpanChannel
}

// Notify - notifies a change in the lifeSpan
func (bls *baseLifeSpanner) Notify(event string) {
	select {
	case bls.lifeSpanChannel <- event:
	case <-bls.closeChannel:
	}
}

// SetCurrentFile - sets the current file for the logLifeSpanner
func (bls *baseLifeSpanner) SetCurrentFile(currentFile string) {
	bls.currentFile = currentFile
}

// Close - closes the lifeSpanChannel
func (bls *baseLifeSpanner) Close() {
	close(bls.closeChannel)
}
