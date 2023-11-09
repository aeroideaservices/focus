package focus

type AssociationType int

const (
	None AssociationType = iota
	BelongsTo
	HasOne
	HasMany
	ManyToMany
)

// Ассоциация модели
type Association struct {
	Type           AssociationType // Type Тип ассоциации, может быть belongsTo, hasOne, hasMany, manyToMany (как в gorm)
	Model          *Model          // Model Модель, к которой происходит ассоциация
	Many2Many      string          // Many2Many Название промежуточной таблицы для отношения many2many
	ForeignKey     string          // ForeignKey Внешний ключ ассоциации
	References     string          // References Референс внешнего ключа
	JoinForeignKey string          // JoinForeignKey Внешний ключ в промежуточной таблице
	JoinReferences string          // JoinReferences Референс в промежуточной таблице
	JoinSort       string          // JoinSort сортировка ассоциации (только для типа many2many)
	modelCode      string          // modelCode код ассоциированной модели
}
