package ffmpeg

import (
	"log"
	"os"
	"os/exec"
)

// Transcodes video but allows passing subtitles to burn them in during the
// downloading/encoding, otherwise ffmpeg has to go through the entire video twice.
// Just a wrapper, requires ffmpeg installation.
func TranscodeWithSubs(dstPath string, srcPath string, videoEncoder string, subtitlePath string) {

	// srcPath can also be a url. I.e to a dash manifest.
	// hardcoded ultrafast preset for now.
	cmd := exec.Command("ffmpeg", "-i", srcPath, "-filter_complex", "subtitles="+subtitlePath, "-c:a", "copy", "-c:v", videoEncoder, "-preset", "ultrafast", dstPath)

	// Pipe command output to stdout
	// Ffmpeg writes info to stderr, not stdout.
	cmd.Stderr = os.Stdout

	// Run still runs the command and waits for completion
	// but the output is instantly piped to Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("could not run command: ", err)
	}

}

// Wrapper for the ffmpeg shell command. Transcodes video stream and copies the other streams to a local multimedia container.
// srcPath can be a url or a filepath, if passed a manifest it will be resolved.
// Just a wrapper, requires ffmpeg installation.
func Transcode(dstPath string, srcPath string, videoEncoder string) {

	// srcPath can also be a url. I.e to a dash manifest.
	// hardcoded ultrafast preset for now.
	cmd := exec.Command("ffmpeg", "-i", srcPath, "-c:a", "copy", "-c:v", videoEncoder, "-preset", "ultrafast", dstPath)

	// Pipe command output to stdout
	// Ffmpeg writes info to stderr, not stdout.
	cmd.Stderr = os.Stdout

	// Run still runs the command and waits for completion
	// but the output is instantly piped to Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("could not run command: ", err)
	}

}
