package golimit

import (
	"sync"
	"time"
)

type Visitors = map[string]*Visitor
type Bans = map[string]time.Time

type Visitor struct {
	id       string
	messages []time.Time
}

type Limiter struct {
	Period   time.Duration
	Limit    int
	Visitors Visitors
	mtx      sync.Mutex
	Bans     Bans
}

func NewLimiter(period time.Duration, limit int) *Limiter {
	lim := &Limiter{
		Period:   period,
		Limit:    limit,
		Visitors: make(Visitors),
		Bans:     make(Bans),
	}

	go lim.clenup()

	return lim
}

func (l *Limiter) visitor(id string) *Visitor {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	v, exists := l.Visitors[id]
	if !exists {
		l.Visitors[id] = &Visitor{id: id}

		return l.Visitors[id]
	}

	return v
}

func (l *Limiter) clenup() {
	for {
		time.Sleep(1 * time.Minute)

		l.clean(time.Now())
	}
}

func (l *Limiter) clean(t time.Time) {
	max := t.Add(-l.Period)
	l.mtx.Lock()

	for key, v := range l.Visitors {
		last := v.messages[len(v.messages)-1]
		if max.After(last) {
			delete(l.Visitors, key)
		}
	}

	l.mtx.Unlock()
}

func (l *Limiter) unban() {
	for {
		time.Sleep(time.Minute * 1)
		l.mtx.Lock()
		now := time.Now()
		for id, banTime := range l.Bans {
			if now.After(banTime) {
				delete(l.Bans, id)
			}
		}
		l.mtx.Unlock()
	}
}

func (l *Limiter) Allow(id string) bool {
	v := l.visitor(id)

	max := time.Now().Add(-l.Period)

	l.mtx.Lock()
	defer l.mtx.Unlock()

	// Remove old messages
	for _, msg := range v.messages {
		if max.After(msg) {
			v.messages = v.messages[1:]
		} else {
			break
		}
	}

	// Will the followin message exceed the limit?
	if len(v.messages)+1 > l.Limit {
		return false
	}

	v.messages = append(v.messages, time.Now())

	return true
}

func (l *Limiter) Ban(id string, d time.Duration) {
	l.mtx.Lock()
	l.Bans[id] = time.Now().Add(d)
	l.mtx.Unlock()
}

func (l *Limiter) IsBanned(id string) bool {
	l.mtx.Lock()
	_, exists := l.Bans[id]
	l.mtx.Unlock()

	return exists
}
