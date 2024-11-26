package music

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lrstanley/go-ytdlp"
	"io"
	"layeh.com/gopus"
	"log"
	"muzicBot/bot/queue"
	"os/exec"
)

const (
	CHANNELS   int = 2
	FRAME_RATE int = 48000
	FRAME_SIZE int = 960
	MAX_BYTES  int = (FRAME_SIZE * 2) * 2
)

func GetAudioURL(url string) (string, error) {
	ytdlp.MustInstall(context.TODO(), nil)

	ytd := ytdlp.New()
	ytd.Format("bestaudio/best")
	ytd.NoPlaylist()
	ytd.Quiet()
	ytd.GetURL()
	ytd.SkipDownload()

	run, err := ytd.Run(context.TODO(), url)
	if err != nil {
		return "", err
	}

	return run.Stdout, nil
}

func GetVoiceChannel(item *queue.Item) *discordgo.Channel {
	guild, err := item.Session.State.Guild(item.Interaction.GuildID)
	if err != nil {
		log.Printf("Error getting guild: %s", err)
		return nil
	}

	for _, state := range guild.VoiceStates {
		if state.UserID == item.Interaction.Member.User.ID {
			channel, err := item.Session.Channel(state.ChannelID)
			if err != nil {
				log.Printf("Error getting channel: %s", err)
				return nil
			}

			return channel
		}
	}

	return nil
}

func PlayAudio(item *queue.Item, url string) error {
	voiceChannel := GetVoiceChannel(item)
	if voiceChannel == nil {
		return UpdateMessage(item, "You must be in a voice channel to use this command")
	}

	voiceConnection, err := item.Session.ChannelVoiceJoin(item.Interaction.GuildID, voiceChannel.ID, false, true)
	if err != nil {
		return UpdateMessage(item, "Error joining voice channel")
	}

	ffmpeg := exec.Command("ffmpeg", "-i", url, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegOut, err := ffmpeg.StdoutPipe()
	if err != nil {
		return UpdateMessage(item, "Error creating ffmpeg stdout pipe")
	}

	buffer := bufio.NewReaderSize(ffmpegOut, 16384)
	err = ffmpeg.Start()
	if err != nil {
		return UpdateMessage(item, "Error starting ffmpeg")
	}

	err = voiceConnection.Speaking(true)
	if err != nil {
		return UpdateMessage(item, "Error speaking in voice channel")
	}

	send := make(chan []int16, CHANNELS)
	go SendPCM(voiceConnection, send)
	for {
		audioBuffer := make([]int16, FRAME_SIZE)
		err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil
		}
		if err != nil {
			return err
		}
		send <- audioBuffer
	}
}

func UpdateMessage(item *queue.Item, content string) error {
	_, err := item.Session.InteractionResponseEdit(item.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		return err
	}
	return nil
}

func SendPCM(voice *discordgo.VoiceConnection, pcm <-chan []int16) {
	encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder error,", err)
		_ = UpdateMessage(CurrentState.Current, "Error creating encoder")
		return
	}
	for {
		receive, ok := <-pcm
		if !ok {
			fmt.Println("PCM channel closed")
			_ = UpdateMessage(CurrentState.Current, "PCM channel closed")
			return
		}
		opus, err := encoder.Encode(receive, FRAME_SIZE, MAX_BYTES)
		if err != nil {
			fmt.Println("Encoding error,", err)
			_ = UpdateMessage(CurrentState.Current, "Error encoding audio")
			return
		}
		if !voice.Ready || voice.OpusSend == nil {
			fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", voice.Ready, voice.OpusSend)
			_ = UpdateMessage(CurrentState.Current, "Discordgo not ready for opus packets")
			return
		}
		voice.OpusSend <- opus
	}
}
