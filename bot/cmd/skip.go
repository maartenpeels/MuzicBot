package cmd

import (
	"muzicBot/bot/core"
)

func SkipCommandHandler(ctx *core.Context) {
	ctx.Respond("Skipping...")

	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)

	if sess == nil {
		ctx.UpdateResponse("I'm not playing anything right now!")
		return
	}

	sess.Skip()
}
