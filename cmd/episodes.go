package cmd

import (
	"github.com/lukas-hen/svtplay-translate/internal/svt"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "episodes",
	Short: "The episodes commmand prompts you for which episode of På Spåret to download.",
	Long:  `The episodes commmand prompts you for which episode of På Spåret to download.`,
	Run: func(cmd *cobra.Command, args []string) {
		svt.RunUrlFetchingFlow()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
