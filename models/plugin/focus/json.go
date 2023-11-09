package focus

import (
	jsoniter "github.com/json-iterator/go"
	"strings"
	"unicode"
)

var json = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	TagKey:                 "-",
}.Froze()

type namingStrategyExtension struct {
	jsoniter.DummyExtension
}

func init() {
	json.RegisterExtension(&namingStrategyExtension{})
}

func (extension *namingStrategyExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		if unicode.IsLower(rune(binding.Field.Name()[0])) || binding.Field.Name()[0] == '_' {
			continue
		}
		tag, hastag := binding.Field.Tag().Lookup("focus")
		if hastag {
			tagParts := parseTagParts(tag)
			if code, ok := tagParts["code"]; ok && code != "" {
				binding.ToNames = []string{code}
				binding.FromNames = []string{code}
				continue
			}
		}
		binding.ToNames = []string{extension.translate(binding.Field.Name())}
		binding.FromNames = []string{extension.translate(binding.Field.Name())}
	}
	extension.DummyExtension.UpdateStructDescriptor(structDescriptor)
}

func parseTagParts(tag string) map[string]string {
	tagPartsValues := make(map[string]string)
	tagParts := strings.Split(tag, ";")
	for _, part := range tagParts {
		keyValue := strings.SplitN(part, ":", 2)
		key := keyValue[0]
		value := ""
		if len(keyValue) == 2 {
			value = keyValue[1]
		}
		tagPartsValues[key] = value
	}

	return tagPartsValues
}

func (extension *namingStrategyExtension) translate(name string) string {
	return firstToLower(name)
}
