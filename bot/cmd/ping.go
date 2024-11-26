package cmd

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func PingCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})

	if err != nil {
		log.Printf("Error sending interaction response: %s", err)
	}
}
