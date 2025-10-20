package cmd

import (
	"github.com/spf13/cobra"
)

// libraryCmd represents the library commands group
var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Manage your library",
	Long:  `Manage your saved tracks, albums, playlists, and shows.`,
}

func init() {
	rootCmd.AddCommand(libraryCmd)
}
