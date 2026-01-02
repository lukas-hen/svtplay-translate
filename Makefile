build:
	go build -o bin ./...

clean:
	rm out

trans:
	ffmpeg -i https://ed0-as2119.cdn.svt.se/d0/se/20251209/e2fa79af-9fc2-4862-8ba7-9d93f26f42a1/dash-full.mpd -filter_complex subtitles=./tmp/translated.srt -c:a copy -c:v libx264 -preset ultrafast ./tmp/video_subbed.mp4