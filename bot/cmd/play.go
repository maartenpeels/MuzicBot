package cmd

import (
	"muzicBot/bot/core"
)

func PlayCommandHandler(ctx *core.Context) {
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)

	if sess == nil {
		vc := ctx.GetVoiceChannel()
		if vc == nil {
			ctx.Respond("You must be in a voice channel to use the bot!")
			return
		}
		newSess, err := ctx.Sessions.Join(ctx.Discord, ctx.Guild.ID, vc.ID, core.JoinProperties{
			Muted:    false,
			Deafened: false,
		})
		sess = newSess
		if err != nil {
			ctx.Respond("Failed to join voice channel")
			return
		}
	}

	url := ctx.Interaction.ApplicationCommandData().Options[0].StringValue()
	sess.Queue.Add(core.Song{Url: url})
	ctx.Respond("Added `" + url + "` to the queue.")
}
