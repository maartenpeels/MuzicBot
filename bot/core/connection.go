package core

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
)

type Connection struct {
	voiceConnection *discordgo.VoiceConnection
	send            chan []int16
	lock            sync.Mutex
	sendPcm         bool
	stopRunning     bool
	playing         bool
}

func NewConnection(voiceConnection *discordgo.VoiceConnection) *Connection {
	connection := new(Connection)
	connection.voiceConnection = voiceConnection
	return connection
}
func (connection *Connection) Disconnect() {
	err := connection.voiceConnection.Disconnect()
	if err != nil {
		log.Printf("Error disconnecting from voice channel: %v", err)
		return
	}
}
