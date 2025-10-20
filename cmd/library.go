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

var libraryAlbumsCmd = &cobra.Command{
	Use:   "albums",
	Short: "List your saved albums",
	Long:  `List all your saved albums.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		return runLibraryAlbums(limit)
	},
}

var libraryTracksCmd = &cobra.Command{
	Use:   "tracks",
	Short: "List your saved tracks",
	Long:  `List all your saved tracks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		return runLibraryTracks(limit)
	},
}

var libraryShowsCmd = &cobra.Command{
	Use:   "shows",
	Short: "List your saved shows",
	Long:  `List all your saved shows (podcasts).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		return runLibraryShows(limit)
	},
}

func init() {
	rootCmd.AddCommand(libraryCmd)
	libraryCmd.AddCommand(libraryPlaylistsCmd)
	libraryCmd.AddCommand(libraryAlbumsCmd)
	libraryCmd.AddCommand(libraryTracksCmd)
	libraryCmd.AddCommand(libraryShowsCmd)

	// Add limit flags
	libraryPlaylistsCmd.Flags().IntP("limit", "l", 50, "Number of results to return")
	libraryAlbumsCmd.Flags().IntP("limit", "l", 50, "Number of results to return")
	libraryTracksCmd.Flags().IntP("limit", "l", 50, "Number of results to return")
	libraryShowsCmd.Flags().IntP("limit", "l", 50, "Number of results to return")
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

func runLibraryAlbums(limit int) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	libraryService := api.NewLibraryService(client)
	ctx := context.Background()

	albums, err := libraryService.GetSavedAlbums(ctx, limit)
	if err != nil {
		return err
	}

	if len(albums.Albums) == 0 {
		ui.PrintInfo("No saved albums found")
		return nil
	}

	fmt.Println("💿 Your Saved Albums:")
	for i, album := range albums.Albums {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatAlbum(album.SimpleAlbum))
	}

	return nil
}

func runLibraryTracks(limit int) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	libraryService := api.NewLibraryService(client)
	ctx := context.Background()

	tracks, err := libraryService.GetSavedTracks(ctx, limit)
	if err != nil {
		return err
	}

	if len(tracks.Tracks) == 0 {
		ui.PrintInfo("No saved tracks found")
		return nil
	}

	fmt.Println("🎵 Your Saved Tracks:")
	for i, track := range tracks.Tracks {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatTrack(track.FullTrack))
	}

	return nil
}

func runLibraryShows(limit int) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	libraryService := api.NewLibraryService(client)
	ctx := context.Background()

	shows, err := libraryService.GetSavedShows(ctx, limit)
	if err != nil {
		return err
	}

	if len(shows.Shows) == 0 {
		ui.PrintInfo("No saved shows found")
		return nil
	}

	fmt.Println("🎙️ Your Saved Shows:")
	for i, show := range shows.Shows {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatShow(show.SimpleShow))
	}

	return nil
}
