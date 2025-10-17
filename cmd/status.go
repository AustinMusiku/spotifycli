package cmd

import (
	"context"
	"fmt"

	"github.com/AustinMusiku/spotifycli/internal/api"
	"github.com/AustinMusiku/spotifycli/internal/ui"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current playback status",
	Long:  `Show the current playback status including track, progress, and playback state.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Add aliases
	statusCmd.Aliases = []string{"s", "now"}
}

func runStatus() error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	playback, playerState, err := playbackService.GetCurrentPlayback(ctx)
	if err != nil {
		return err
	}

	if playback == nil {
		ui.PrintInfo("No playback")
		return nil
	}

	status := ui.FormatPlaybackState(playerState)
	fmt.Println(status)

	return nil
}
