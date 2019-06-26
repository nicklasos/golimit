package golimit

import (
	"testing"
	"time"
)

func TestBan(t *testing.T) {
	lim := NewLimiter(1, 1)

	if lim.IsBanned("user") {
		t.Error("User should not be banned")
	}

	if lim.BannedCount() != 0 {
		t.Error("Banned count should be 0")
	}

	lim.Ban("user", time.Second*1)

	if !lim.IsBanned("user") {
		t.Error("User should be banned")
	}

	if lim.BannedCount() != 1 {
		t.Error("Banned count should be 1")
	}
}

func TestClean(t *testing.T) {
	lim := NewLimiter(1*time.Second, 2)

	lim.Allow("id")
	lim.Allow("id")

	lim.clean(time.Now())

	if len(lim.Visitors) != 1 {
		t.Error("Should not be cleaned")
	}

	time.Sleep(time.Second * 1)

	lim.clean(time.Now())

	if len(lim.Visitors) != 0 {
		t.Error("Should be cleaned")
	}
}

func TestAllow(t *testing.T) {
	lim := NewLimiter(1*time.Second, 2)

	if lim.Allow("id user") != true {
		t.Error("Should be allowed at first time")
	}

	if lim.Allow("id user") != true {
		t.Error("Should be allowed 2 messages in a second")
	}

	if lim.Allow("id user") != false {
		t.Error("Should be blocked on third message in a second")
	}

	time.Sleep(1 * time.Second)

	if lim.Allow("id user") != true {
		t.Error("Should be allowed after a second sleep")
	}
}
