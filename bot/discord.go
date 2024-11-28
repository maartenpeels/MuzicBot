package bot

import (
	"fmt"
	"log"
	"muzicBot/bot/cmd"
	"muzicBot/bot/core"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// TODO: handle disconnect / kick from voice channel

var Sessions *core.SessionManager

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

	Sessions = core.NewSessionManager()

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

	permissions := discordgo.PermissionVoiceUseVAD | discordgo.PermissionVoiceConnect | discordgo.PermissionViewChannel | discordgo.PermissionVoiceSpeak | discordgo.PermissionSendMessages | discordgo.PermissionManageMessages
	invite := fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=%d", discord.State.User.ID, permissions)
	log.Printf("Invite URL: %s", invite)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Printf("Bot is now running. Press CTRL+C to exit.")
	<-stop

	log.Printf("Shutting down")
	Sessions.LeaveAll()
}

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handler, ok := cmd.CommandHandlers[i.ApplicationCommandData().Name]
	if !ok {
		log.Printf("No handler for command %s", i.ApplicationCommandData().Name)
		return
	}

	channel, err := s.State.Channel(i.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel,", err)
		return
	}

	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Error getting guild,", err)
		return
	}

	ctx := core.NewContext(s, i, guild, Sessions)
	handler(ctx)
}
