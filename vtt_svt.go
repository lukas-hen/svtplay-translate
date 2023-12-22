package main

// For parsing svt flavored .vtt subtitles.

// Example entry. each of these are split with newline.

// 82962cbb-7a7f-49e9-bf73-4691c13c7ea3
// 00:26:21.240 --> 00:26:26.920 align:left position:18%
// <c.teletext>Familjerna vÃ¤ntar spÃ¤nt utanfÃ¶r. Men
// de senaste Ã¥ren har protesterna Ã¶katâ€“</c>

func ScanWTTHead() {}

func ScanStyle() {}

func ScanSubtitle() {}

func ScanSubUuid() {}

func ScanSubMeta() { /* Duration, Alignment & Position} */ }

func ScanSubText() {}

type VTTEntry struct {
	Uuid             string
	SubtitleMetadata SubtitleMetadata
	Subtitle         Subtitle
}

type Subtitle struct {
	Text  string
	Class string
}

type SubtitleMetadata struct {
	Duration  Duration
	Alignment string
	Position  string
}

type Duration struct {
	From string // should be changed to time.
	To   string
}
