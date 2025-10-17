package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/AustinMusiku/spotifycli/internal/api"
	"github.com/AustinMusiku/spotifycli/internal/auth"
	"github.com/AustinMusiku/spotifycli/internal/config"
	"github.com/AustinMusiku/spotifycli/internal/ui"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Spotify",
	Long:  `Authenticate with Spotify using OAuth2 PKCE flow. This will open your browser for authorization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLogin()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if already authenticated
	if cfg.IsAuthenticated() && !cfg.IsTokenExpired() {
		ui.PrintInfo("Already authenticated with Spotify")
		return nil
	}

	// Get client ID if not set
	if cfg.ClientID == "" {
		ui.PrintInfo("Please enter your Spotify Client ID:")
		fmt.Print("Client ID: ")

		var clientID string
		fmt.Scanln(&clientID)

		if clientID == "" {
			return fmt.Errorf("client ID is required")
		}

		cfg.ClientID = clientID
	}

	// Create PKCE auth handler
	pkceAuth := auth.NewPKCEAuth(cfg.ClientID, cfg.RedirectURI)

	// Fire up callback server
	server := auth.NewCallbackServer("8080")
	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start callback server: %w", err)
	}
	defer server.Stop()

	// Get authorization URL
	authURL := pkceAuth.GetAuthURL()

	ui.PrintInfo("Opening browser for authentication...")
	ui.PrintInfo(fmt.Sprintf("If the browser doesn't open automatically, visit: %s", authURL))

	if err := openBrowser(authURL); err != nil {
		ui.PrintWarning(fmt.Sprintf("Failed to open browser: %v", err))
		ui.PrintInfo(fmt.Sprintf("Please visit: %s", authURL))
	}

	ui.PrintInfo("Waiting for authentication...")
	code, state, err := server.WaitForCallback(5 * time.Minute)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Verify state
	if state != pkceAuth.State {
		return fmt.Errorf("invalid state parameter")
	}

	// Exchange code for token
	ctx := context.Background()
	token, err := pkceAuth.ExchangeCode(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}

	cfg.SetTokens(token.AccessToken, token.RefreshToken, token.TokenType, 3600) // Default 1 hour
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}

	// Test authentication
	client := api.NewClient(cfg)
	if err := client.Authenticate(token.AccessToken); err != nil {
		return fmt.Errorf("authentication test failed: %w", err)
	}

	ui.PrintSuccess("Successfully authenticated with Spotify!")
	return nil
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
