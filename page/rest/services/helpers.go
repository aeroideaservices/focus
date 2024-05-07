package services

import "github.com/google/uuid"

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
