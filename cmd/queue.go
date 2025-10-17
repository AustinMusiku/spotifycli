package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/AustinMusiku/spotifycli/internal/api"
	"github.com/AustinMusiku/spotifycli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

// queueCmd represents the queue command
var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Show current queue",
	Long:  `Show the current playback queue.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQueue()
	},
}

// queueAddCmd represents the queue add command
var queueAddCmd = &cobra.Command{
	Use:   "add <URI or search query>",
	Short: "Add track to queue",
	Long:  `Add a track to the current queue. Can be a Spotify URI or a search query.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQueueAdd(args[0])
	},
}

func init() {
	rootCmd.AddCommand(queueCmd)
	queueCmd.AddCommand(queueAddCmd)

	// Add aliases
	queueCmd.Aliases = []string{"q"}
}

func runQueue() error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	queue, err := playbackService.GetQueue(ctx)
	if err != nil {
		return err
	}

	if queue == nil {
		ui.PrintInfo("No queue")
		return nil
	}

	// Display current track
	if queue.CurrentlyPlaying.Name != "" {
		fmt.Println("Currently Playing:")
		fmt.Printf("  %s\n", ui.FormatTrack(queue.CurrentlyPlaying))
	}

	// Display queue
	if len(queue.Items) > 0 {
		fmt.Println("\nQueue:")
		for i, track := range queue.Items {
			fmt.Printf("  %d. %s\n", i+1, ui.FormatTrack(track))
		}
	} else {
		ui.PrintInfo("Queue is empty")
	}

	return nil
}

func runQueueAdd(query string) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	// Get active device
	deviceID, err := getActiveDevice(ctx, client)
	if err != nil {
		return err
	}

	var uri spotify.URI

	// Check if it's a URI
	if strings.HasPrefix(query, "spotify:") {
		uri = spotify.URI(query)
	} else {
		// Search for track
		searchService := api.NewSearchService(client)
		results, err := searchService.SearchTracks(ctx, query, 1)
		if err != nil {
			return err
		}

		if len(results.Tracks.Tracks) == 0 {
			return fmt.Errorf("no tracks found for query: %s", query)
		}

		uri = results.Tracks.Tracks[0].URI
	}

	id, _, err := api.ParseURI(string(uri))
	if err != nil {
		return err
	}

	// Add to queue
	err = playbackService.AddToQueue(ctx, id, deviceID)
	if err != nil {
		return err
	}

	ui.PrintSuccess("Added to queue")
	return nil
}
