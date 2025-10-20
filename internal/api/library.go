package api

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

// LibraryService handles library-related API calls
type LibraryService struct {
	client *Client
}

func NewLibraryService(client *Client) *LibraryService {
	return &LibraryService{client: client}
}

// GetUserPlaylists gets the user's playlists
func (l *LibraryService) GetUserPlaylists(ctx context.Context, limit int) (*spotify.SimplePlaylistPage, error) {
	if err := l.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	playlists, err := l.client.GetSpotifyClient().CurrentUsersPlaylists(ctx, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return playlists, nil
}

// GetSavedAlbums gets the user's saved albums
func (l *LibraryService) GetSavedAlbums(ctx context.Context, limit int) (*spotify.SavedAlbumPage, error) {
	if err := l.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	albums, err := l.client.GetSpotifyClient().CurrentUsersAlbums(ctx, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return albums, nil
}

// GetSavedTracks gets the user's saved tracks
func (l *LibraryService) GetSavedTracks(ctx context.Context, limit int) (*spotify.SavedTrackPage, error) {
	if err := l.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	tracks, err := l.client.GetSpotifyClient().CurrentUsersTracks(ctx, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return tracks, nil
}

// GetSavedShows gets the user's saved shows (podcasts)
func (l *LibraryService) GetSavedShows(ctx context.Context, limit int) (*spotify.SavedShowPage, error) {
	if err := l.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	shows, err := l.client.GetSpotifyClient().CurrentUsersShows(ctx, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return shows, nil
}

// SaveTrack saves a track to the user's library
func (l *LibraryService) SaveTrack(ctx context.Context, trackID spotify.ID) error {
	if err := l.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := l.client.GetSpotifyClient().AddTracksToLibrary(ctx, trackID)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// RemoveTrack removes a track from the user's library
func (l *LibraryService) RemoveTrack(ctx context.Context, trackID spotify.ID) error {
	if err := l.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := l.client.GetSpotifyClient().RemoveTracksFromLibrary(ctx, trackID)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// SaveAlbum saves an album to the user's library
func (l *LibraryService) SaveAlbum(ctx context.Context, albumID spotify.ID) error {
	if err := l.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := l.client.GetSpotifyClient().AddAlbumsToLibrary(ctx, albumID)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// RemoveAlbum removes an album from the user's library
func (l *LibraryService) RemoveAlbum(ctx context.Context, albumID spotify.ID) error {
	if err := l.client.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	err := l.client.GetSpotifyClient().RemoveAlbumsFromLibrary(ctx, albumID)
	if err != nil {
		return HandleAPIError(err)
	}

	return nil
}

// SaveShow saves a show to the user's library
func (l *LibraryService) SaveShow(ctx context.Context, showID spotify.ID) error {
	// Shows not supported in current API
	return fmt.Errorf("shows not supported in current API")
}

// RemoveShow removes a show from the user's library
func (l *LibraryService) RemoveShow(ctx context.Context, showID spotify.ID) error {
	// Shows not supported in current API
	return fmt.Errorf("shows not supported in current API")
}
