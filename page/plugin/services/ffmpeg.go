package services

import (
	"bytes"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"strconv"
)

func GetNFrameBlurred(in io.Reader, frameNum int, blurRatio int) (*bytes.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("pipe:").
		WithInput(in).
		Output("pipe:", ffmpeg.KwArgs{"vframes": frameNum, "f": "image2", "vf": "gblur=sigma=" + strconv.Itoa(blurRatio)}).
		WithOutput(buf).
		Run()
	return bytes.NewReader(buf.Bytes()), err
}

func GetNFrame(in io.Reader, frameNum int) (*bytes.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("pipe:").
		WithInput(in).
		Output("pipe:", ffmpeg.KwArgs{"vframes": frameNum, "f": "image2"}).
		WithOutput(buf).
		Run()
	return bytes.NewReader(buf.Bytes()), err
}

//ffmpeg -i input.mp4 -c:v libx265 -preset ultrafast -crf 28 -c:a aac -b:a 250k output.mp4
//-tag:v hvc1 -c:a eac3

func CompressVideo(in io.Reader) (*bytes.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("pipe:").
		WithInput(in).
		Output("pipe:", ffmpeg.KwArgs{
			"movflags": "frag_keyframe+empty_moov",
			"format":   "mp4",
			"c:v":      "libx265",
			"preset":   "ultrafast",
			"crf":      "28",
			"c:a":      "aac",
			"b:a":      "250k",
			"tag:v":    "hvc1",
		}).
		WithOutput(buf).
		Run()
	return bytes.NewReader(buf.Bytes()), err
}
