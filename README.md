# SpotifyCLI

A comprehensive command-line interface for Spotify built with Go. Control your Spotify playback, manage your library, and search content directly from the terminal.

## Features

- **Playback Control**: Play, pause, next, previous, volume, shuffle, repeat
- **Advanced Search**: Search across tracks, albums, artists, playlists, shows, and episodes
- **Library Management**: Save and manage your tracks, albums, and shows
- **Device Management**: List and switch between available devices
- **Queue Management**: View current queue and add tracks
- **Secure Authentication**: OAuth2 PKCE flow with encrypted token storage

## Installation

### Prerequisites

- Go 1.21 or later
- A Spotify account
- A Spotify app (for authentication)

### Building from Source

```bash
git clone https://github.com/AustinMusiku/spotifycli.git
cd spotifycli
go build -o spotifycli
```

### Setting up Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Add redirect URI: `http://127.0.0.1:8080/callback`
4. Copy your Client ID

## Quick Start

1. **Authenticate with Spotify**:
   ```bash
   ./spotifycli login
   ```
   Enter your Client ID when prompted, and the app will open your browser for authentication.

2. **Check your status**:
   ```bash
   ./spotifycli status
   ```

3. **Play music**:
   ```bash
   ./spotifycli play "your favorite song"
   ```

## Commands

### Authentication

- `spotifycli login` - Authenticate with Spotify
- `spotifycli logout` - Clear stored credentials

### Playback Control

- `spotifycli play [query]` - Start/resume playback or play specific content
- `spotifycli pause` - Pause playback
- `spotifycli next` - Skip to next track
- `spotifycli previous` - Go to previous track
- `spotifycli volume <0-100>` - Set volume
- `spotifycli shuffle <on|off>` - Toggle shuffle
- `spotifycli repeat <off|track|context>` - Set repeat mode

### Status & Queue

- `spotifycli status` - Show current playback status
- `spotifycli queue` - Show current queue
- `spotifycli queue add <query>` - Add track to queue

### Search

- `spotifycli search <query>` - Search all content types
- `spotifycli search <query> --type track` - Search specific content type
- `spotifycli search <query> --limit 10` - Limit results

### Library Management

- `spotifycli library playlists` - List your playlists
- `spotifycli library albums` - List saved albums
- `spotifycli library tracks` - List saved tracks
- `spotifycli library shows` - List saved shows (podcasts)
- `spotifycli library save <URI>` - Save item to library
- `spotifycli library remove <URI>` - Remove item from library

### Device Management

- `spotifycli devices` - List available devices
- `spotifycli device <name>` - Switch to specific device

## Command Aliases

For faster usage, common commands have aliases:

- `play` → `p`
- `pause` → `pa`
- `next` → `n`
- `previous` → `prev`, `b`
- `status` → `s`, `now`
- `queue` → `q`

## Examples

### Basic Playback

```bash
# Play a specific song
spotifycli play "watendawili Cham thum"

# Play from a Spotify URI
spotifycli play "spotify:track:4uLU6hMCjMI75M1A2tKUQC"

# Resume playback
spotifycli play

# Pause
spotifycli pause

# Skip to next track
spotifycli next
```

### Search and Discovery

```bash
# Search for tracks
spotifycli search "This side feat YG" --type track

# Search for albums only
spotifycli search "Treasure Self Love" --type album

# Search for podcasts
spotifycli search "The Pragmatic Engineer" --type show
```

### Library Management

```bash
# List your saved tracks
spotifycli library tracks

# Save a track to your library
spotifycli library save "spotify:track:4uLU6hMCjMI75M1A2tKUQC"

# List your playlists
spotifycli library playlists
```

### Device Control

```bash
# List available devices
spotifycli devices

# Switch to a specific device
spotifycli device "My Computer"
```

## Configuration

The app stores configuration in `~/.config/spotifycli.json`. This includes:
- Client ID (from your Spotify developer dashboard app)
- Encrypted access and refresh tokens
- Token expiration times

Notes:
- Tokens are encrypted using AES-256-GCM.
- You can provide an encryption key via the SPOTIFYCLI_KEY environment variable; otherwise a machine-derived key is used.”

## Security

- Uses OAuth2 PKCE flow for secure authentication
- Encryption uses SPOTIFYCLI_KEY if set, otherwise uses machine-derived key.
- Tokens are encrypted using AES-256-GCM
- No client secrets stored in the application
- Automatic token refresh handling (in active development)

## Troubleshooting

### Authentication Issues

If you encounter authentication problems:

1. Make sure your Spotify app has the correct redirect URI: `http://127.0.0.1:8080/callback`. Spotify does not accept `localhost` unless explicitly bound to either the IPv4 or IPv6 loopback address.
2. Check that your Client ID is correct
3. Try logging out and logging back in: `spotifycli logout && spotifycli login`

### No Devices Found

If you see "No devices found":

1. Make sure Spotify is running on at least one device
2. Try refreshing your devices: `spotifycli devices`

### Rate Limiting

If you hit rate limits:

1. Wait a few minutes before trying again
2. The app will automatically handle rate limiting

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
- Uses [zmb3/spotify](https://github.com/zmb3/spotify) for Spotify Web API integration
- Inspired by [brianstrauch/spotify-cli](https://github.com/brianstrauch/spotify-cli)

