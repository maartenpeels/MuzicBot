package music

import (
	"log"
	"muzicBot/bot/queue"
)

type State struct {
	Playing bool
	Current *queue.Item
}

var Queue = &queue.MusicQueue{}
var CurrentState = &State{}

func Add(item queue.Item) {
	if CurrentState.Playing {
		Queue.Add(item)
		log.Printf("Added %s to queue", item.URL)
		_ = UpdateMessage(&item, "Added to queue: "+item.URL)
		return
	}

	play(&item)
}

func next() {
	item := Queue.Pop()
	if item == nil {
		CurrentState.Playing = false
		CurrentState.Current = nil
		log.Printf("Queue is empty")
		return
	}

	play(item)
}

func play(item *queue.Item) {
	CurrentState.Playing = true
	CurrentState.Current = item

	log.Printf("Playing %s", item.URL)
	err := UpdateMessage(item, "Playing "+item.URL)
	if err != nil {
		next()
		return
	}

	audioUrl, err := GetAudioURL(item.URL)
	if err != nil {
		log.Printf("Error getting audio URL: %s", err)
		_ = UpdateMessage(item, "Error getting audio URL")
		next()
		return
	}

	err = PlayAudio(item, audioUrl)
	if err != nil {
		log.Printf("Error playing audio: %s", err)
		_ = UpdateMessage(item, "Error playing audio")
		next()
		return
	}

	next()
}
