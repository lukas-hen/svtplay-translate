package svt

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/lukas-hen/svtplay-translate/internal/utils"
)

// Writes fetched urls into a target file /
func RunUrlFetchingFlow() {

	// This application has no stability requirements.
	// Most issues throw panics.

	// Init
	base_api_url := "https://api.svt.se"
	client := &http.Client{}

	// Start of control flow
	episodes := CrawlPSEpisodes()
	episode_urls := utils.GetKeysFromMap[string, bool](episodes)
	chosen_url := utils.PromptUserDecision[string](episode_urls, "Which episode do you want to translate?")

	// Path has lots of junk. E.g /video/eXYgnwk/pa-sparet/fre-3-nov-20-00?id=Kw74Q4o
	// We only need /video/<id> for the api.
	split_path := strings.Split(chosen_url, "/")
	video_path := "/" + split_path[1] + "/" + split_path[2]
	full_req_path := base_api_url + video_path

	log.Printf("Requesting cdn addresses from %s\n", full_req_path)
	body := MakeSyncGet(client, full_req_path)

	var episodeResponse EpisodeResponse
	json.Unmarshal([]byte(body), &episodeResponse)

	allVideoObjects := episodeResponse.VideoReferences
	allSubtitleObjects := episodeResponse.SubtitleReferences

	videos := make([]string, len(allVideoObjects))
	for idx, video := range allVideoObjects {
		videos[idx] = video.Url
	}

	selectedVideo := utils.PromptUserDecision[string](videos, "Select which streaming protocol to get:")

	subtitles := make([]string, len(allSubtitleObjects))
	for idx, subtitle := range allSubtitleObjects {
		subtitles[idx] = subtitle.Url
	}

	selectedSubtitle := utils.PromptUserDecision[string](subtitles, "Select which subtitles to get:")

	// For future container purposes these should be set in the env.
	tmp_path := filepath.Join(".", "tmp")
	err := os.MkdirAll(tmp_path, os.ModePerm)
	check(err)

	videoPath := tmp_path + "/video-" + time.Now().Format(time.RFC3339) + ".txt"
	subtitlePath := tmp_path + "/subtitle-" + time.Now().Format(time.RFC3339) + ".txt"

	log.Printf("Writing video url to: %s\n", videoPath)
	err = os.WriteFile(videoPath, []byte(selectedVideo+"\n"), os.ModeAppend)
	check(err)

	log.Printf("Writing subtitle url to: %s\n", subtitlePath)
	err = os.WriteFile(subtitlePath, []byte(selectedSubtitle+"\n"), os.ModeAppend)
	check(err)
}

// Implementation here is so tied to SvtPlay so it does not make sense to parameterize it.
// Returns a dict with the unique episodes and a boolean value = true.
// This is a common convention when using map as a set.
func CrawlPSEpisodes() map[string]bool {
	// What is returned in the end.
	// Using map to only gather unique episodes.
	all_urls := make(map[string]bool)

	c := colly.NewCollector(
		colly.AllowedDomains(
			"svtplay.se",
			"www.svtplay.se",
			"svtplay.se/pa-sparet",
		),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		outbound_url := e.Attr("href")
		if strings.HasPrefix(outbound_url, "/video") && strings.Contains(outbound_url, "id=") {
			all_urls[outbound_url] = true
		}
	})

	pa_sparet_url := "https://svtplay.se/pa-sparet"

	log.Printf("Fetching episodes from: %s\n", pa_sparet_url)

	err := c.Visit(pa_sparet_url)
	check(err)

	return all_urls
}

// Makes Simple Blocking Get Request
func MakeSyncGet(client *http.Client, url string) []byte {

	req, err := http.NewRequest("GET", url, nil)
	check(err)

	res, err := client.Do(req)
	check(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	check(err)

	return body
}

type EpisodeResponse struct {
	SvtId              string              `json:"svtId"`
	ProgramTitle       string              `json:"programTitle"`
	EpisodeTitle       string              `json:"episodeTitle"`
	VideoReferences    []VideoReference    `json:"videoReferences"`
	SubtitleReferences []SubtitleReference `json:"subtitleReferences"`
}

type VideoReference struct {
	Url      string `json:"url"`
	Redirect string `json:"redirect"`
	Resolve  string `json:"resolve"`
	Format   string `json:"format"`
}

type SubtitleReference struct {
	Url          string `json:url`
	Format       string `json:format`
	Language     string `json:language`
	LanguageName string `json:languageName`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
