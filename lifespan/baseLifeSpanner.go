package lifespan

type baseLifeSpanner struct {
	currentFile     string
	lifeSpanChannel chan string
}

func newBaseLifeSpanner() *baseLifeSpanner {
	return &baseLifeSpanner{
		lifeSpanChannel: make(chan string),
	}
}

// GetNotification gets the channel associated with a log recreate event
func (bls *baseLifeSpanner) GetNotification() <-chan string {
	return bls.lifeSpanChannel
}

// SetCurrentFile sets the current file for the logLifeSpanner
func (bls *baseLifeSpanner) SetCurrentFile(currentFile string) {
	bls.currentFile = currentFile
}

// Close closes the lifeSpanChannel
func (bls *baseLifeSpanner) Close() {
	close(bls.lifeSpanChannel)
}
