package cmd

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"muzicBot/bot/music"
	"muzicBot/bot/queue"
)

func PlayCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Play command received")

	go music.Add(queue.Item{
		RequestedBy: i.Member.User.Username,
		Session:     s,
		Interaction: i.Interaction,
		URL:         i.ApplicationCommandData().Options[0].StringValue(),
	})

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Working...",
		},
	})

	if err != nil {
		log.Printf("Error sending interaction response: %s", err)
	}
}
