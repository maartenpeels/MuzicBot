package core

import (
	"fmt"
	"time"
)

type Song struct {
	Url string
}

type SongQueue struct {
	list           []Song
	current        *Song
	Running        bool
	sessionManager *SessionManager
}

func (queue *SongQueue) Get() []Song {
	return queue.list
}

func (queue *SongQueue) Set(list []Song) {
	queue.list = list
}

func (queue *SongQueue) Add(song Song) {
	queue.list = append(queue.list, song)
	fmt.Printf("Added song to queue. Queue length: %d, Running: %v\n", len(queue.list), queue.Running)
}

func (queue *SongQueue) HasNext() bool {
	return len(queue.list) > 0
}

func (queue *SongQueue) Next() Song {
	song := queue.list[0]
	queue.list = queue.list[1:]
	queue.current = &song
	return song
}

func (queue *SongQueue) Start(sess *Session) {
	queue.Running = true
	idleTimeout := 10 * time.Second
	idleTimer := time.NewTimer(idleTimeout)

	for queue.Running {
		if !queue.HasNext() {
			select {
			case <-idleTimer.C:
				sess.SendMessage("No activity for a while, leaving the channel.")
				queue.sessionManager.Leave(*sess)
				return
			default:
				time.Sleep(200 * time.Millisecond)
			}
			continue
		}

		song := queue.Next()
		sess.SendMessage("Now playing `" + song.Url + "`.")
		err := sess.Play(song.Url)
		if err != nil {
			sess.SendMessage("Failed to play `" + song.Url + "`.")
			fmt.Printf("Failed to play `%s`: %v\n", song.Url, err)
			return
		}

		idleTimer.Reset(idleTimeout)
	}

	sess.SendMessage("Stopped playing.")
}

func (queue *SongQueue) Stop() {
	queue.Running = false
}

func NewSongQueue(sm *SessionManager) *SongQueue {
	queue := new(SongQueue)
	queue.list = make([]Song, 0)
	queue.sessionManager = sm
	return queue
}
