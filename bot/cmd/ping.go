package cmd

import (
	"muzicBot/bot/core"
)

func PingCommandHandler(ctx *core.Context) {
	ctx.Respond("Pong!")
}
