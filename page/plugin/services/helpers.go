package services

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func GetIdsFromStrings(ids []string) ([]uuid.UUID, error) {
	var result []uuid.UUID
	for _, idStr := range ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, nil
}

//func GetUrlByFilePath(path string, mediaProvider actions.MediaProvider) string {
//	var url string
//
//	if env.CdnUrl != "" && !strings.Contains(
//		path, env.CdnUrl,
//	) {
//		url = strings.Replace(
//			mediaProvider.GetUrlByFilepath(path), "?", "", 1,
//		)
//		return url
//	}
//	if env.CdnUrl == "" && !strings.Contains(
//		path, env.AwsEndpoint+"/"+env.AwsBucket,
//	) {
//		url = strings.Replace(
//			mediaProvider.GetUrlByFilepath(path), "?", "", 1,
//		)
//		return url
//	}
//	return path
//}

type Copier struct{}

func (c Copier) Copy(toValue interface{}, fromValue interface{}) (err error) {
	return copier.Copy(toValue, fromValue)
}
