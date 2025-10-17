package cmd

import (
	"context"
	"fmt"
	"strconv"
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

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause playback",
	Long:  `Pause the current playback.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPause()
	},
}

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Skip to next track",
	Long:  `Skip to the next track in the current queue.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNext()
	},
}

// previousCmd represents the previous command
var previousCmd = &cobra.Command{
	Use:   "previous",
	Short: "Go to previous track",
	Long:  `Go to the previous track in the current queue.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPrevious()
	},
}

// volumeCmd represents the volume command
var volumeCmd = &cobra.Command{
	Use:   "volume <0-100>",
	Short: "Set volume",
	Long:  `Set the volume for the active device (0-100).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVolume(args[0])
	},
}

// shuffleCmd represents the shuffle command
var shuffleCmd = &cobra.Command{
	Use:   "shuffle <on|off>",
	Short: "Toggle shuffle",
	Long:  `Toggle shuffle mode on or off.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runShuffle(args[0])
	},
}

// repeatCmd represents the repeat command
var repeatCmd = &cobra.Command{
	Use:   "repeat <off|track|context>",
	Short: "Set repeat mode",
	Long:  `Set repeat mode: off (no repeat), track (repeat current track), or context (repeat current playlist/album).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRepeat(args[0])
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(previousCmd)
	rootCmd.AddCommand(volumeCmd)
	rootCmd.AddCommand(shuffleCmd)
	rootCmd.AddCommand(repeatCmd)

	// Add aliases
	playCmd.Aliases = []string{"p"}
	pauseCmd.Aliases = []string{"p"}
	nextCmd.Aliases = []string{"n"}
	previousCmd.Aliases = []string{"prev", "b"}
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

func runPause() error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	deviceID, err := getActiveDevice(ctx, client)
	if err != nil {
		return err
	}

	err = playbackService.Pause(ctx, deviceID)
	if err != nil {
		return err
	}

	ui.PrintSuccess("Paused playback")
	return nil
}

func runNext() error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	deviceID, err := getActiveDevice(ctx, client)
	if err != nil {
		return err
	}

	err = playbackService.Next(ctx, deviceID)
	if err != nil {
		return err
	}

	ui.PrintSuccess("Skipped to next track")
	return nil
}

func runPrevious() error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	deviceID, err := getActiveDevice(ctx, client)
	if err != nil {
		return err
	}

	err = playbackService.Previous(ctx, deviceID)
	if err != nil {
		return err
	}

	ui.PrintSuccess("Went to previous track")
	return nil
}

func runVolume(volumeStr string) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	volume, err := strconv.Atoi(volumeStr)
	if err != nil {
		return fmt.Errorf("invalid volume: %s (must be a number)", volumeStr)
	}

	if volume < 0 || volume > 100 {
		return fmt.Errorf("volume must be between 0 and 100")
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	deviceID, err := getActiveDevice(ctx, client)
	if err != nil {
		return err
	}

	err = playbackService.SetVolume(ctx, deviceID, volume)
	if err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Set volume to %d%%", volume))
	return nil
}

func runShuffle(state string) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	shuffle := state == "on"
	if state != "on" && state != "off" {
		return fmt.Errorf("invalid shuffle state: %s (must be 'on' or 'off')", state)
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	deviceID, err := getActiveDevice(ctx, client)
	if err != nil {
		return err
	}

	err = playbackService.SetShuffle(ctx, deviceID, shuffle)
	if err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Shuffle %s", state))
	return nil
}

func runRepeat(state string) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	if state != "off" && state != "track" && state != "context" {
		return fmt.Errorf("invalid repeat state: %s (must be 'off', 'track', or 'context')", state)
	}

	playbackService := api.NewPlaybackService(client)
	ctx := context.Background()

	deviceID, err := getActiveDevice(ctx, client)
	if err != nil {
		return err
	}

	err = playbackService.SetRepeat(ctx, deviceID, state)
	if err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Repeat set to %s", state))
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
