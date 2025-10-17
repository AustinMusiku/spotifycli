package cmd

import (
	"context"
	"fmt"

	"github.com/AustinMusiku/spotifycli/internal/api"
	"github.com/AustinMusiku/spotifycli/internal/ui"
	"github.com/spf13/cobra"
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

func init() {
	rootCmd.AddCommand(queueCmd)

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
