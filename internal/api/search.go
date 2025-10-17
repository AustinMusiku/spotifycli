package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/zmb3/spotify/v2"
)

// SearchService handles search-related API calls
type SearchService struct {
	client *Client
}

// NewSearchService creates a new search service
func NewSearchService(client *Client) *SearchService {
	return &SearchService{client: client}
}

// SearchResult represents a search result
type SearchResult struct {
	Tracks    *spotify.FullTrackPage
	Albums    *spotify.SimpleAlbumPage
	Artists   *spotify.FullArtistPage
	Playlists *spotify.SimplePlaylistPage
	Shows     *spotify.SimpleShowPage
	Episodes  *spotify.SimpleEpisodePage
}

// Search performs a search across all content types
func (s *SearchService) Search(ctx context.Context, query string, limit int) (*SearchResult, error) {
	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	// Search all content types
	results, err := s.client.GetSpotifyClient().Search(ctx, query, spotify.SearchTypeTrack|spotify.SearchTypeAlbum|spotify.SearchTypeArtist|spotify.SearchTypePlaylist|spotify.SearchTypeShow|spotify.SearchTypeEpisode, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return &SearchResult{
		Tracks:    results.Tracks,
		Albums:    results.Albums,
		Artists:   results.Artists,
		Playlists: results.Playlists,
		Shows:     results.Shows,
		Episodes:  results.Episodes,
	}, nil
}

// SearchTracks searches for tracks only
func (s *SearchService) SearchTracks(ctx context.Context, query string, limit int) (*spotify.SearchResult, error) {
	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	results, err := s.client.GetSpotifyClient().Search(ctx, query, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return results, nil
}

// SearchAlbums searches for albums only
func (s *SearchService) SearchAlbums(ctx context.Context, query string, limit int) (*spotify.SearchResult, error) {
	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	results, err := s.client.GetSpotifyClient().Search(ctx, query, spotify.SearchTypeAlbum, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return results, nil
}

// SearchArtists searches for artists only
func (s *SearchService) SearchArtists(ctx context.Context, query string, limit int) (*spotify.SearchResult, error) {
	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	results, err := s.client.GetSpotifyClient().Search(ctx, query, spotify.SearchTypeArtist, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return results, nil
}

// SearchPlaylists searches for playlists only
func (s *SearchService) SearchPlaylists(ctx context.Context, query string, limit int) (*spotify.SearchResult, error) {
	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	results, err := s.client.GetSpotifyClient().Search(ctx, query, spotify.SearchTypePlaylist, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return results, nil
}

// SearchShows searches for shows (podcasts) only
func (s *SearchService) SearchShows(ctx context.Context, query string, limit int) (*spotify.SearchResult, error) {
	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	results, err := s.client.GetSpotifyClient().Search(ctx, query, spotify.SearchTypeShow, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return results, nil
}

// SearchEpisodes searches for episodes only
func (s *SearchService) SearchEpisodes(ctx context.Context, query string, limit int) (*spotify.SearchResult, error) {
	if err := s.client.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	results, err := s.client.GetSpotifyClient().Search(ctx, query, spotify.SearchTypeEpisode, spotify.Limit(limit))
	if err != nil {
		return nil, HandleAPIError(err)
	}

	return results, nil
}

// ParseURI parses a Spotify URI and returns the type and ID
func ParseURI(uri string) (spotify.ID, spotify.URI, error) {
	// Remove spotify: prefix if present
	uri = strings.TrimPrefix(uri, "spotify:")

	// Split by colon
	parts := strings.Split(uri, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid Spotify URI format: %s", uri)
	}

	contentType := parts[0]
	id := parts[1]

	// Validate content type
	validTypes := []string{"track", "album", "artist", "playlist", "show", "episode"}
	for _, validType := range validTypes {
		if contentType == validType {
			return spotify.ID(id), spotify.URI("spotify:" + uri), nil
		}
	}

	return "", "", fmt.Errorf("unsupported content type: %s", contentType)
}
