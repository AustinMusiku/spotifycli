package cmd

import (
	"context"
	"fmt"

	"github.com/AustinMusiku/spotifycli/internal/api"
	"github.com/AustinMusiku/spotifycli/internal/ui"
	"github.com/spf13/cobra"
)

// libraryCmd represents the library commands group
var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Manage your library",
	Long:  `Manage your saved tracks, albums, playlists, and shows.`,
}

var libraryPlaylistsCmd = &cobra.Command{
	Use:   "playlists",
	Short: "List your playlists",
	Long:  `List all your saved playlists.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		return runLibraryPlaylists(limit)
	},
}

func init() {
	rootCmd.AddCommand(libraryCmd)
	libraryCmd.AddCommand(libraryPlaylistsCmd)

	// Add limit flags
	libraryPlaylistsCmd.Flags().IntP("limit", "l", 50, "Number of results to return")
}

func runLibraryPlaylists(limit int) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	libraryService := api.NewLibraryService(client)
	ctx := context.Background()

	playlists, err := libraryService.GetUserPlaylists(ctx, limit)
	if err != nil {
		return err
	}

	if len(playlists.Playlists) == 0 {
		ui.PrintInfo("No playlists found")
		return nil
	}

	fmt.Println("📋 Your Playlists:")
	for i, playlist := range playlists.Playlists {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatPlaylist(playlist))
	}

	return nil
}
