package main

import (
	"github.com/lukas-hen/svtplay-translate/cmd"
	"github.com/lukas-hen/svtplay-translate/pkg/server"
	"github.com/lukas-hen/svtplay-translate/pkg/svt"
)

func main() {
	cmd.Execute()
}

func handleEpisodes() {
	svt.RunUrlFetchingFlow()
}

func handleServe() {
	server.RunServer()
}
