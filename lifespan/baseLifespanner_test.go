package lifespan

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseLifeSpanner_GetChannelShouldWork(t *testing.T) {
	t.Parallel()

	bls := newBaseLifeSpanner()
	open := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_, open = <-bls.GetNotification()
		wg.Done()
	}()

	bls.lifeSpanChannel <- ""

	wg.Wait()

	assert.True(t, open)
}

func TestBaseLifeSpanner_CloseShouldCloseChannel(t *testing.T) {
	t.Parallel()

	bls := newBaseLifeSpanner()
	open := true
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_, open = <-bls.GetNotification()
		wg.Done()
	}()

	bls.Close()

	wg.Wait()

	assert.False(t, open)
}
