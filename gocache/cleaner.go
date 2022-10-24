package gocache

import (
	"time"
)

// Cleaner Object defines the interval a Cleaner should clean the Cache
type Cleaner struct {
	interval time.Duration
	cache    *Cache
	quit     chan chan interface{}
}

// NewCleaner creates a new Cleaner Object
func NewCleaner(i time.Duration) *Cleaner {
	cleaner := Cleaner{
		interval: i,
		quit:     make(chan chan interface{}),
	}
	go cleaner.Start()
	return &cleaner
}

// Clean removes all expired objects from the Cache
func (c *Cleaner) Clean() {
	keys := c.cache.GetAll()
	for _, key := range keys {
		value, ok := c.cache.Get(key)
		if !ok {
			continue
		}
		switch value := value.(type) {
		case ExpirableData:
			if time.Now().UnixNano() > value.Expiration.UnixNano() {
				c.cache.Remove(key)
			}
		}
	}
}

// Stop halts the Cleaner
func (c *Cleaner) Stop() {
	ChanReceive(c.quit)
}

// Run starts the Cleaner, periodically cleaning the Cache of expired data
func (c *Cleaner) Start() {
	ticker := time.NewTicker(c.interval)
	for {
		select {
		case <-ticker.C:
			c.Clean()
		case ch := <-c.quit:
			close(c.quit)
			ch <- true
			return
		}
	}
}
