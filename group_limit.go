package golimit

import "time"

type GroupLimiter struct {
	limiters []*Limiter
}

func NewGroupLimiter(limiters ...*Limiter) *GroupLimiter {
	return &GroupLimiter{limiters}
}

func (g *GroupLimiter) Allow(id string) bool {
	for _, l := range g.limiters {
		if !l.Allow(id) {
			return false
		}
	}

	return true
}

func (g *GroupLimiter) Ban(id string, d time.Duration) {
	for _, l := range g.limiters {
		l.Ban(id, d)
	}
}

func (g *GroupLimiter) IsBanned(id string) bool {
	for _, l := range g.limiters {
		if l.IsBanned(id) {
			return true
		}
	}

	return false
}
