package cmd

import (
	"log"
	"os"

	"github.com/lukas-hen/svtplay-translate/internal/episodes"
	"github.com/lukas-hen/svtplay-translate/internal/utils"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var episodesCmd = &cobra.Command{
	Use:   "episodes",
	Short: "Tool for downloading på spåret episodes.",
	Long: `The episodes commmand prompts you for which episode of På Spåret to download.
	It then proceeds to download the episode DASH manifest & swedish subtitles.`,
	Run: func(cmd *cobra.Command, args []string) {

		s, _ := cmd.Flags().GetString("subtitle_path")
		v, _ := cmd.Flags().GetString("video_url_path")

		availableEpisodes := episodes.GetPåSpåretEpisodes()
		episode := utils.PromptUserDecision[string](availableEpisodes, "Which episode do you want to download & translate?")

		res := episodes.ResolveEpisode(episode)

		manifestUrl := episodes.GetDashManifestUrl(res.VideoReferences)
		subsUrl := episodes.GetSubtitleUrl(res.SubtitleReferences)

		vf, err := os.Create(v)
		if err != nil {
			log.Fatalln(err)
		}
		defer vf.Close()
		vf.WriteString(manifestUrl + "\n")

		episodes.DownloadSubtitleFile(subsUrl, s)

	},
}

func init() {
	rootCmd.AddCommand(episodesCmd)
	episodesCmd.PersistentFlags().StringP("subtitle_path", "s", "", "Path of where to store subtitles.")
	episodesCmd.PersistentFlags().StringP("video_url_path", "v", "", "Path of where to store video manifest url.")
	episodesCmd.MarkPersistentFlagRequired("subtitle_path")
	episodesCmd.MarkPersistentFlagRequired("video_url_path")
}
