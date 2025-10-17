package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/zmb3/spotify/v2"
)

// Colors for different output types
var (
	SuccessColor = color.New(color.FgGreen)
	ErrorColor   = color.New(color.FgRed)
	InfoColor    = color.New(color.FgCyan)
	WarningColor = color.New(color.FgYellow)
	BoldColor    = color.New(color.Bold)
	DimColor     = color.New(color.Faint)
)

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	SuccessColor.Println("✓", message)
}

// PrintError prints an error message
func PrintError(message string) {
	ErrorColor.Println("✗", message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	InfoColor.Println("ℹ", message)
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	WarningColor.Println("⚠", message)
}

// FormatDuration formats a duration in milliseconds to MM:SS
func FormatDuration(ms int) string {
	duration := time.Duration(ms) * time.Millisecond
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// FormatTrack formats a track for display
func FormatTrack(track spotify.FullTrack) string {
	artistNames := make([]string, len(track.Artists))
	for i, artist := range track.Artists {
		artistNames[i] = artist.Name
	}

	artists := strings.Join(artistNames, ", ")
	duration := FormatDuration(int(track.Duration))

	return fmt.Sprintf("%s - %s (%s) [%s]",
		BoldColor.Sprint(artists),
		track.Name,
		DimColor.Sprint(track.Album.Name),
		duration)
}

// FormatSimpleTrack formats a simple track for display
func FormatSimpleTrack(track spotify.SimpleTrack) string {
	artistNames := make([]string, len(track.Artists))
	for i, artist := range track.Artists {
		artistNames[i] = artist.Name
	}

	artists := strings.Join(artistNames, ", ")
	duration := FormatDuration(int(track.Duration))

	return fmt.Sprintf("%s - %s [%s]",
		BoldColor.Sprint(artists),
		track.Name,
		duration)
}

// FormatAlbum formats an album for display
func FormatAlbum(album spotify.SimpleAlbum) string {
	artistNames := make([]string, len(album.Artists))
	for i, artist := range album.Artists {
		artistNames[i] = artist.Name
	}

	artists := strings.Join(artistNames, ", ")
	year := album.ReleaseDate
	if len(year) > 4 {
		year = year[:4]
	}

	return fmt.Sprintf("%s - %s (%s)",
		BoldColor.Sprint(artists),
		album.Name,
		DimColor.Sprint(year))
}

// FormatArtist formats an artist for display
func FormatArtist(artist spotify.FullArtist) string {
	genres := strings.Join(artist.Genres, ", ")
	if genres != "" {
		return fmt.Sprintf("%s (%s)",
			BoldColor.Sprint(artist.Name),
			DimColor.Sprint(genres))
	}
	return BoldColor.Sprint(artist.Name)
}

// FormatPlaylist formats a playlist for display
func FormatPlaylist(playlist spotify.SimplePlaylist) string {
	trackCount := playlist.Tracks.Total
	return fmt.Sprintf("%s (%d tracks)",
		BoldColor.Sprint(playlist.Name),
		trackCount)
}

// FormatShow formats a show for display
func FormatShow(show spotify.SimpleShow) string {
	return fmt.Sprintf("%s - %s",
		BoldColor.Sprint(show.Name),
		show.Description)
}

// FormatEpisode formats an episode for display
func FormatEpisode(episode spotify.EpisodePage) string {
	duration := FormatDuration(int(episode.Duration_ms))
	return fmt.Sprintf("%s - %s [%s]",
		BoldColor.Sprint(episode.Name),
		episode.Show.Name,
		duration)
}

// FormatDevice formats a device for display
func FormatDevice(device spotify.PlayerDevice) string {
	status := "inactive"
	if device.Active {
		status = "active"
	}

	volume := ""
	if device.Volume > 0 {
		volume = fmt.Sprintf(" (%d%%)", device.Volume)
	}

	return fmt.Sprintf("%s - %s (%s)%s",
		BoldColor.Sprint(device.Name),
		device.Type,
		status,
		volume)
}

// FormatPlaybackState formats the current playback state
func FormatPlaybackState(playback *spotify.PlayerState) string {
	if playback == nil || playback.Item == nil {
		return "No playback"
	}

	track := playback.Item
	artistNames := make([]string, len(track.Artists))
	for i, artist := range track.Artists {
		artistNames[i] = artist.Name
	}

	artists := strings.Join(artistNames, ", ")

	// Progress bar
	progress := float64(playback.Progress) / float64(track.Duration)
	progressBar := createProgressBar(progress, 20)

	// Shuffle and repeat indicators
	shuffle := ""
	if playback.ShuffleState {
		shuffle = " 🔀"
	}

	repeat := ""
	switch playback.RepeatState {
	case "track":
		repeat = " 🔁"
	case "context":
		repeat = " 🔂"
	}

	return fmt.Sprintf("♪ %s - %s\n   %s %s%s%s",
		BoldColor.Sprint(artists),
		track.Name,
		progressBar,
		FormatDuration(int(playback.Progress)),
		DimColor.Sprint("/"+FormatDuration(int(track.Duration))),
		shuffle+repeat)
}

// createProgressBar creates a visual progress bar
func createProgressBar(progress float64, width int) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filled := int(progress * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s]", bar)
}

// PrintTable prints a formatted table
func PrintTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		InfoColor.Println("No results found")
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	headerRow := ""
	for i, header := range headers {
		headerRow += fmt.Sprintf("%-*s", widths[i]+2, header)
	}
	BoldColor.Println(headerRow)

	// Print separator
	separator := ""
	for _, width := range widths {
		separator += strings.Repeat("-", width+2)
	}
	DimColor.Println(separator)

	// Print rows
	for _, row := range rows {
		rowStr := ""
		for i, cell := range row {
			rowStr += fmt.Sprintf("%-*s", widths[i]+2, cell)
		}
		fmt.Println(rowStr)
	}
}
