package logger

import (
	"sync"
)

type LifeSpanWrapper struct {
	mutLogLifeSpanner sync.Mutex
	lifeSpanner       LogLifeSpanner
}

// SetLifeSpanner sets the current lifeSpanner
func (lsw *LifeSpanWrapper) SetLifeSpanner(spanner LogLifeSpanner) {
	lsw.mutLogLifeSpanner.Lock()
	lsw.lifeSpanner = spanner
	lsw.mutLogLifeSpanner.Unlock()
}

// SetCurrentFile sets the new logger file
func (lsw *LifeSpanWrapper) SetCurrentFile(newFile string) {
	lsw.mutLogLifeSpanner.Lock()
	lsw.lifeSpanner.SetCurrentFile(newFile)
	lsw.mutLogLifeSpanner.Unlock()
}

// GetNotificationChannel gets the notification channel for the log lifespan
func (lsw *LifeSpanWrapper) GetNotificationChannel() <-chan string {
	lsw.mutLogLifeSpanner.Lock()
	defer lsw.mutLogLifeSpanner.Unlock()
	return lsw.lifeSpanner.GetNotification()
}
