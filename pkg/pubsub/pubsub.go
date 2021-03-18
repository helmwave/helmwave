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

func NewReleasePubSub() *ReleasePubSub {
	return &ReleasePubSub{
		subs: make(map[string][]chan ReleaseStatus),
	}
}

func (ps *ReleasePubSub) Subscribe(release string) <-chan ReleaseStatus {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan ReleaseStatus, 1)
	ps.subs[release] = append(ps.subs[release], ch)
	return ch
}

func (ps *ReleasePubSub) publish(release string, status ReleaseStatus) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range ps.subs[release] {
		ch <- status
		close(ch)
	}
}

func (ps *ReleasePubSub) PublishSuccess(release string) {
	ps.publish(release, ReleaseSuccess)
}

func (ps *ReleasePubSub) PublishFailed(release string) {
	ps.publish(release, ReleaseFailed)
}
