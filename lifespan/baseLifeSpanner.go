package lifespan

type baseLifeSpanner struct {
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

// Close closes the lifeSpanChannel
func (bls *baseLifeSpanner) Close() {
	close(bls.lifeSpanChannel)
}
