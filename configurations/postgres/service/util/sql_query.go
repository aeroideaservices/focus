package util

import (
	"fmt"
	"strings"
)

func PrepareLikeValue(value *string) {
	*value = strings.TrimSpace(*value)
	*value = strings.ReplaceAll(*value, "_", "\\_")
	*value = strings.ReplaceAll(*value, "%", "\\%")
	*value = fmt.Sprintf("%%%s%%", strings.ToLower(*value))
}
