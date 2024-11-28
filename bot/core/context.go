package core

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Context struct {
	Discord     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Guild       *discordgo.Guild
	Sessions    *SessionManager
}

func NewContext(discord *discordgo.Session, integration *discordgo.InteractionCreate, guild *discordgo.Guild, sessions *SessionManager) *Context {
	return &Context{
		Discord:     discord,
		Interaction: integration,
		Guild:       guild,
		Sessions:    sessions,
	}
}

// Respond sends a response to the interaction (use only once in the handler)
func (ctx *Context) Respond(response string) {
	err := ctx.Discord.InteractionRespond(ctx.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}
}

// UpdateResponse updates the response to the interaction (use multiple times after Respond)
func (ctx *Context) UpdateResponse(response string) {
	_, err := ctx.Discord.InteractionResponseEdit(ctx.Interaction.Interaction, &discordgo.WebhookEdit{
		Content: &response,
	})
	if err != nil {
		log.Printf("Error updating interaction response: %v", err)
		return
	}
}

func (ctx *Context) GetVoiceChannel() *discordgo.Channel {
	for _, state := range ctx.Guild.VoiceStates {
		if state.UserID == ctx.Interaction.Member.User.ID {
			channel, _ := ctx.Discord.State.Channel(state.ChannelID)
			return channel
		}
	}

	return nil
}
