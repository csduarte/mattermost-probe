package util

import (
	"sync"
	"time"
)

// MessageMap will syncronize multireads and writes of message times
type MessageMap struct {
	Items map[string]int64
	sync.RWMutex
}

// NewMessageMap returns a fresh message map
func NewMessageMap() *MessageMap {
	return &MessageMap{map[string]int64{}, sync.RWMutex{}}
}

// Add will insert an item into the message map
func (mm *MessageMap) Add(guid string, t int64) {
	mm.Lock()
	defer mm.Unlock()
	mm.Items[guid] = t
}

// Get will get the key and pass a bool back if it was found
func (mm *MessageMap) Get(key string) (int64, bool) {
	mm.RLock()
	defer mm.RUnlock()
	val, ok := mm.Items[key]
	return val, ok
}

// Delete reoves key if can be found, otherwise returns false
func (mm *MessageMap) Delete(key string) (int64, bool) {
	mm.Lock()
	defer mm.Unlock()
	start, ok := mm.Items[key]
	if ok {
		delete(mm.Items, key)
		return start, true
	}
	return 0, false
}

// FistOverdue will return if it finds any message beyond cutoff
func (mm *MessageMap) FistOverdue(cutoff int64) (string, int64) {
	mm.RLock()
	defer mm.RUnlock()
	now := time.Now().UnixNano() / 1000000
	for id, start := range mm.Items {
		delay := now - start
		if delay >= cutoff {
			return id, delay
		}
	}
	return "", 0
}

// Empty will reset the contents of the message map
func (mm *MessageMap) Empty() {
	mm.Lock()
	defer mm.Unlock()
	mm.Items = map[string]int64{}
}
