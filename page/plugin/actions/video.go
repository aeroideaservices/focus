package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	mediaActions "github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/page/plugin/services"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func (uc VideoUseCase) GenerateSubtitles(ctx context.Context, mediaIds []uuid.UUID) error {

	for _, id := range mediaIds {
		// get file from s3
		fileName, err := uc.medias.Download(ctx, mediaActions.GetMedia{Id: id})
		if err != nil {
			os.Remove(fileName)
			return err
		}
		if !strings.Contains(fileName, ".mp4") {
			os.Remove(fileName)
			return errors.New("only mp4 files are supported")
		}

		// get audio from video
		audio, audioFN, err := uc.getAudioFromVideo(fileName)

		os.Remove(fileName)

		// save audio to s3
		uri, err := uc.medias.Upload(
			ctx, mediaActions.CreateMedia{
				Filename: audioFN,
				Size:     audio.Size(),
				File:     audio,
			},
		)

		os.Remove(audioFN)

		//uri = "https://storage.yandexcloud.net/aerosite/audio_testcreatesub.mp3"
		uri = "https://storage.yandexcloud.net/speechkittest404/audio_testcreatesub.mp3"

		//  url of audio to yandex speech
		operation, err := uc.requestYandexSpeech(uri)

		if err != nil {
			return err
		}

		operation, err = uc.getResultOfYandexSpeech(operation.Id)
		if err != nil {
			return err
		}

		subJson, err := json.Marshal(operation)
		if err != nil {
			return err
		}

		updSubtitles := &mediaActions.UpdateMediaSubtitles{
			Id: id,
		}

		err = updSubtitles.Subtitles.Scan(subJson)
		if err != nil {
			return err
		}

		err = uc.medias.UpdateSubtitles(ctx, *updSubtitles)
		if err != nil {
			return err
		}
	}

	// TODO: save subtitles to db
	return nil
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

func (uc VideoUseCase) getAudioFromVideo(fname string) (*bytes.Reader, string, error) {
	audio, audioFN, err := services.GetAudioFromVideo(fname)
	return audio, audioFN, err
}

func (uc VideoUseCase) requestYandexSpeech(uri string) (*YandexSpeechOperationResult, error) {
	requestBody := RecognitionRequest{
		Config: RecognitionConfig{
			Specification: Specification{
				ProfanityFilter: false,
				LiteratureText:  true,
				AudioEncoding:   "MP3",
				RawResults:      false,
			},
		},
		Audio: RecognitionAudio{
			Uri: uri,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return nil, err
	}

	url := "https://transcribe.api.cloud.yandex.net/speech/stt/v2/longRunningRecognize"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ajecivs7c9aafdp32rt7", "AQVNzJf4diVKe3Ixb6ezRboT-tsf2HcWgYyVv167")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	operation := &YandexSpeechOperationResult{}
	err = json.Unmarshal(body, operation)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

func (uc VideoUseCase) getResultOfYandexSpeech(id string) (*YandexSpeechOperationResult, error) {
	url := "https://operation.api.cloud.yandex.net/operations/" + id
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Api-Key AQVN1uQlxgRnlJpv43S4jtdH5cgvtcjaVVq6zEFZ")
	operation := &YandexSpeechOperationResult{}

	for true {
		time.Sleep(time.Second * 30)
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, operation)
		if err != nil {
			return nil, err
		}
		if operation.Done {
			res.Body.Close()
			return operation, nil
		}
		res.Body.Close()
	}

	return operation, nil
}
