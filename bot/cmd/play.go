package cmd

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func PlayCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Play command received")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Playing: " + i.ApplicationCommandData().Options[0].StringValue(),
		},
	})

	if err != nil {
		log.Printf("Error sending interaction response: %s", err)
	}
}
