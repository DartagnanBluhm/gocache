package gocache

import (
	"time"
)

// ExpirableData Object defines a expirable data object, these are the only objects removed during the Cleaning process
type ExpirableData struct {
	Entry      interface{} `json:"entry"`
	Expiration time.Time   `json:"expiration"`
}

// DataEntry Object defines a Cache entry. This object is used to transmit Cache entries over channels in one transaction
type DataEntry struct {
	Key   interface{} `json:"key"`
	Value interface{} `json:"value"`
	Valid bool        `json:"valid"`
}

// Cache Object defines the channels, Cleaner and map of a Cache
type Cache struct {
	add     chan chan interface{}
	ask     chan chan interface{}
	keys    chan chan interface{}
	remove  chan chan interface{}
	cache   map[interface{}]interface{}
	Cleaner *Cleaner
	reset   chan chan interface{}
	quit    chan chan interface{}
}

// NewCache initalises a new Cache
func NewCache() *Cache {
	c := Cache{
		add:    make(chan chan interface{}),
		ask:    make(chan chan interface{}),
		keys:   make(chan chan interface{}),
		remove: make(chan chan interface{}),
		cache:  make(map[interface{}]interface{}),
		reset:  make(chan chan interface{}),
		quit:   make(chan chan interface{}),
	}
	go c.Start()
	return &c
}

// AddCleaner initalises a new Cleaner and adds it to the Cache
func (c *Cache) AddCleaner(cleanInterval time.Duration) {
	c.Cleaner = NewCleaner(cleanInterval)
	c.Cleaner.cache = c
}

// ChanSend opens a connection on the provided channel, sending the provided object, then waits to receive a response, returning the response
func ChanSend(dest chan chan interface{}, value interface{}) interface{} {
	ch := make(chan interface{})
	dest <- ch
	ch <- value
	response := <-ch
	close(ch)
	return response
}

// ChanReceive opens a connection on the provided channel and waits to receive a response, returning the response
func ChanReceive(dest chan chan interface{}) interface{} {
	ch := make(chan interface{})
	dest <- ch
	response := <-ch
	close(ch)
	return response
}

// Add a new entry to the Cache
func (c *Cache) Add(key interface{}, value interface{}) {
	ChanSend(c.add, DataEntry{
		Key:   key,
		Value: value,
	})
}

// Remove an entry from the Cache
func (c *Cache) Remove(key interface{}) {
	ChanSend(c.remove, key)
}

// Get a DataEntry Object from the Cache
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	response := ChanSend(c.ask, key).(DataEntry)
	return response.Value, response.Valid
}

// Reset removes all entries in the Cache
func (c *Cache) Reset() {
	ChanReceive(c.reset)
}

// Stop halts the Cache
func (c *Cache) Stop() {
	ChanReceive(c.quit)
}

// GetAll extracts all keys from the Cache
func (c *Cache) GetAll() []interface{} {
	return ChanReceive(c.keys).([]interface{})
}

// Start the cache, listening on its channels for requests and starts its cleaner
func (c *Cache) Start() {
	for {
		select {
		case ch := <-c.add:
			entry := (<-ch).(DataEntry)
			c.cache[entry.Key] = entry.Value
			ch <- entry
		case ch := <-c.remove:
			key := (<-ch).(string)
			delete(c.cache, key)
			ch <- key
		case ch := <-c.ask:
			key := (<-ch).(string)
			entry := DataEntry{}
			entry.Value, entry.Valid = c.cache[key]
			ch <- entry
		case ch := <-c.keys:
			keys := make([]interface{}, len(c.cache))
			i := 0
			for k := range c.cache {
				keys[i] = k
				i++
			}
			ch <- keys
		case ch := <-c.reset:
			c.cache = make(map[interface{}]interface{})
			ch <- true
		case ch := <-c.quit:
			close(c.add)
			close(c.remove)
			close(c.ask)
			close(c.keys)
			close(c.reset)
			close(c.quit)
			c.cache = nil // detach map so it is garbage collected
			ch <- true
			return
		}
	}
}
