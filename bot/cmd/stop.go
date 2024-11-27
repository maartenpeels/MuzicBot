package cmd

import (
	"muzicBot/bot/core"
)

func StopCommandHandler(ctx *core.Context) {
	ctx.Respond("Stopping playback...")

	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil {
		ctx.UpdateResponse("I'm not playing anything right now!")
		return
	}

	ctx.Sessions.Leave(*sess)
	ctx.UpdateResponse("Stopped playback.")
}
