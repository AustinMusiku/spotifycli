package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

// Client wraps the Spotify API client with authentication handling
type Client struct {
	spotifyClient *spotify.Client
	config        ConfigProvider
	httpClient    *http.Client
}

// ConfigProvider interface for accessing configuration
type ConfigProvider interface {
	GetAccessToken() string
	GetRefreshToken() string
	GetTokenType() string
	IsTokenExpired() bool
	SetTokens(accessToken, refreshToken, tokenType string, expiresIn int64)
	Save() error
}

func NewClient(config ConfigProvider) *Client {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Client{
		config:     config,
		httpClient: httpClient,
	}
}

func (c *Client) Authenticate(accessToken string) error {
	ctx := context.Background()

	token := &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	c.spotifyClient = spotify.New(client)

	// Test authentication by getting user profile
	_, err := c.spotifyClient.CurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	return nil
}

// GetSpotifyClient returns the underlying Spotify client
func (c *Client) GetSpotifyClient() *spotify.Client {
	return c.spotifyClient
}

// RefreshToken refreshes the access token if needed
func (c *Client) RefreshToken(ctx context.Context) error {
	if !c.config.IsTokenExpired() {
		return nil
	}

	// TODO: For now, we'll require re-authentication but later it will use the refresh token
	return fmt.Errorf("token expired, please run 'spotify login' to re-authenticate")
}

// EnsureAuthenticated ensures the client is authenticated AND token is valid
func (c *Client) EnsureAuthenticated(ctx context.Context) error {
	if c.spotifyClient == nil {
		return fmt.Errorf("not authenticated, please run 'spotify login'")
	}

	// Check if token needs refresh
	if err := c.RefreshToken(ctx); err != nil {
		return err
	}

	return nil
}

// HandleAPIError just wraps Spotify API errors and provides user-friendly messages
func HandleAPIError(err error) error {
	if err == nil {
		return nil
	}

	// Check for specific Spotify API errors
	if spotifyErr, ok := err.(spotify.Error); ok {
		switch spotifyErr.Status {
		case 401:
			return fmt.Errorf("authentication failed, please run 'spotify login'")
		case 403:
			return fmt.Errorf("insufficient permissions, please check your Spotify app settings")
		case 404:
			return fmt.Errorf("resource not found")
		case 429:
			return fmt.Errorf("rate limit exceeded, please wait a moment and try again")
		case 500, 502, 503:
			return fmt.Errorf("spotify service is temporarily unavailable")
		default:
			return fmt.Errorf("spotify API error: %s", spotifyErr.Message)
		}
	}

	return fmt.Errorf("API error: %w", err)
}
