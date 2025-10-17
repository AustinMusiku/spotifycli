package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/AustinMusiku/spotifycli/internal/api"
	"github.com/AustinMusiku/spotifycli/internal/config"
	"github.com/AustinMusiku/spotifycli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play [track/album/playlist URI or search query]",
	Short: "Start or resume playback",
	Long:  `Start or resume playback. If no argument is provided, resumes current playback. If a URI or search query is provided, plays that content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPlay(args)
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	// Add aliases
	playCmd.Aliases = []string{"p"}
}

func runPlay(args []string) error {
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

	if len(args) == 0 {
		// Resume playback
		err = playbackService.Play(ctx, deviceID, "")
		if err != nil {
			return err
		}
		ui.PrintSuccess("Resumed playback")
	} else {
		// Play specific content
		query := strings.Join(args, " ")

		// Check if it's a URI
		if strings.HasPrefix(query, "spotify:") {
			// Direct URI
			err = playbackService.Play(ctx, deviceID, spotify.URI(query))
			if err != nil {
				return err
			}
			ui.PrintSuccess(fmt.Sprintf("Playing %s", query))
		} else {
			// Search and play first result
			searchService := api.NewSearchService(client)
			results, err := searchService.SearchTracks(ctx, query, 1)
			if err != nil {
				return err
			}

			if len(results.Tracks.Tracks) == 0 {
				return fmt.Errorf("no tracks found for query: %s", query)
			}

			track := results.Tracks.Tracks[0]
			err = playbackService.Play(ctx, deviceID, track.URI)
			if err != nil {
				return err
			}

			ui.PrintSuccess(fmt.Sprintf("Playing: %s", ui.FormatTrack(track)))
		}
	}

	return nil
}

// getAuthenticatedClient returns an authenticated API client
func getAuthenticatedClient() (*config.Config, *api.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		return nil, nil, fmt.Errorf("not authenticated, please run 'spotify login'")
	}

	client := api.NewClient(cfg)
	if err := client.Authenticate(cfg.GetAccessToken()); err != nil {
		return nil, nil, fmt.Errorf("authentication failed: %w", err)
	}

	return cfg, client, nil
}

// getActiveDevice gets the active device ID
func getActiveDevice(ctx context.Context, client *api.Client) (spotify.ID, error) {
	deviceService := api.NewDeviceService(client)
	devices, err := deviceService.GetDevices(ctx)
	if err != nil {
		return "", err
	}

	if len(devices) == 0 {
		return "", fmt.Errorf("no devices found, please start Spotify on a device")
	}

	// Find active device
	for _, device := range devices {
		if device.Active {
			return device.ID, nil
		}
	}

	// If no active device, use the first one
	return devices[0].ID, nil
}
