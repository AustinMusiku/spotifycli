package cmd

import (
	"fmt"

	"github.com/AustinMusiku/spotifycli/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of spotifycli",
	Long:  `Show the current version, commit hash, and build date of spotifycli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("spotifycli build\n\nVersion: %s\nCommit: %s\nDate: %s\n", version.Version, version.Commit, version.BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Aliases = []string{"v"}
}
