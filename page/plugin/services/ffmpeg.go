package services

import (
	"bufio"
	"bytes"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
	"strconv"
	"strings"
)

func GetNFrameBlurred(fname string, frameNum int, blurRatio int) (*bytes.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("file:"+fname).
		Output(
			"pipe:", ffmpeg.KwArgs{"vframes": frameNum, "f": "image2", "vf": "gblur=sigma=" + strconv.Itoa(blurRatio)},
		).
		WithOutput(buf).
		WithErrorOutput(os.Stdout).
		Run()
	return bytes.NewReader(buf.Bytes()), err
}

func GetNFrame(fname string, frameNum int) (*bytes.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("file:"+fname).
		Output("pipe:", ffmpeg.KwArgs{"vframes": frameNum, "f": "image2"}).
		WithOutput(buf).
		WithErrorOutput(os.Stdout).
		Run()
	return bytes.NewReader(buf.Bytes()), err
}

//ffmpeg -i input.mp4 -c:v libx265 -preset ultrafast -crf 28 -c:a aac -b:a 250k output.mp4
//-tag:v hvc1 -c:a eac3

func CompressVideo(in io.Reader) (*bytes.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("pipe:").
		WithInput(in).
		Output(
			"pipe:", ffmpeg.KwArgs{
				"movflags": "frag_keyframe+empty_moov",
				"format":   "mp4",
				"c:v":      "libx265",
				"preset":   "ultrafast",
				"crf":      "28",
				"c:a":      "aac",
				"b:a":      "250k",
				"tag:v":    "hvc1",
			},
		).
		WithOutput(buf).
		Run()
	return bytes.NewReader(buf.Bytes()), err
}

func GetAudioFromVideo(fname string) (*bytes.Reader, string, error) {
	//audio := bytes.NewBuffer(nil)

	outputFile := strings.ReplaceAll(fname, ".mp4", ".mp3")

	outputFile = "audio_" + outputFile

	err := ffmpeg.Input("file:"+fname).Output(
		outputFile, ffmpeg.KwArgs{"q:a": 0, "map": "a"},
	).Run()

	audioOs, err := os.Open(outputFile)

	stat, err := audioOs.Stat()
	if err != nil {
		return nil, "", err
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(audioOs).Read(bs)
	if err != nil && err != io.EOF {
		return nil, "", err
	}

	return bytes.NewReader(bs), outputFile, err
}
