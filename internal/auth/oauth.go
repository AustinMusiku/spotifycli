package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	authURL  = "https://accounts.spotify.com/authorize"
	tokenURL = "https://accounts.spotify.com/api/token"
)

// Required scopes for the Spotify CLI
var scopes = []string{
	"user-read-playback-state",
	"user-modify-playback-state",
	"user-read-currently-playing",
	"user-library-read",
	"user-library-modify",
	"playlist-read-private",
	"playlist-modify-public",
	"playlist-modify-private",
	"user-read-email",
	"user-read-private",
	"user-read-recently-played",
	"user-top-read",
	"user-follow-read",
	"user-follow-modify",
	"streaming",
}

type PKCEAuth struct {
	ClientID      string
	RedirectURI   string
	State         string
	CodeVerifier  string
	CodeChallenge string
}

func NewPKCEAuth(clientID, redirectURI string) *PKCEAuth {
	state := generateRandomString(32)
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	return &PKCEAuth{
		ClientID:      clientID,
		RedirectURI:   redirectURI,
		State:         state,
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
	}
}

// GetAuthURL returns the authorization URL for the user to visit
func (a *PKCEAuth) GetAuthURL() string {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", a.ClientID)
	params.Add("scope", strings.Join(scopes, " "))
	params.Add("redirect_uri", a.RedirectURI)
	params.Add("state", a.State)
	params.Add("code_challenge_method", "S256")
	params.Add("code_challenge", a.CodeChallenge)

	return fmt.Sprintf("%s?%s", authURL, params.Encode())
}

// ExchangeCode exchanges the authorization code for tokens
func (a *PKCEAuth) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	config := &oauth2.Config{
		ClientID: a.ClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: a.RedirectURI,
	}

	// Create a custom HTTP client for PKCE
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Exchange code for token with PKCE
	token, err := config.Exchange(
		context.WithValue(ctx, oauth2.HTTPClient, client),
		code,
		oauth2.SetAuthURLParam("code_verifier", a.CodeVerifier),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[mathrand.Intn(len(charset))]
	}
	return string(b)
}

// generateCodeVerifier generates a code verifier for PKCE
func generateCodeVerifier() string {
	return generateRandomString(128)
}

// generateCodeChallenge generates a code challenge from the verifier
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
