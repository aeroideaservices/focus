package form

// FieldType
type FieldType string

const (
	None            FieldType = "none"
	Checkbox        FieldType = "checkbox"
	IntInput        FieldType = "intInput"
	UintInput       FieldType = "uintInput"
	FloatInput      FieldType = "floatInput"
	Rating          FieldType = "rating"
	Select          FieldType = "select"
	DatePickerInput FieldType = "datePickerInput"
	DateTimePicker  FieldType = "dateTimePicker"
	TimePicker      FieldType = "timePickerInput"
	Textarea        FieldType = "textarea"
	TextInput       FieldType = "textInput"
	Wysiwyg         FieldType = "wysiwyg"
	EditorJs        FieldType = "editorJs"
	Media           FieldType = "media"
	PhoneInput      FieldType = "phoneInput"
	EmailInput      FieldType = "emailInput"
	ColorsInput     FieldType = "colorsInput"
)

var (
	// FieldTypes
	FieldTypes = []FieldType{
		None,
		Checkbox,
		IntInput,
		UintInput,
		FloatInput,
		Rating,
		Select,
		DatePickerInput,
		DateTimePicker,
		TimePicker,
		Textarea,
		TextInput,
		Wysiwyg,
		EditorJs,
		Media,
		PhoneInput,
		EmailInput,
		ColorsInput,
	}
)

// Field
type Field struct {
	Code     string         `json:"code"`
	Title    string         `json:"title"`
	Type     FieldType      `json:"type"`
	Multiple bool           `json:"multiple"`
	Sortable bool           `json:"sortable"`
	Extra    map[string]any `json:"extra"`
}

// Request объект описания запроса. Используется для описания запроса в полях типа Select, Media, EditorJs
type Request struct {
	URI       string `json:"uri"`                 // URI относительный или полный путь
	Meth      string `json:"meth"`                // Meth метод http-запроса
	Service   string `json:"service,omitempty"`   // Service код сервиса, к которому идет запрос. Только для относительного URI
	Body      any    `json:"body,omitempty"`      // Body тело запроса
	Paginated bool   `json:"paginated,omitempty"` // Paginated нужно ли пагинировать запрос, пагинация передается в query параметрах запроса
}

// ViewExtras параметры отображения
type ViewExtras map[string]any

// ViewsExtras параметры отображений
type ViewsExtras map[string]ViewExtras

// Resolve получение параметров отображения
func (ve ViewsExtras) Resolve(code string) ViewExtras {
	if e, ok := ve[code]; ok {
		return e
	}
	return nil
}

// viewsExtras параметры отображений
var viewsExtras ViewsExtras

// RegisterViewsExtras регистрация новых параметров отображений
func RegisterViewsExtras(vx ViewsExtras) {
	viewsExtras = vx
}

// ResolveViewExtras получение параметров отображения
func ResolveViewExtras(code string) ViewExtras {
	return viewsExtras.Resolve(code)
}

type SelectData []struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}
