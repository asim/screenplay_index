package app

import (
	"encoding/json"
	"github.com/asim/screenplay_index/rankings"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type trendingUpdate struct {
	s  *script
	ip string
}

type trendingManager struct {
	trending []script
	seen     map[string]bool
	rm       *rankings.RankingsManager
	mtx      sync.RWMutex
	u        chan *trendingUpdate
}

var (
	trendingFile           = "trending.aof"
	numRanked              = 20
	rankWindow             = 48
	tickInterval           = 1800
	saveInterval           = 3600
	defaultTrendingManager = newTrendingManager()
)

func init() {
	defaultTrendingManager.init()
	go defaultTrendingManager.run()
}

func newTrendingManager() *trendingManager {
	return &trendingManager{
		seen: make(map[string]bool),
		u:    make(chan *trendingUpdate, 100),
		rm:   rankings.New(numRanked, rankWindow),
	}
}

func (t *trendingManager) click(s *script, ip string) {
	t.u <- &trendingUpdate{s, ip}
}

func (t *trendingManager) init() {
	dir, _ := filepath.Abs(workingDir)
	if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
		log.Println(err)
		panic(err)
	}

	t.load()
}

func (t *trendingManager) load() {
	dir, _ := filepath.Abs(workingDir)
	path := filepath.Join(dir, trendingFile)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}

	var trending []script
	err = json.Unmarshal(b, &trending)
	if err != nil {
		log.Println(err)
		return
	}

	t.trending = trending
}

func (t *trendingManager) getTrending() []script {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.trending
}

func (t *trendingManager) run() {
	tick := time.NewTicker(time.Duration(tickInterval) * time.Second)
	save := time.NewTicker(time.Duration(saveInterval) * time.Second)
	for {
		select {
		case <-tick.C:
			log.Println("Sliding trending")
			t.slide()
		case <-save.C:
			log.Println("Saving trending")
			t.save()
		case u := <-t.u:
			log.Println("Updating trending")
			t.update(u.s, u.ip)
		}
	}
}

func (t *trendingManager) save() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	b, err := json.Marshal(t.trending)
	if err != nil {
		log.Println(err)
		return
	}

	dir, _ := filepath.Abs(workingDir)
	path := filepath.Join(dir, trendingFile)
	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		log.Println(err)
	}
}

func (t *trendingManager) slide() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.updateTrending()
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
		t.updateTrending()
	}
}

func (t *trendingManager) updateTrending() {
	items := t.rm.GetRankings()

	// If we have enough items and the rankings aren't returning
	if required := numRanked / 2; len(*items) < required && len(t.trending) >= required {
		return
	}

	var trending []script

	for i := len(*items); i > 0; i-- {
		item := (*items)[i-1]
		s := item.Value().(*script)
		trending = append(trending, *s)
	}

	t.trending = trending
}
