package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"muzicBot/bot/cmd"
	"os"
	"os/signal"
)

func Init() {
	discord, err := discordgo.New("Bot " + Env.Token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %s", err)
	}

	err = discord.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %s", err)
	}
	defer discord.Close()

	discord.AddHandler(onLoggedIn)
	discord.AddHandler(InteractionCreate)

	log.Printf("Adding commands")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(cmd.Commands))
	for i, v := range cmd.Commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Printf("Error creating command: %s", err)
			continue
		}

		registeredCommands[i] = cmd
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Printf("Bot is now running. Press CTRL+C to exit.")
	<-stop

	log.Printf("Shutting down")

	// TODO: Make removing commands a setting
	log.Printf("Removing commands")
	for _, v := range registeredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, "", v.ID)
		if err != nil {
			log.Printf("Error deleting command: %s", err)
		}
	}
}

func onLoggedIn(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s", event.User.String())
}

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handler, ok := cmd.CommandHandlers[i.ApplicationCommandData().Name]
	if !ok {
		log.Printf("No handler for command %s", i.ApplicationCommandData().Name)
		return
	}

	handler(s, i)
}
