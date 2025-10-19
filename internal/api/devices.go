package api

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

// DeviceService handles device-related API calls
type DeviceService struct {
	client *Client
}

func NewDeviceService(client *Client) *DeviceService {
	return &DeviceService{client: client}
}

// GetDevices gets all available devices for the authenticated user
func (d *DeviceService) GetDevices(ctx context.Context) ([]spotify.PlayerDevice, error) {
	if err := d.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	devices, err := d.client.GetSpotifyClient().PlayerDevices(ctx)
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return devices, nil
}

// TransferPlayback transfers playback to a specific device
func (d *DeviceService) TransferPlayback(ctx context.Context, deviceID spotify.ID, play bool) error {
	if err := d.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := d.client.GetSpotifyClient().TransferPlayback(ctx, deviceID, play)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}
