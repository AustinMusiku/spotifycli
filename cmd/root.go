/*
Copyright © 2025 AustinMusiku

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spotifycli",
	Short: "A Spotify controller for the terminal",
	Long: `Spotify CLI - Control your Spotify from the terminal.

A comprehensive command-line interface for Spotify that supports:
- Playback control (play, pause, next, previous, volume, shuffle, repeat)
- Search across all content types (tracks, albums, artists, playlists, shows, episodes)
- Library management (save/remove tracks, albums, shows)
- Device management (list and switch between devices)
- Queue management (view and add to queue)
- Modern OAuth2 PKCE authentication

Get started by running 'spotifycli login' to authenticate with Spotify.`,
}

// Adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "1.0.0"
}
