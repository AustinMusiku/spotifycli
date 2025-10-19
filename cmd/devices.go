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

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "List available devices",
	Long:  `List all available Spotify devices.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDevices()
	},
}

// deviceCmd represents the device command
var deviceCmd = &cobra.Command{
	Use:   "device <name>",
	Short: "Switch active device",
	Long:  `Switch playback to a specific device by name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDevice(args[0])
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	rootCmd.AddCommand(deviceCmd)
}

func runDevices() error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	deviceService := api.NewDeviceService(client)
	ctx := context.Background()

	devices, err := deviceService.GetDevices(ctx)
	if err != nil {
		return err
	}

	if len(devices) == 0 {
		ui.PrintInfo("No devices found. Please start Spotify on a device.")
		return nil
	}

	fmt.Println("🎵 Available Devices:")
	for i, device := range devices {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatDevice(device))
	}

	return nil
}

func runDevice(deviceName string) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	deviceService := api.NewDeviceService(client)
	ctx := context.Background()

	devices, err := deviceService.GetDevices(ctx)
	if err != nil {
		return err
	}

	// Find device by name (case-insensitive)
	var targetDevice spotify.PlayerDevice
	found := false

	for _, device := range devices {
		if strings.EqualFold(device.Name, deviceName) {
			targetDevice = device
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("device not found: %s", deviceName)
	}

	// Transfer playback to the device
	err = deviceService.TransferPlayback(ctx, targetDevice.ID, false)
	if err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Switched to device: %s", targetDevice.Name))
	return nil
}
