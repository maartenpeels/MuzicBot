package cmd

import (
	"github.com/bwmarrin/discordgo"
	"muzicBot/bot/core"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Ping the bot",
	},
	{
		Name:        "play",
		Description: "Play music",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "The url to stream music from",
				Required:    true,
			},
		},
	},
	{
		Name:		 "skip",
		Description: "Skip the current song"
	}
	{
		Name:        "stop",
		Description: "Stop the music",
	},
}

var CommandHandlers = map[string]func(ctx *core.Context){
	"ping": PingCommandHandler,
	"play": PlayCommandHandler,
	"skip": SkipCommandHandler,
	"stop": StopCommandHandler,
}
