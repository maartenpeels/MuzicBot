package core

import "fmt"

type Song struct {
	Url string
	Ctx *Context
}

type SongQueue struct {
	list    []Song
	current *Song
	Running bool
}

func (queue *SongQueue) Get() []Song {
	return queue.list
}

func (queue *SongQueue) Set(list []Song) {
	queue.list = list
}

func (queue *SongQueue) Add(song Song) {
	if len(queue.Get()) > 0 {
		song.Ctx.UpdateResponse("Added `" + song.Url + "` to the queue.")
	}
	queue.list = append(queue.list, song)
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
	for queue.HasNext() && queue.Running {
		song := queue.Next()
		song.Ctx.UpdateResponse("Now playing `" + song.Url + "`.")
		err := sess.Play(song.Url)
		if err != nil {
			song.Ctx.UpdateResponse("Failed to play `" + song.Url + "`.")
			fmt.Printf("Failed to play `%s`: %v\n", song.Url, err)
			return
		}
	}

	//if !queue.Running {
	//	callback("Stopped playing.")
	//} else {
	//	callback("Finished queue.")
	//}
}

func NewSongQueue() *SongQueue {
	queue := new(SongQueue)
	queue.list = make([]Song, 0)
	return queue
}
