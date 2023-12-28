package cmd

import (
	"github.com/spf13/cobra"
)

// Takes a .VTT subtitle input file, translates it and writes the translated subtitles
// to a .SRT file that ffmpeg can burn in.
var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cmd.Flags().GetString("sourceVtt")
		cmd.Flags().GetString("outSrt")
		// TODO

	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	transcodeCmd.PersistentFlags().StringP("source", "s", "", "Path to source .VTT file to translate.")
	transcodeCmd.PersistentFlags().StringP("out", "o", "", "Path to output .SRT file destination.")
	transcodeCmd.MarkPersistentFlagRequired("source")
	transcodeCmd.MarkPersistentFlagRequired("out")
}
