package actions

import (
	"bytes"
	"context"
	"fmt"
	mediaActions "github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/page/plugin/services"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type VideoUseCase struct {
	medias *mediaActions.Medias
	logger *zap.SugaredLogger
}

func NewVideoUseCase(
	medias *mediaActions.Medias,
	logger *zap.SugaredLogger,
) *VideoUseCase {
	return &VideoUseCase{
		medias: medias,
		logger: logger,
	}
}

func (uc VideoUseCase) Create(request CreateVideoRequest) (*CreateVideoResponse, error) {
	ctx := context.Background()

	uc.logger.Debug("Creating video samples", "fileName", request.Filename)
	videoSamples, err := uc.createVideoSamples(request.File)
	if err != nil {
		return nil, err
	}

	uc.logger.Debug("Uploading medias", "fileName", request.Filename)
	fileExt := filepath.Ext(request.Filename)
	fileTitle := strings.TrimSuffix(request.Filename, fileExt)
	ids, err := uc.medias.UploadList(ctx, mediaActions.CreateMediasList{
		FolderId: request.FolderId,
		Files: []mediaActions.MediaFile{
			{
				Filename: request.Filename,
				Size:     request.Size,
				File:     request.File,
			},
			//{
			//	Filename: fileTitle + "_compressed" + fileExt,
			//	Size:     videoSamples.CompressedVideo.Size(),
			//	File:     videoSamples.CompressedVideo,
			//},
			{
				Filename: fileTitle + "_preview" + ".jpg",
				Size:     videoSamples.Preview.Size(),
				File:     videoSamples.Preview,
			},
			{
				Filename: fileTitle + "_preview_blurred" + ".jpg",
				Size:     videoSamples.PreviewBlurred.Size(),
				File:     videoSamples.PreviewBlurred,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(ids) != 3 {
		return nil, fmt.Errorf("error uploading medias")
	}

	return &CreateVideoResponse{
		VideoId: ids[0],
		//VideoLiteId:      ids[1],
		PreviewId:        ids[1],
		PreviewBlurredId: ids[2],
	}, nil
}

func (uc VideoUseCase) createVideoSamples(video io.ReadSeeker) (*VideoSamples, error) {
	f, fname, err := createTempVideoFile(video)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	uc.logger.Debug("Creating video preview")
	preview, err := services.GetNFrame(fname, 1)
	if err != nil {
		return nil, err
	}
	_, err = video.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	uc.logger.Debug("Creating video blurred preview")
	blurred, err := services.GetNFrameBlurred(fname, 1, 35)
	if err != nil {
		return nil, err
	}
	_, err = video.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	//uc.logger.Debug("Creating compressed video")
	//compressed, err := services.CompressVideo(video)
	//if err != nil {
	//	return nil, err
	//}
	//_, err = video.Seek(0, io.SeekStart)
	//if err != nil {
	//	return nil, err
	//}

	return &VideoSamples{
		//CompressedVideo: compressed,
		PreviewBlurred: blurred,
		Preview:        preview,
	}, nil
}

func createTempVideoFile(in io.Reader) (*os.File, string, error) {
	fname := fmt.Sprintf("video_" + uuid.New().String() + ".mp4")
	f, err := os.Create(fname)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(f, in)
	if err != nil {
		return nil, "", err
	}
	return f, fname, nil
}

type VideoSamples struct {
	CompressedVideo *bytes.Reader
	PreviewBlurred  *bytes.Reader
	Preview         *bytes.Reader
}
