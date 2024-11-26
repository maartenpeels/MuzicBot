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

	sess.Queue.Clear()
	sess.Stop()
	ctx.UpdateResponse("Stopped playback.")
}
