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

type Copier struct{}

func (c Copier) Copy(toValue interface{}, fromValue interface{}) (err error) {
	return copier.Copy(toValue, fromValue)
}
