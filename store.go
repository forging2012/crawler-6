package crawler

import (
	"net/url"
	"sync"
	"time"
)

type status int

const (
	U_Init     status = iota
	U_Waiting         // in waiting queue
	U_Enqueued        // in main queue
	U_Sieving         // in filter
	U_Finished
	U_Redirected
	U_Error
)

type storeHandle interface {
	V() *URL
	Unlock()
}

type URLStore interface {
	Get(u url.URL) (URL, bool)
	Put(u URL)
	// Watch locks the entry located by u and returns a pointer to data.
	Watch(u url.URL) storeHandle
	// WatchP locks the entry(if not exist, create)
	WatchP(u URL) storeHandle
}

type URL struct {
	Loc     url.URL
	Score   int64
	Freq    time.Duration
	Visited struct {
		Count int
		Time  time.Time
	}
	Depth        int
	LastModified time.Time
	Status       status
	nextTime     time.Time
}

type entry struct {
	url URL
	sync.Mutex
}

func (entry *entry) V() *URL {
	return &entry.url
}

type store struct {
	sync.RWMutex
	m map[url.URL]*entry
}

func newURL(u url.URL) *URL {
	u.Fragment = ""
	return &URL{
		Loc:    u,
		status: U_Init,
	}
}

func newMemStore() *store {
	return &store{
		m: make(map[url.URL]*entry),
	}
}

func (p *store) Watch(u url.URL) (h storeHandle) {
	p.RLock()
	defer p.RUnlock()
	entry, ok := p.m[u]
	if !ok {
		return
	}
	entry.Lock()
	h = entry
	return
}

func (p *store) WatchP(u URL) storeHandle {
	p.Lock()
	defer p.Unlock()
	u.Loc.Fragment = ""
	entry, ok := p.m[u.Loc]
	if ok {
		entry.Lock()
		return entry
	}

	entry = &entry{url: u}
	entry.Lock()
	p.m[u.Loc] = entry
	return entry
}

func (p *store) Put(u URL) {
	u.Loc.Fragment = ""
	p.Lock()
	p.m[u.Loc] = &entry{url: u}
	p.Unlock()
}

func (p *store) Get(u url.URL) (uu URL, ok bool) {
	p.RLock()
	entry, present := p.m[u]
	if present {
		uu, ok = entry.url, true
	}
	p.RUnlock()
	return
}
