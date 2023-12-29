package cmd

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/lukas-hen/svtplay-translate/pkg/translation"
	"github.com/lukas-hen/svtplay-translate/pkg/vtt"
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
		openaiApiKey := os.Getenv("OPENAI_API_KEY")

		webvtt := vtt.ParseFile(src)
		cues := webvtt.Cues
		translateN := len(cues)

		translator := translation.NewOpenaiTranslator(openaiApiKey, "Swedish", "English")

		// Pre-allocate buffer. One goroutine will write to one index.
		// Since only one goroutine operates on one index this should be memory safe.
		translatedBuf := make([]string, translateN)

		var wg sync.WaitGroup

		for i := 0; i < translateN; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()

				s, err := translator.Translate(cues[n].TextWithoutTags())
				if err != nil {
					log.Printf("Non rate-limiting error translating. Err: %s", err)
					return
				}

				translatedBuf[n] = s
			}(i)
		}

		wg.Wait()

		var translatedCues []*vtt.Cue

		for idx, s := range translatedBuf {
			c := &vtt.Cue{
				Id:      strconv.Itoa(idx),
				Timings: cues[idx].Timings,
				Text:    s,
			}

			translatedCues = append(translatedCues, c)
		}

		f, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		for _, cue := range translatedCues {
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
