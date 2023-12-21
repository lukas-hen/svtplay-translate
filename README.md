# SvtPlay English Subtitle Streaming

This CLI tool fetches a stream and subtitles from SVT Play, transcodes the content, translates subtitles, and hosts it on your computer's IP. Enjoy the stream by accessing your computer's IP through a local network browser or a media player that supports .mp4 playback.

## How it Works

1. **Interacting with SVT Site:**
    - Go navigates the SVT site, prompts for the episode, and retrieves the [DASH](https://en.wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP) manifest (manifest.mpd) and subtitles.vtt through the SVT CDN.

2. **Parallel Processes:**
    - FFMPEG decodes the stream using the manifest into /videos/<video_name>.mp4.
    - Simultaneously, a Go program converts subtitles from .vtt to .srt, translates them, and stores them in /subtitles/<sub_name>.srt.

3. **Embedding Subtitles:**
    - FFMPEG embeds the subtitles with the video & audio, saving the result to /serve/<video_w_subs>.mp4.

4. **Server Launch:**
    - A straightforward Go file server is initiated to serve the specific video. (Streaming enhancements may be considered in future iterations.)