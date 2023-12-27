package main

import (
	"fmt"
	"os"
	"strconv"
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

	// 10 req per sec
	rl := NewRateLimiter(5)

	rl.Run()

	req_id := 1
	for {
		if req_id > 20 {
			break
		}
		req := "Req: " + strconv.Itoa(req_id)
		rl.Send(req)
		req_id++
	}
	rl.Wait()

	//v := vtt.ParseFile(subtitlePath)
	//cues := v.Cues
	// strbuf := ""

	// // Translate one chunk of 5 cues.
	// for i := 0; i < 5; i++ {
	// 	strbuf += cues[i].TextWithoutTags() + "\n\n"
	// }

	// res := Translate(strbuf)

	// fmt.Println(res)
	//v.WriteSrtFile("./subtitles.srt")

}

func handleServe() {
	RunServer()
}
