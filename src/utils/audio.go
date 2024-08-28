package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func MergeAudioBuffer(audioBuf [][]byte) []byte {
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)

	ffmpeg_cmd := "ffmpeg"
	for i, audio := range audioBuf {
		path := path.Join(tempDir, string(i)+".wav")
		os.WriteFile(path, audio, 0777)
		ffmpeg_cmd += " -i " + path
	}

	outputDir := path.Join(tempDir, "output.wav")
	ffmpeg_cmd += " " + outputDir

	cmd := exec.Command(ffmpeg_cmd)
	cmd_out, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// Print the output
	fmt.Println(string(cmd_out))

	audio, err := os.ReadFile(outputDir)
	if err != nil {
		return nil
	}

	return audio
}
