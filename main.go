package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/lukas-hen/svtplay-translate/vtt"
)

func main() {

	argsWoProg := os.Args[1:]
	n_args := len(argsWoProg)
	expected_args := "episodes|translate|serve"

	if n_args == 0 {
		fmt.Printf("Invalid number of args passed (%d). Expected %s\n", n_args, expected_args)
		os.Exit(1)
	}

	command := argsWoProg[0]

	switch command {
	case "episodes":
		handleEpisodes()
	case "translate":
		subtitlePath := argsWoProg[1]
		handleTranslate(subtitlePath)
	case "serve":
		handleServe()
	default:
		fmt.Printf("Invalid command \"%s\" passed. Please pass %s\n", command, expected_args)
	}
}

func handleEpisodes() {
	RunUrlFetchingFlow()
}

func handleTranslate(subtitlePath string) {

	v := vtt.ParseFile(subtitlePath)
	cues := v.Cues
	translate_n := len(cues)

	translatedBuf := make([]string, translate_n)

	var wg sync.WaitGroup

	for i := 0; i < translate_n; i++ {
		wg.Add(1)
		go Translate(cues[i].TextWithoutTags(), translatedBuf, i, &wg)
	}

	wg.Wait()

	for idx, s := range translatedBuf {
		c := vtt.Cue{
			Id:      strconv.Itoa(idx),
			Timings: cues[idx].Timings,
			Text:    s,
		}

		c.WriteToSRTFile("./subtitle.srt")
	}

}

func handleServe() {
	RunServer()
}
