package app

import (
	"github.com/asim/screenplay_index/rankings"
	"sync"
	"time"
)

type trendingUpdate struct {
	s  *script
	ip string
}

type trendingManager struct {
	trending []*script
	seen     map[string]bool
	rm       *rankings.RankingsManager
	mtx      sync.RWMutex
	u        chan *trendingUpdate
}

var (
	numRanked              = 20
	rankWindow             = 48
	tickInterval           = 1800
	defaultTrendingManager = newTrendingManager()
)

func init() {
	go defaultTrendingManager.run()
}

func newTrendingManager() *trendingManager {
	return &trendingManager{
		seen: make(map[string]bool),
		u:    make(chan *trendingUpdate, 100),
		rm:   rankings.New(numRanked, rankWindow),
	}
}

func (t *trendingManager) getTrending() []*script {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.trending
}

func (t *trendingManager) slide() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	items := t.rm.GetRankings()
	var trending []*script

	for _, item := range *items {
		trending = append(trending, item.Value().(*script))
	}

	t.trending = trending
	t.seen = make(map[string]bool)
	t.rm.Slide()
}

func (t *trendingManager) update(s *script, ip string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if !t.seen[s.Short+ip] {
		t.rm.Update(s)
		t.seen[s.Short+ip] = true
	}

	if len(t.trending) < numRanked {
		items := t.rm.GetRankings()
		var trending []*script

		for _, item := range *items {
			trending = append(trending, item.Value().(*script))
		}

		t.trending = trending
	}
}

func (t *trendingManager) click(s *script, ip string) {
	t.u <- &trendingUpdate{s, ip}
}

func (t *trendingManager) run() {
	tick := time.NewTicker(time.Duration(tickInterval) * time.Second)

	for {
		select {
		case <-tick.C:
			t.slide()
		case u := <-t.u:
			t.update(u.s, u.ip)
		}
	}
}
