package cmd

import (
	"log"
	"os"

	"github.com/lukas-hen/svtplay-translate/internal/vtt"
	"github.com/lukas-hen/svtplay-translate/pkg/translation"
	"github.com/spf13/cobra"
)

// Takes a .VTT subtitle input file, translates it and writes the translated subtitles
// to a .SRT file that ffmpeg can burn in.
var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		src, _ := cmd.Flags().GetString("source")
		dst, _ := cmd.Flags().GetString("out")

		webvtt := vtt.ParseFile(src)
		translator := translation.NewOpenaiTranslator()

		translated := translation.ParTranslateCues(webvtt.Cues, translator)

		f, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		for _, cue := range translated {
			f.WriteString(cue.ToSRT())
		}
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.PersistentFlags().StringP("source", "s", "", "Path to source .VTT file to translate.")
	translateCmd.PersistentFlags().StringP("out", "o", "", "Path to output .SRT file destination.")
	translateCmd.MarkPersistentFlagRequired("source")
	translateCmd.MarkPersistentFlagRequired("out")
}
