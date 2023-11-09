package utils

import (
	"strings"
)

func PrepareLikeQuery(query string) string {
	query = strings.TrimSpace(query)
	query = strings.ReplaceAll(query, `%`, `\%`)
	query = strings.ReplaceAll(query, `_`, `\_`)

	return "%" + query + "%"
}
