package model

import (
	"slices"
	"sync"
)

type Tracker struct {
	id string

	activesMu sync.RWMutex
	actives   []*GrantInfo

	expiredMu sync.RWMutex
	expired   []GrantInfo
}

func NewTracker(id string) *Tracker {
	return &Tracker{
		id: id,
	}
}

// ID returns the ID of the player who owns the tracker.
func (t *Tracker) ID() string {
	return t.id
}

// Actives returns the active grants of the player.
func (t *Tracker) Actives() []*GrantInfo {
	t.activesMu.RLock()
	defer t.activesMu.RUnlock()

	return t.actives
}

// Expired returns the expired grants of the player.
func (t *Tracker) Expired() []GrantInfo {
	t.expiredMu.RLock()
	defer t.expiredMu.RUnlock()

	return t.expired
}

// AddActive adds a grant to the active grants of the player.
func (t *Tracker) AddActive(gi *GrantInfo) {
	t.activesMu.Lock()
	t.actives = append(t.actives, gi)
	t.activesMu.Unlock()
}

// AddExpired adds a grant to the expired grants of the player.
func (t *Tracker) AddExpired(gi GrantInfo) {
	t.expiredMu.Lock()
	t.expired = append(t.expired, gi)
	t.expiredMu.Unlock()
}

// RemoveActive removes a grant from the active grants of the player.
func (t *Tracker) RemoveActive(gi *GrantInfo) {
	t.activesMu.Lock()
	defer t.activesMu.Unlock()

	if idx := slices.Index(t.actives, gi); idx != -1 {
		t.actives = append(t.actives[:idx], t.actives[idx+1:]...)
	}
}
