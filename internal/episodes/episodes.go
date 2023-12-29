package episodes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/lukas-hen/svtplay-translate/internal/utils"
)

const påSparetUrl = "https://svtplay.se/pa-sparet"
const svtApiUrl = "https://api.svt.se"

var allowedDomains = colly.AllowedDomains(
	"svtplay.se",
	"www.svtplay.se",
	"svtplay.se/pa-sparet",
)

func GetPåSpåretEpisodes() []string {
	episodes := crawlEpisodes()
	return utils.GetKeysFromMap[string, bool](episodes)
}

// Implementation here is so tied to SvtPlay & på spåret so it does not make sense to parameterize it.
// Returns a dict with the unique episodes and a boolean value = true.
// This is a common convention when using map as a set.
func crawlEpisodes() map[string]bool {
	// What is returned in the end.
	// Using map to only gather unique episodes.
	all_urls := make(map[string]bool)

	c := colly.NewCollector(allowedDomains)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		outbound_url := e.Attr("href")
		if strings.HasPrefix(outbound_url, "/video") && strings.Contains(outbound_url, "id=") {
			all_urls[outbound_url] = true
		}
	})

	log.Printf("Fetching episodes from: %s\n", påSparetUrl)

	err := c.Visit(påSparetUrl)
	if err != nil {
		log.Fatalf("Could not fetch episodes. ", err)
	}

	return all_urls
}

// Takes a på spåret url, constructs the api url and resolves the episode data.
func ResolveEpisode(videoPathUrl string) EpisodeResponse {

	u := videoPathToApiUrl(videoPathUrl)

	log.Printf("Requesting episode manifest from %s\n", u)
	res, err := http.Get(u)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var episodeResponse EpisodeResponse
	json.Unmarshal(body, &episodeResponse)

	return episodeResponse

}

// Extracts the first occurence of a dash-full manifest in a list of url filepaths.
func GetDashManifestUrl(videoList []VideoReference) string {

	for _, video := range videoList {
		if strings.HasSuffix(video.Url, "dash-full.mpd") {
			return video.Url
		}
	}

	log.Fatalln("Couldn't find a dash manifest.")
	return ""
}

func GetSubtitleUrl(subtitleList []SubtitleReference) string {

	for _, subtitle := range subtitleList {
		if strings.HasSuffix(subtitle.Url, "text-closed.vtt") {
			return subtitle.Url
		}
	}

	log.Fatalln("Couldn't find a dash manifest.")
	return ""
}

func DownloadSubtitleFile(subtitleUrl string, subtitlePath string) {

	res, err := http.Get(subtitleUrl)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	os.WriteFile(subtitlePath, body, 0644)
}

// Takes an svt video url from the site and constructs a url that can request video metadata from the api.
func videoPathToApiUrl(url string) string {
	// Path has lots of junk. E.g /video/eXYgnwk/pa-sparet/fre-3-nov-20-00?id=Kw74Q4o
	// We only need /video/<id> for the api.
	splitPath := strings.Split(url, "/")
	videoPath := "/" + splitPath[1] + "/" + splitPath[2]
	return svtApiUrl + videoPath
}

// func VideoPathToEpisodeName(url string) string {
// 	// Path has lots of junk. E.g /video/eXYgnwk/pa-sparet/fre-3-nov-20-00?id=Kw74Q4o
// 	// We only need /video/<id> for the api.
// 	splitPath := strings.Split(url, "/")
// 	return splitPath[4]
// }

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
