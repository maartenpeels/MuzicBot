# MuzicBot

MuzicBot is a Discord bot that streams music from various sources directly into your voice channel. It uses `yt-dlp` to download audio and `ffmpeg` to process it before sending it to Discord.

## Features

- Stream music from URLs
- Basic commands: play, stop, ping
- Handles voice channel connections

## Requirements

- Go 1.18+
- `yt-dlp`
- `ffmpeg`
- A Discord bot token

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/maartenpeels/muzicBot.git
   cd muzicBot
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Ensure `yt-dlp` and `ffmpeg` are installed and available in your PATH.

## Configuration

Create a `.env` file in the root directory with your Discord bot token:
```
TOKEN=your_discord_bot_token
```

## Usage

### Docker

1. Run the bot:
   ```sh
   docker run -e TOKEN=your_discord_bot_token maarten1012/muzicbot
   ```

### Development

1. Build the bot:
   ```sh
   go build -o muzicBot
   ```

2. Run the bot:
   ```sh
   ./muzicBot
   ```

3. Invite the bot to your server using the invite URL printed in the console.

## Commands

- `/ping`: Check if the bot is responsive.
- `/play <url>`: Play music from the provided URL.
- `/skip`: Skip the currently playing song.
- `/stop`: Stop playing music.

## Development

To add new commands or modify existing ones, edit the files in the `cmd` and `core` directories. The main entry point for the bot is in `bot/discord.go`.

## Contributing

Feel free to submit issues or pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.