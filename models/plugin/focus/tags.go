package focus

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"github.com/aeroideaservices/focus/models/plugin/form"
	focusStrings "github.com/aeroideaservices/focus/services/formatting/strings"
)

type FieldFiller struct {
	Fill    func(field *Field, value string) // Fill проставление значения на основе значения
	Default func(field *Field)               // Default проставление значения по умолчанию
	Code    string
}

type FillFunc func(*Field, string)

var (
	commonTags = []FieldFiller{
		{Code: "title", Fill: titleFill, Default: titleDefault},
		{Code: "column", Fill: columnFill, Default: columnDefault},
		{Code: "code", Fill: codeFill, Default: codeDefault},
		{Code: "time", Fill: timeFill, Default: timeDefault},
		{Code: "media", Fill: mediaFill},
		{Code: "view", Fill: viewFill, Default: viewDefault},
		{Code: "viewExtra", Fill: viewExtraFill},
		{Code: "multiple", Fill: multipleFill, Default: multipleDefault},
		{Code: "filterable", Fill: filterFill},
		{Code: "sortable", Fill: sortFill},
		{Code: "primaryKey", Fill: primaryKeyFill},
		{Code: "disabled", Fill: disabledFill, Default: disabledDefault},
		{Code: "hidden", Fill: hiddenFill, Default: hiddenDefault},
		{Code: "position", Fill: positionFill},
		{Code: "block", Fill: blockFill},
		{Code: "unique", Fill: uniqueFill},
		{Code: "precision", Fill: precisionFill},
		{Code: "step", Fill: stepFill},
	}

	associationTags = []FieldFiller{
		{Code: "many2many", Fill: many2manyFill, Default: many2manyDefault},
		{Code: "association", Fill: associationFill, Default: associationDefault},
		{Code: "foreignKey", Fill: foreignKeyFill, Default: foreignKeyDefault},
		{Code: "references", Fill: referencesFill, Default: referencesDefault},
		{Code: "joinForeignKey", Fill: joinForeignKeyFill, Default: joinForeignKeyDefault},
		{Code: "joinReferences", Fill: joinReferencesFill, Default: joinReferencesDefault},
		{Code: "joinSort", Fill: joinSortFill},
	}
)

func titleFill(field *Field, value string) { field.Title = value }

func columnFill(field *Field, value string) { field.Column = value }

func codeFill(field *Field, value string) { field.Code = value }

func timeFill(field *Field, value string) {
	if value == "" {
		field.IsTime = true
		return
	}
	field.IsTime, _ = strconv.ParseBool(value)
}

func mediaFill(field *Field, value string) {
	if value == "" {
		field.IsMedia = true
		return
	}
	field.IsMedia, _ = strconv.ParseBool(value)
}

func viewFill(field *Field, value string) {
	if !slices.Contains(form.FieldTypes, form.FieldType(value)) {
		log.Panicf("'view' tag must be one of %s, got %s", form.FieldTypes, value)
	}
	field.View = form.FieldType(value)
}

func viewExtraFill(field *Field, value string) {
	ve := form.ResolveViewExtras(value)
	if ve == nil {
		panic("wrong viewExtra field definition")
	}
	field.ViewExtra = ve
}

func multipleFill(field *Field, value string) {
	if value == "" || value == "true" {
		field.Multiple = true
	}
}

func filterFill(field *Field, value string) {
	if field.IsMedia || field.Association != nil {
		field.Filterable = false
		return
	}
	if value == "" || value == "true" {
		field.Filterable = true
		return
	}
}

func sortFill(field *Field, value string) {
	if value == "" || value == "true" {
		field.Sortable = true
	}
}

func primaryKeyFill(field *Field, value string) {
	if value == "" || value == "true" {
		field.Model.PrimaryKey = field
		field.primaryKey = true
		return
	}
}

func disabledFill(field *Field, value string) {
	value = strings.ReplaceAll(value, " ", "")
	values := strings.Split(value, ",")
	for _, v := range values {
		switch v {
		case "create":
			field.Disabled = append(field.Disabled, CreateView)
		case "update":
			field.Disabled = append(field.Disabled, UpdateView)
		case "false":
			field.Disabled = nil
			return
		case "true", "":
			field.Disabled = []view{CreateView, UpdateView}
			return
		default:
			panic("wrong tag value for disabled field")
		}
	}
}

func hiddenFill(field *Field, value string) {
	value = strings.ReplaceAll(value, " ", "")
	values := strings.Split(value, ",")
	for _, v := range values {
		switch v {
		case "create":
			field.Hidden = append(field.Hidden, CreateView)
		case "update":
			field.Hidden = append(field.Hidden, UpdateView)
		case "list":
			field.Hidden = append(field.Hidden, ListView)
		case "false":
			field.Hidden = nil
			return
		case "true", "":
			field.Hidden = []view{CreateView, UpdateView, ListView}
			return
		default:
			panic("wrong tag value for hidden field")
		}
	}
}

func positionFill(field *Field, value string) {
	field.Position, _ = strconv.Atoi(value)
}

func blockFill(field *Field, value string) {
	field.Block = value
}

func uniqueFill(field *Field, value string) {
	if value == "" {
		field.IsUnique = true
		return
	}
	field.IsUnique, _ = strconv.ParseBool(value)
}

func titleDefault(field *Field) { field.Title = field.name }

func columnDefault(field *Field) { field.Column = focusStrings.CamelToSnakeCase(field.name) }

func codeDefault(field *Field) { field.Code = firstToLower(field.name) }

func timeDefault(field *Field) { field.IsTime = field.RawType() == reflect.TypeOf(time.Time{}) }

func multipleDefault(field *Field) {
	field.Multiple = field.t.Kind() == reflect.Slice
}

func disabledDefault(field *Field) {
	if field.primaryKey {
		field.Disabled = []view{CreateView, UpdateView}
	}
}

func hiddenDefault(field *Field) {
	if field.primaryKey {
		field.Hidden = []view{CreateView, UpdateView}
	}
}

func viewDefault(field *Field) {
	switch {
	case field.IsMedia:
		field.View = form.Media
		return
	case field.IsTime:
		field.View = form.DateTimePicker
		return
	case field.Association != nil:
		field.View = form.Select
		return
	case field.t == reflect.TypeOf(uuid.UUID{}) || field.t == reflect.TypeOf(&uuid.UUID{}):
		field.View = form.TextInput
		return
	}

	switch field.RawKind() {
	case reflect.Bool:
		field.View = form.Checkbox
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.View = form.IntInput
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.View = form.UintInput
	case reflect.Float32, reflect.Float64:
		field.View = form.FloatInput
	case reflect.String:
		field.View = form.TextInput
	case reflect.Struct:
		field.View = form.Select
	default:
		panic("you must specify field view")
	}
}

func precisionFill(field *Field, value string) {
	switch field.RawKind() {
	case reflect.Float32, reflect.Float64:
		if field.FloatProperties == nil {
			field.FloatProperties = &FloatProperties{}
		}

		v, err := strconv.Atoi(value)
		if err != nil {
			panic("the tag precision value should be int")
		}

		field.Precision = v
	default:
		panic("precision tag can only be applied to a float field")
	}
}

func stepFill(field *Field, value string) {
	switch field.RawKind() {
	case reflect.Float32, reflect.Float64:
		if field.FloatProperties == nil {
			field.FloatProperties = &FloatProperties{}
		}

		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic("the tag precision value should be float")
		}

		field.Step = v
	default:
		panic("precision tag can only be applied to a float field")
	}
}

func many2manyFill(field *Field, value string) {
	field.Association.Many2Many = value
	field.Association.Type = ManyToMany
}

func associationFill(field *Field, value string) {
	if field.Association.Type != None {
		return
	}
	switch value {
	case "belongsTo":
		field.Association.Type = BelongsTo
	case "hasOne":
		field.Association.Type = HasOne
	case "hasMany":
		field.Association.Type = HasMany
	case "many2many", "manyToMany":
		field.Association.Type = ManyToMany
	}
}

func foreignKeyFill(field *Field, value string) { field.Association.ForeignKey = value }

func referencesFill(field *Field, value string) { field.Association.References = value }

func joinForeignKeyFill(field *Field, value string) { field.Association.JoinForeignKey = value }

func joinReferencesFill(field *Field, value string) { field.Association.JoinReferences = value }

func joinSortFill(field *Field, value string) {
	field.Association.JoinSort = value
}

func many2manyDefault(field *Field) {
	assoc := field.Association
	if assoc.Many2Many != "" || assoc.Type != ManyToMany {
		return
	}
	assoc.Many2Many = field.Model.TableName + "_" + assoc.Model.TableName
}

func associationDefault(field *Field) {
	if field.Association.Type != None {
		return
	}
	switch field.t.Kind() {
	case reflect.Slice: // HasMany
		field.Association.Type = HasMany
	default: // HasOne or BelongsTo
		field.Association.Type = BelongsTo
	}
}

func foreignKeyDefault(field *Field) {
	if field.Association.ForeignKey != "" {
		return
	}
	model := field.Model
	assoc := field.Association
	switch assoc.Type {
	case BelongsTo:
		assoc.ForeignKey = field.Column + "_" + assoc.Model.PrimaryKey.Column
	case HasOne, HasMany:
		assoc.ForeignKey = strings.TrimSuffix(model.TableName, "s") + "_" + model.PrimaryKey.Column
	case ManyToMany:
		assoc.ForeignKey = model.PrimaryKey.Column
	}
}

func referencesDefault(field *Field) {
	if field.Association.References != "" {
		return
	}
	model := field.Model
	association := field.Association
	switch field.Association.Type {
	case BelongsTo, ManyToMany:
		association.References = association.Model.PrimaryKey.Column
	case HasOne, HasMany:
		association.References = model.PrimaryKey.Column
	}
}

func joinForeignKeyDefault(field *Field) {
	assoc := field.Association
	if assoc.JoinForeignKey != "" {
		return
	}
	if assoc.Type != ManyToMany {
		return
	}

	assoc.JoinForeignKey = strings.TrimSuffix(field.Model.TableName, "s") + "_" + assoc.ForeignKey
}

func joinReferencesDefault(field *Field) {
	assoc := field.Association
	if assoc.JoinReferences != "" {
		return
	}
	if assoc.Type != ManyToMany {
		return
	}

	assoc.JoinReferences = strings.TrimSuffix(assoc.Model.TableName, "s") + "_" + assoc.References
}

func firstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}
