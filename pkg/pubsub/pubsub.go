package pubsub

import "sync"

const (
	ReleaseSuccess ReleaseStatus = iota
	ReleaseFailed
)

type ReleaseStatus int

type ReleasePubSub struct {
	mu   sync.RWMutex
	subs map[string][]chan ReleaseStatus
}

// NewReleasePubSub creates new PubSub structure.
func NewReleasePubSub() *ReleasePubSub {
	return &ReleasePubSub{
		mu:   sync.RWMutex{},
		subs: make(map[string][]chan ReleaseStatus),
	}
}

// Subscribe adds new subscription for defined key and returns notification channel.
func (ps *ReleasePubSub) Subscribe(release string) <-chan ReleaseStatus {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan ReleaseStatus, 1)
	ps.subs[release] = append(ps.subs[release], ch)

	return ch
}

// publish publishes notification for all subscribers.
func (ps *ReleasePubSub) publish(release string, status ReleaseStatus) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range ps.subs[release] {
		ch <- status
		close(ch)
	}
}

// PublishSuccess publishes success notification for all subscribers.
func (ps *ReleasePubSub) PublishSuccess(release string) {
	ps.publish(release, ReleaseSuccess)
}

// PublishFailed publishes failed notification for all subscribers.
func (ps *ReleasePubSub) PublishFailed(release string) {
	ps.publish(release, ReleaseFailed)
}
