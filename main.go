package main

import (
	"fmt"
	"os"

	"github.com/lukas-hen/svtplay-translate/vtt"
)

func main() {

	argsWoProg := os.Args[1:]
	n_args := len(argsWoProg)
	expected_args := "episodes|translate|serve"

	if n_args > 1 || n_args == 0 {
		fmt.Printf("Invalid number of args passed (%d). Expected %s\n", n_args, expected_args)
		os.Exit(1)
	}

	command := argsWoProg[0]

	switch command {
	case "episodes":
		handleEpisodes()
	case "translate":
		handleTranslate()
	case "serve":
		handleServe()
	default:
		fmt.Printf("Invalid command \"%s\" passed. Please pass %s\n", command, expected_args)
	}
}

func handleEpisodes() {
	RunUrlFetchingFlow()
}

func handleTranslate() {
	v := vtt.ParseFile("./subtitles.vtt")
	cues := v.Cues
	for _, c := range cues {
		c.Text = c.Text + "123"
	}

	err := v.WriteSrtFile("./subtitles.srt")
	if err != nil {
		panic(err)
	}
}

func handleServe() {
	RunServer()
}
