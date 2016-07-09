package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type pending struct {
	Script script
	State  string
}

type pendingManager struct {
	mtx     sync.RWMutex
	pending map[string]pending
	aol     io.ReadWriter
	enc     *json.Encoder
}

var (
	workingDir            = "./work"
	pendingFile           = "pending.aof"
	defaultPendingManager = newPendingManager()
)

func init() {
	defaultPendingManager.init()
}

func newPending(s script) pending {
	return pending{
		Script: s,
		State:  "pending",
	}
}

func newPendingManager() *pendingManager {
	return &pendingManager{
		pending: make(map[string]pending),
	}
}

func (p *pendingManager) add(s script) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if _, ok := p.pending[s.Url]; ok {
		return fmt.Errorf("Script is already waiting to be indexed")
	}

	np := newPending(s)
	p.pending[s.Url] = np

	err := p.enc.Encode(&np)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("There was a problem adding your script")
	}

	return nil
}

func (p *pendingManager) approve(id float64, url string) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if _, ok := p.pending[url]; !ok {
		return fmt.Errorf("Script not in pending queue")
	}

	np := p.pending[url]
	if id != np.Script.Id {
		return fmt.Errorf("Id provided %d does not match script id %d", id, np.Script.Id)
	}

	if np.State != "pending" {
		return fmt.Errorf("Script not pending approval. State %s", np.State)
	}

	np.State = "approved"
	delete(p.pending, url)

	if err := addScript(np.Script.Title, np.Script.Url); err != nil {
		log.Println(err)
		return fmt.Errorf("Error saving script", err)
	}

	if err := p.enc.Encode(&np); err != nil {
		log.Println(err)
		return fmt.Errorf("Error saving pending state", err)
	}

	return nil
}

func (p *pendingManager) init() {
	dir, _ := filepath.Abs(workingDir)
	path := filepath.Join(dir, pendingFile)

	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Println(err)
		panic(err)
	}

	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		p.aol = f
	} else {
		f, err := os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		p.aol = f
	}

	p.enc = json.NewEncoder(p.aol)
	p.load()
}

func (p *pendingManager) load() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	var pp pending
	d := json.NewDecoder(p.aol)
	for {
		err := d.Decode(&pp)
		if err != nil {
			break
		}
		p.pending[pp.Script.Url] = pp
		if pp.State == "approved" || pp.State == "rejected" {
			delete(p.pending, pp.Script.Url)
		}
	}
}

func (p *pendingManager) read() []script {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	var scripts []script
	for _, s := range p.pending {
		scripts = append(scripts, s.Script)
	}

	return scripts
}

func (p *pendingManager) reject(id float64, url string) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if _, ok := p.pending[url]; !ok {
		return fmt.Errorf("Script not in pending queue")
	}

	np := p.pending[url]
	if id != np.Script.Id {
		return fmt.Errorf("Id provided %d does not match script id %d", id, np.Script.Id)
	}

	if np.State != "pending" {
		return fmt.Errorf("Script not pending approval. State %s", np.State)
	}

	np.State = "rejected"
	delete(p.pending, url)

	if err := p.enc.Encode(&np); err != nil {
		log.Println(err)
		return fmt.Errorf("Error saving pending state", err)
	}

	return nil
}
