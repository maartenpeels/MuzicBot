package core

type Song string
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

func (queue *SongQueue) Add(url Song) {
	queue.list = append(queue.list, url)
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

func (queue *SongQueue) Clear() {
	queue.list = make([]Song, 0)
	queue.Running = false
	queue.current = nil
}

func (queue *SongQueue) Start(sess *Session, callback func(string)) {
	queue.Running = true
	for queue.HasNext() && queue.Running {
		song := queue.Next()
		callback(string("Now playing `" + song + "`."))
		err := sess.Play(song)
		if err != nil {
			callback(string("Failed to play `" + song + "`."))
			return
		}
	}
	if !queue.Running {
		callback("Stopped playing.")
	} else {
		callback("Finished queue.")
	}
}

func (queue *SongQueue) Current() *Song {
	return queue.current
}

func (queue *SongQueue) Pause() {
	queue.Running = false
}

func NewSongQueue() *SongQueue {
	queue := new(SongQueue)
	queue.list = make([]Song, 0)
	return queue
}
