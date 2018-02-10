package topic

import "fmt"

type Topic interface {
	AddListener() (list <-chan Post, canc func())
	GetPipe() chan<- Post
	Run()
	Close()
	String() string
}

type topic struct {
	kind, name string
	b          chan Post
	listeners  map[chan Post]bool // Race condition.
}

// AddListener creates a channel where all new messages from users will be
// pushed. It returns a listening cancel function.
func (t *topic) AddListener() (list <-chan Post, canc func()) {
	p := make(chan Post, 10)
	t.listeners[p] = true
	return p, func() { delete(t.listeners, p) }
}

// GetPipe returns a channel where new messages from clients can be pushed.
func (t *topic) GetPipe() chan<- Post {
	return t.b
}

func (t *topic) Close() {
	close(t.b)
}

func (t *topic) String() string {
	return fmt.Sprintf("%s / %s", t.kind, t.name)
}

// Run listen to a incoming messages and broadcasts them to all listeners.
func (t *topic) Run() {
	for m := range t.b {
		for l := range t.listeners {
			l <- m
		}
	}
}

func NewTopic(kind, name string) Topic {
	return &topic{
		kind,
		name,
		make(chan Post),
		make(map[chan Post]bool)}
}
