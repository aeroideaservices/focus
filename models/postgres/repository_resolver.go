package postgres

import (
	"fmt"
	"github.com/aeroideaservices/focus/models/plugin/actions"
	"gorm.io/gorm"
)

type repositoryResolver struct {
	repositories map[string]actions.Repository
}

// NewRepositoryResolver конструктор repositoryResolver
func NewRepositoryResolver(db *gorm.DB, modelsRegistry actions.ModelsRegistry) *repositoryResolver {
	models := modelsRegistry.ListModels()
	repositories := make(map[string]actions.Repository)
	for _, model := range models {
		repositories[model.Code] = newElementsRepository(db, model)
	}
	return &repositoryResolver{
		repositories: repositories,
	}
}

// Resolve получение репозитория по коду модели
func (r repositoryResolver) Resolve(modelCode string) actions.Repository {
	if repo, ok := r.repositories[modelCode]; ok {
		return repo
	}

	panic(fmt.Sprintf("model with code '%s' does not registered", modelCode))
}
