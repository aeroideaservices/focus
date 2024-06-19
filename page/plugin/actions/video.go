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

const (
	yandexSpeechUrl          = "https://transcribe.api.cloud.yandex.net/speech/stt/v2/longRunningRecognize"
	yandexSpeechOperationUrl = "https://operation.api.cloud.yandex.net/operations/"
	chunks                   = 20
	audioEncoding            = "MP3"
	videoFormats             = ".mp4"
)

type VideoUseCase struct {
	medias       *mediaActions.Medias
	logger       *zap.SugaredLogger
	yandexApiKey string
}

func NewVideoUseCase(
	medias *mediaActions.Medias,
	logger *zap.SugaredLogger,
	yandexApiKey string,
) *VideoUseCase {
	return &VideoUseCase{
		medias:       medias,
		logger:       logger,
		yandexApiKey: yandexApiKey,
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

	go uc.GenerateSubtitles(ctx, []uuid.UUID{ids[0]})

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
		if !strings.Contains(videoFormats, filepath.Ext(fileName)) {
			os.Remove(fileName)
			return errors.New("file is not video")
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

		//  url of audio to yandex speech
		operation, err := uc.requestYandexSpeech(uri)

		if err != nil {
			return err
		}

		operation, err = uc.getResultOfYandexSpeech(operation.Id)
		if err != nil {
			return err
		}

		saveOperations, err := uc.splitSubtitles(*operation, chunks)
		if err != nil {
			return err
		}

		subJson, err := json.Marshal(saveOperations)
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
	return nil
}

func (uc VideoUseCase) UpdateSubtitles(ctx context.Context, subtitles SubtitlesToSave, mediaId uuid.UUID) error {

	updatedSubtitles := uc.updateChunks(subtitles, chunks)

	subJson, err := json.Marshal(updatedSubtitles)
	if err != nil {
		return err
	}

	updSubtitles := &mediaActions.UpdateMediaSubtitles{
		Id: mediaId,
	}

	err = updSubtitles.Subtitles.Scan(subJson)
	if err != nil {
		return err
	}

	return uc.medias.UpdateSubtitles(ctx, *updSubtitles)
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
				AudioEncoding:   audioEncoding,
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

	req, err := http.NewRequest("POST", yandexSpeechUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Api-Key "+uc.yandexApiKey)

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
	url := yandexSpeechOperationUrl + id
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Api-Key "+uc.yandexApiKey)
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

func (uc VideoUseCase) splitSubtitles(operation YandexSpeechOperationResult, chunks int) (*SubtitlesToSave, error) {
	if len(operation.Response.Chunks) < 0 {
		return nil, errors.New("no chunks in response")
	}
	if len(operation.Response.Chunks[0].Alternatives[0].Words) == 0 {
		return nil, errors.New("no words in response")
	}
	result := &SubtitlesToSave{
		FullText: operation.Response.Chunks[0].Alternatives[0].Text,
		Chunks:   make([]ChunkToSave, chunks),
	}

	chunkSize := len(operation.Response.Chunks[0].Alternatives[0].Words) / chunks
	extraWords := len(operation.Response.Chunks[0].Alternatives[0].Words) % chunks

	chunkIndex := 0

	for i := 0; i < len(operation.Response.Chunks[0].Alternatives[0].Words); i += chunkSize {
		end := i + chunkSize
		if end > len(operation.Response.Chunks[0].Alternatives[0].Words) {
			end = len(operation.Response.Chunks[0].Alternatives[0].Words)
		}
		extraWordFlag := chunkIndex < extraWords

		if extraWordFlag {
			end++
		}

		startTime := operation.Response.Chunks[0].Alternatives[0].Words[i].StartTime
		endTime := operation.Response.Chunks[0].Alternatives[0].Words[end-1].EndTime
		chunkText := ""
		for j := i; j < end; j++ {
			chunkText += operation.Response.Chunks[0].Alternatives[0].Words[j].Word + " "
		}

		result.Chunks[chunkIndex] = ChunkToSave{
			StartTime: startTime,
			EndTime:   endTime,
			Text:      chunkText,
		}

		if extraWordFlag {
			i++
		}

		chunkIndex++
	}
	return result, nil
}

func (uc VideoUseCase) updateChunks(subtitles SubtitlesToSave, chunks int) *SubtitlesToSave {
	newWords := strings.Split(subtitles.FullText, " ")
	chunkSize := len(newWords) / chunks
	extraWords := len(newWords) % chunks
	chunkIndex := 0

	for i := 0; i < len(newWords); i += chunkSize {
		end := i + chunkSize
		if end > len(newWords) {
			end = len(newWords)
		}
		if chunkIndex < extraWords {
			end++
		}

		chunkText := ""
		for j := i; j < end; j++ {
			chunkText += newWords[j] + " "
		}
		subtitles.Chunks[chunkIndex].Text = chunkText

		if chunkIndex < extraWords {
			i++
		}
		chunkIndex++
	}

	return &subtitles
}
