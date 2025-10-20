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

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for content",
	Long:  `Search for tracks, albums, artists, playlists, shows, and episodes.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		contentType, _ := cmd.Flags().GetString("type")
		return runSearch(strings.Join(args, " "), limit, contentType)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().IntP("limit", "l", 20, "Number of results to return")
	searchCmd.Flags().StringP("type", "t", "all", "Content type to search (track, album, artist, playlist, show, episode, all)")
}

func runSearch(query string, limit int, contentType string) error {
	_, client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	searchService := api.NewSearchService(client)
	ctx := context.Background()

	if contentType == "all" {
		// Search all content types
		results, err := searchService.Search(ctx, query, limit)
		if err != nil {
			return err
		}

		displaySearchResults(results)
	} else {
		// Search specific content type
		switch contentType {
		case "track":
			results, err := searchService.SearchTracks(ctx, query, limit)
			if err != nil {
				return err
			}
			displayTracks(results.Tracks.Tracks)

		case "album":
			results, err := searchService.SearchAlbums(ctx, query, limit)
			if err != nil {
				return err
			}
			displayAlbums(results.Albums.Albums)

		case "artist":
			results, err := searchService.SearchArtists(ctx, query, limit)
			if err != nil {
				return err
			}
			displayArtists(results.Artists.Artists)

		case "playlist":
			results, err := searchService.SearchPlaylists(ctx, query, limit)
			if err != nil {
				return err
			}
			displayPlaylists(results.Playlists.Playlists)

		case "show":
			results, err := searchService.SearchShows(ctx, query, limit)
			if err != nil {
				return err
			}
			displayShows(results.Shows.Shows)

		case "episode":
			results, err := searchService.SearchEpisodes(ctx, query, limit)
			if err != nil {
				return err
			}
			displayEpisodes(results.Episodes.Episodes)

		default:
			return fmt.Errorf("invalid content type: %s", contentType)
		}
	}

	return nil
}

func displaySearchResults(results *api.SearchResult) {
	if results.Tracks != nil && results.Tracks.Tracks != nil {
		fmt.Println("🎵 Tracks:")
		displayTracks(results.Tracks.Tracks)
		fmt.Println()
	}

	if results.Albums != nil && results.Albums.Albums != nil {
		fmt.Println("💿 Albums:")
		displayAlbums(results.Albums.Albums)
		fmt.Println()
	}

	if results.Artists != nil && results.Artists.Artists != nil {
		fmt.Println("🎤 Artists:")
		displayArtists(results.Artists.Artists)
		fmt.Println()
	}

	if results.Playlists != nil && results.Playlists.Playlists != nil {
		fmt.Println("📋 Playlists:")
		displayPlaylists(results.Playlists.Playlists)
		fmt.Println()
	}

	if results.Shows != nil && results.Shows.Shows != nil {
		fmt.Println("🎙️ Shows:")
		displayShows(results.Shows.Shows)
		fmt.Println()
	}

	if results.Episodes != nil && results.Episodes.Episodes != nil {
		fmt.Println("🎧 Episodes:")
		displayEpisodes(results.Episodes.Episodes)
		fmt.Println()
	}

}

func displayTracks(tracks []spotify.FullTrack) {
	for i, track := range tracks {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatTrack(track))
	}
}

func displayAlbums(albums []spotify.SimpleAlbum) {
	for i, album := range albums {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatAlbum(album))
	}
}

func displayArtists(artists []spotify.FullArtist) {
	for i, artist := range artists {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatArtist(artist))
	}
}

func displayPlaylists(playlists []spotify.SimplePlaylist) {
	for i, playlist := range playlists {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatPlaylist(playlist))
	}
}

func displayShows(shows []spotify.FullShow) {
	for i, show := range shows {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatShow(show.SimpleShow))
	}
}

func displayEpisodes(episodes []spotify.EpisodePage) {
	for i, episode := range episodes {
		fmt.Printf("  %d. %s\n", i+1, ui.FormatEpisode(episode))
	}
}
