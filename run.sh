#!/bin/bash

export OPENAI_API_KEY=<KEY_HERE>

# Serving the video content requires superuser. Run me as sudo.

SUBTITLE_PATH=./tmp/subtitles.vtt
VIDEO_PATH=./tmp/video.url
TRANSLATED_PATH=./tmp/translated.srt
VIDEO_W_SUBS=./tmp/video_subbed.mp4
IFACE_NAME=en0

#./bin/svtplay-translate episodes -s $SUBTITLE_PATH -v $VIDEO_PATH

./bin/svtplay-translate translate -s $SUBTITLE_PATH -o $TRANSLATED_PATH

./bin/svtplay-translate transcode -s $(cat $VIDEO_PATH) -o $VIDEO_W_SUBS --subtitles $TRANSLATED_PATH 

#./bin/svtplay-translate serve -f $VIDEO_W_SUBS -i $IFACE_NAME

