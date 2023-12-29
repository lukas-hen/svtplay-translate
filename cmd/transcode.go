package cmd

import (
	"github.com/lukas-hen/svtplay-translate/pkg/ffmpeg"
	"github.com/spf13/cobra"
)

var transcodeCmd = &cobra.Command{
	Use:   "transcode",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		s, _ := cmd.Flags().GetString("source")
		o, _ := cmd.Flags().GetString("out")
		e, _ := cmd.Flags().GetString("encoder")
		sub, _ := cmd.Flags().GetString("subtitles")

		if sub == "" {
			ffmpeg.Transcode(o, s, e)
		} else {
			ffmpeg.TranscodeWithSubs(o, s, e, sub)
		}

	},
}

func init() {
	rootCmd.AddCommand(transcodeCmd)
	transcodeCmd.PersistentFlags().StringP("source", "s", "", "Path to source file to transcode.")
	transcodeCmd.PersistentFlags().StringP("out", "o", "", "Path to transcoded output file.")
	transcodeCmd.PersistentFlags().StringP("encoder", "e", "libx264", "The video encoder to use.")
	transcodeCmd.PersistentFlags().String("subtitles", "", "Path to subtitles (.srt | .ass) file.")
	transcodeCmd.MarkPersistentFlagRequired("source")
	transcodeCmd.MarkPersistentFlagRequired("out")
}
