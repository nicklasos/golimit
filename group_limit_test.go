package golimit

import (
	"testing"
	"time"
)

func TestAllowGroup(t *testing.T) {
	lim := NewGroupLimiter(
		NewLimiter(1*time.Second, 2),
		NewLimiter(1*time.Minute, 30),
	)

	if lim.limiters[1].Limit != 30 {
		t.Error("Group limiters is not set")
	}

	lim.Allow("id")

	if _, ok := lim.limiters[0].Visitors["id"]; !ok {
		t.Error("Group limiter is not proxies Allow to Limiter #1")
	}

	if _, ok := lim.limiters[1].Visitors["id"]; !ok {
		t.Error("Group limiter is not proxies Allow to Limiter #2")
	}
}

func TestAllowGroupBan(t *testing.T) {
	lim := NewGroupLimiter(
		NewLimiter(1*time.Second, 2),
		NewLimiter(1*time.Minute, 30),
	)

	lim.Ban("id", 1*time.Minute)

	if !lim.limiters[0].IsBanned("id") {
		t.Error("Group limiter ban is not working")
	}

	if !lim.limiters[1].IsBanned("id") {
		t.Error("Group limiter ban is not working")
	}

	if !lim.IsBanned("id") {
		t.Error("Group limiter ban is not working")
	}
}
