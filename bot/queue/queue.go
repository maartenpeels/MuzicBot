package queue

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type Item struct {
	RequestedBy string
	Session     *discordgo.Session
	Interaction *discordgo.Interaction
	URL         string
}

type MusicQueue struct {
	items []Item
	mu    sync.Mutex
}

func (q *MusicQueue) Add(item Item) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *MusicQueue) Skip() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) > 0 {
		q.items = q.items[1:]
	}
}

func (q *MusicQueue) Pop() *Item {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return &item
}

func (q *MusicQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}
