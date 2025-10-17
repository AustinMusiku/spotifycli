package api

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

// PlaybackService handles playback-related API calls
type PlaybackService struct {
	client *Client
}

// NewPlaybackService creates a new playback service
func NewPlaybackService(client *Client) *PlaybackService {
	return &PlaybackService{client: client}
}

// GetCurrentPlayback gets the current playback state
func (p *PlaybackService) GetCurrentPlayback(ctx context.Context) (*spotify.CurrentlyPlaying, *spotify.PlayerState, error) {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return nil, nil, err
	}

	playback, err := p.client.GetSpotifyClient().PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return nil, nil, HandleAPIError(err)
	}

	playerState, error := p.client.GetSpotifyClient().PlayerState(ctx)
	if error != nil {
		return nil, nil, HandleAPIError(error)
	}

	return playback, playerState, nil
}

// Play starts or resumes playback
func (p *PlaybackService) Play(ctx context.Context, deviceID spotify.ID, uri spotify.URI) error {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	opts := &spotify.PlayOptions{
		DeviceID: &deviceID,
	}

	if uri != "" {
		opts.URIs = []spotify.URI{uri}
	}

	err := p.client.GetSpotifyClient().PlayOpt(ctx, opts)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// Pause pauses playback
func (p *PlaybackService) Pause(ctx context.Context, deviceID spotify.ID) error {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := p.client.GetSpotifyClient().PauseOpt(ctx, &spotify.PlayOptions{
		DeviceID: &deviceID,
	})
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// Next skips to the next track
func (p *PlaybackService) Next(ctx context.Context, deviceID spotify.ID) error {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := p.client.GetSpotifyClient().NextOpt(ctx, &spotify.PlayOptions{
		DeviceID: &deviceID,
	})
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// Previous goes to the previous track
func (p *PlaybackService) Previous(ctx context.Context, deviceID spotify.ID) error {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := p.client.GetSpotifyClient().PreviousOpt(ctx, &spotify.PlayOptions{
		DeviceID: &deviceID,
	})
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// SetVolume sets the volume for a device
func (p *PlaybackService) SetVolume(ctx context.Context, deviceID spotify.ID, volume int) error {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if volume < 0 || volume > 100 {
		return fmt.Errorf("volume must be between 0 and 100")
	}

	err := p.client.GetSpotifyClient().Volume(ctx, volume)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// SetShuffle sets the shuffle state
func (p *PlaybackService) SetShuffle(ctx context.Context, deviceID spotify.ID, shuffle bool) error {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := p.client.GetSpotifyClient().Shuffle(ctx, shuffle)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// SetRepeat sets the repeat state
func (p *PlaybackService) SetRepeat(ctx context.Context, deviceID spotify.ID, state string) error {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if state != "off" && state != "track" && state != "context" {
		return fmt.Errorf("invalid repeat state: %s (must be 'off', 'track', or 'context')", state)
	}

	err := p.client.GetSpotifyClient().Repeat(ctx, state)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// GetQueue gets the current queue
func (p *PlaybackService) GetQueue(ctx context.Context) (*spotify.Queue, error) {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	queue, err := p.client.GetSpotifyClient().GetQueue(ctx)
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return queue, nil
}

// AddToQueue adds a track to the queue
func (p *PlaybackService) AddToQueue(ctx context.Context, id spotify.ID, deviceID spotify.ID) error {
	if err := p.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := p.client.GetSpotifyClient().QueueSong(ctx, id)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}
