package translation

import (
	"strconv"
	"sync"

	"github.com/lukas-hen/svtplay-translate/internal/vtt"
)

func TranslateVTTtoSRT(subtitlePath string) {

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
