package core

import "fmt"

type SongQueue struct {
	list    []string
	current *string
	Running bool
}

func (queue *SongQueue) Get() []string {
	return queue.list
}

func (queue *SongQueue) Set(list []string) {
	queue.list = list
}

func (queue *SongQueue) Add(url string) {
	queue.list = append(queue.list, url)
}

func (queue *SongQueue) HasNext() bool {
	return len(queue.list) > 0
}

func (queue *SongQueue) Next() string {
	url := queue.list[0]
	queue.list = queue.list[1:]
	queue.current = &url
	return url
}

func (queue *SongQueue) Start(sess *Session, callback func(string)) {
	queue.Running = true
	for queue.HasNext() && queue.Running {
		url := queue.Next()
		callback("Now playing `" + url + "`.")
		err := sess.Play(url)
		if err != nil {
			callback("Failed to play `" + url + "`.")
			fmt.Printf("Failed to play `%s`: %v\n", url, err)
			return
		}
	}
	if !queue.Running {
		callback("Stopped playing.")
	} else {
		callback("Finished queue.")
	}
}

func (queue *SongQueue) Pause() {
	queue.Running = false
}

func NewSongQueue() *SongQueue {
	queue := new(SongQueue)
	queue.list = make([]string, 0)
	return queue
}
