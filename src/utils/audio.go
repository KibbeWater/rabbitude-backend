package utils

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"text/template"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func MergeAudioBuffer(audioBuf [][]byte) []byte {
	tempDir := path.Join(os.TempDir(), "rabbit")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	var inputFilePaths []string

	// Write each audio buffer to a file and collect the file paths
	for i, audio := range audioBuf {
		filePath := path.Join(tempDir, strconv.Itoa(i)+".wav")
		err := os.WriteFile(filePath, audio, os.ModePerm)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return nil
		}
		inputFilePaths = append(inputFilePaths, filePath)
	}

	// Create a file list for FFmpeg concat demuxer
	listFilePath := path.Join(tempDir, "filelist.txt")
	listFile, err := os.Create(listFilePath)
	if err != nil {
		fmt.Println("Error creating file list:", err)
		return nil
	}
	defer listFile.Close()

	// Write the file paths to the list file in the required format
	tmpl := template.Must(template.New("filelist").Parse("file '{{.}}'\n"))
	for _, filePath := range inputFilePaths {
		tmpl.Execute(listFile, filePath)
	}

	// Define output file path
	outputPath := path.Join(tempDir, "output.wav")

	// Run FFmpeg with the concat demuxer
	err = ffmpeg.Input(listFilePath, ffmpeg.KwArgs{"f": "concat", "safe": "0"}).
		Output(outputPath, ffmpeg.KwArgs{"c:a": "pcm_s16le"}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	if err != nil {
		fmt.Println("Error running ffmpeg:", err)
		return nil
	}

	// Read and return the merged output audio file
	outputAudio, err := os.ReadFile(outputPath)
	if err != nil {
		fmt.Println("Error reading output file:", err)
		return nil
	}

	return outputAudio
}
