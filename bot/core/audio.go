package core

import (
	"errors"
	"fmt"
	"io"
	"layeh.com/gopus"
	"os/exec"
)

func (connection *Connection) Play(song Song) error {
	if connection.playing {
		return errors.New("song already playing")
	}

	ytCmd := exec.Command("./yt-dlp", "-f", "bestaudio", "-o", "-", string(song))
	ffmpegCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")

	// Pipe yt-dlp output to ffmpeg
	ytPipe, err := ytCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("yt-dlp stdout pipe error: %w", err)
	}
	ffmpegCmd.Stdin = ytPipe

	// Start yt-dlp and ffmpeg
	if err := ytCmd.Start(); err != nil {
		return fmt.Errorf("yt-dlp start error: %w", err)
	}
	ffmpegPipe, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("ffmpeg stdout pipe error: %w", err)
	}
	if err := ffmpegCmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg start error: %w", err)
	}

	// Start speaking in the voice channel
	connection.voiceConnection.Speaking(true)
	defer connection.voiceConnection.Speaking(false)

	// Opus encoder setup
	encoder, err := gopus.NewEncoder(48000, 2, gopus.Audio)
	if err != nil {
		return fmt.Errorf("opus encoder creation error: %w", err)
	}
	encoder.SetBitrate(128000) // 128 kbps
	//encoder.SetComplexity(10)  // Maximum complexity

	buffer := make([]byte, 3840) // PCM buffer size for 48kHz stereo
	for {
		if connection.stopRunning {
			connection.stopRunning = false
			break
		}
		n, err := ffmpegPipe.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading ffmpeg output: %w", err)
		}

		// Convert []byte to []int16 for Opus encoding
		pcm := make([]int16, n/2)
		for i := 0; i < len(pcm); i++ {
			pcm[i] = int16(buffer[i*2]) | int16(buffer[i*2+1])<<8
		}

		// Encode to Opus
		opusData, err := encoder.Encode(pcm, 960, 3840) // 20ms frame (960 samples at 48kHz)
		if err != nil {
			return fmt.Errorf("opus encoding error: %w", err)
		}

		connection.voiceConnection.OpusSend <- opusData
	}

	// Wait for yt-dlp and ffmpeg to finish
	ytCmd.Wait()
	ffmpegCmd.Wait()

	connection.playing = false

	return nil
}

func (connection *Connection) Stop() {
	connection.stopRunning = true
	connection.playing = false
}
