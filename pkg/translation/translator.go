package translation

import (
	"strconv"
	"sync"

	"github.com/cenkalti/backoff"
	"github.com/lukas-hen/svtplay-translate/internal/vtt"
)

func translate(cueText string, t Translator, translatedBuf []string, i int, wg *sync.WaitGroup) {

	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
	defer ticker.Stop()

	// Instead of rate-limiting, keep retrying with an exponential backoff.
	for range ticker.C {
		translated, err := t.Translate(cueText)
		if err == nil {
			translatedBuf[i] = translated
			break
		}
	}

	wg.Done()
}

// Parallel implementation for translating cues.
func ParTranslateCues(cues []*vtt.Cue, t Translator) []*vtt.Cue {

	translateN := len(cues)

	// Pre-allocate buffer. One goroutine will write to one index.
	// Since only one goroutine operates on one index this should be memory safe.
	translatedBuf := make([]string, translateN)

	var wg sync.WaitGroup

	for i := 0; i < translateN; i++ {
		wg.Add(1)
		go translate(cues[i].TextWithoutTags(), t, translatedBuf, i, &wg)
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

	return translatedCues
}
