package core

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Connection struct {
	voiceConnection *discordgo.VoiceConnection
	stopRunning     bool
	shouldSkip      bool
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
