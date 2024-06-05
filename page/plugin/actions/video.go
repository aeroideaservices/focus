package actions

import (
	"bytes"
	"context"
	"fmt"
	mediaActions "github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/page/plugin/services"
	"go.uber.org/zap"
	"io"
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
	ids, err := uc.medias.UploadList(
		ctx, mediaActions.CreateMediasList{
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
		},
	)
	if err != nil {
		return nil, err
	}

	if len(ids) != 4 {
		return nil, fmt.Errorf("error uploading medias")
	}

	return &CreateVideoResponse{
		VideoId: ids[0],
		//VideoLiteId:      ids[1],
		PreviewId:        ids[2],
		PreviewBlurredId: ids[3],
	}, nil
}

func (uc VideoUseCase) createVideoSamples(video io.ReadSeeker) (*VideoSamples, error) {
	uc.logger.Debug("Creating video preview")
	preview, err := services.GetNFrame(video, 1)
	if err != nil {
		return nil, err
	}
	_, err = video.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	uc.logger.Debug("Creating video blurred preview")
	blurred, err := services.GetNFrameBlurred(video, 1, 35)
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

type VideoSamples struct {
	CompressedVideo *bytes.Reader
	PreviewBlurred  *bytes.Reader
	Preview         *bytes.Reader
}
