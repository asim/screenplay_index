package app

import (
	"log"
	"sync"
	"time"
)

type trackingUpdate struct {
	ip string
	s  string
}

type trackingManager struct {
	mtx  sync.RWMutex
	seen map[string]int
	u    chan *trackingUpdate
}

var (
	slideInterval          = 3600
	defaultTrackingManager = newTrackingManager()
)

func init() {
	go defaultTrackingManager.run()
}

func newTrackingManager() *trackingManager {
	return &trackingManager{
		seen: make(map[string]int),
		u:    make(chan *trackingUpdate, 100),
	}
}

func (t *trackingManager) get(ip, s string) int {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.seen[ip+s]
}

func (t *trackingManager) run() {
	tick := time.NewTicker(time.Duration(slideInterval) * time.Second)

	for {
		select {
		case <-tick.C:
			log.Println("Sliding tracking")
			t.slide()
		case u := <-t.u:
			log.Println("Updating tracking")
			t.update(u.ip, u.s)
		}
	}
}

func (t *trackingManager) slide() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.seen = make(map[string]int)
}

func (t *trackingManager) track(ip, s string) {
	t.u <- &trackingUpdate{ip, s}
}

func (t *trackingManager) update(ip, s string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.seen[ip+s]++
}
