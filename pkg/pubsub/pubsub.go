package pubsub

import (
	"sync"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
)

const (
	// ReleaseSuccess is a ReleaseStatus for success release.
	ReleaseSuccess ReleaseStatus = iota

	// ReleaseFailed is a ReleaseStatus for failed release.
	ReleaseFailed
)

// ReleaseStatus is used to code release status - success or failed.
// Please use ReleaseSuccess and ReleaseFailed contants.
type ReleaseStatus int

// ReleasePubSub is a structure for simple handling pub/sub for release dependencies.
type ReleasePubSub struct {
	subs map[uniqname.UniqName][]chan ReleaseStatus
	mu   sync.RWMutex
}

// NewReleasePubSub creates new PubSub structure.
func NewReleasePubSub() *ReleasePubSub {
	return &ReleasePubSub{
		mu:   sync.RWMutex{},
		subs: make(map[uniqname.UniqName][]chan ReleaseStatus),
	}
}

// Subscribe adds new subscription for defined key and returns notification channel.
func (ps *ReleasePubSub) Subscribe(release uniqname.UniqName) <-chan ReleaseStatus {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan ReleaseStatus, 1)
	ps.subs[release] = append(ps.subs[release], ch)

	return ch
}

// publish publishes notification for all subscribers.
func (ps *ReleasePubSub) publish(release uniqname.UniqName, status ReleaseStatus) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range ps.subs[release] {
		ch <- status
		close(ch)
	}
}

// PublishSuccess publishes success notification for all subscribers.
func (ps *ReleasePubSub) PublishSuccess(release uniqname.UniqName) {
	ps.publish(release, ReleaseSuccess)
}

// PublishFailed publishes failed notification for all subscribers.
func (ps *ReleasePubSub) PublishFailed(release uniqname.UniqName) {
	ps.publish(release, ReleaseFailed)
}
