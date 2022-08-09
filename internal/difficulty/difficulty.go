package difficulty

import (
	"sync"
)

const diffMaxLimit = 100

type Difficulty struct {
	mu         sync.RWMutex
	difficulty int
	maxConn    int
	curConn    int
}

func New(maxConn int) *Difficulty {
	return &Difficulty{maxConn: maxConn}
}

func (d *Difficulty) NewConn() {
	d.mu.Lock()
	d.curConn++
	if d.curConn > d.maxConn {
		d.difficulty++
	}
	if d.difficulty > diffMaxLimit {
		d.difficulty = diffMaxLimit
	}
	d.mu.Unlock()
}

func (d *Difficulty) ConnOver() {
	d.mu.Lock()
	d.curConn--
	if d.curConn < d.maxConn {
		d.difficulty--
	}
	if d.difficulty < 1 {
		d.difficulty = 1
	}
	d.mu.Unlock()
}

func (d *Difficulty) Difficulty() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.difficulty
}
