package lifespan

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseLifeSpanner_GetNotificationShouldWork(t *testing.T) {
	t.Parallel()

	bls := newBaseLifeSpanner()
	open := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_, open = <-bls.GetNotification()
		wg.Done()
	}()

	bls.Notify("")

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
		_, open = <-bls.closeChannel
		wg.Done()
	}()

	bls.Close()

	wg.Wait()

	assert.False(t, open)
}

func TestBaseLifeSpanner_NotifyShouldWriteOnChannel(t *testing.T) {
	t.Parallel()

	bls := newBaseLifeSpanner()
	open := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	notificationMessage := "notification"
	receivedNotification := ""
	go func() {
		receivedNotification, open = <-bls.GetNotification()
		wg.Done()
	}()

	bls.Notify(notificationMessage)

	wg.Wait()

	assert.True(t, open)
	assert.Equal(t, notificationMessage, receivedNotification)
}

func TestBaseLifeSpanner_NotifyIfClosedShouldNotPanic(t *testing.T) {
	t.Parallel()

	bls := newBaseLifeSpanner()
	open := true
	wg := sync.WaitGroup{}
	wg.Add(1)
	notificationMessage := "notification"
	go func() {
		select {
		case <-bls.GetNotification():
		case <-bls.closeChannel:
			open = false
		}
		wg.Done()
	}()

	bls.Close()

	bls.Notify(notificationMessage)

	wg.Wait()

	assert.False(t, open)
}
