# Примеры реализации моделей

В этом модуле приведены примеры использования тегов `focus` и `validate`
для управления моделями.

### Подключение моделей

Подключить модели после создания реестра моделей `focus.ModelsRegistry`

```go
registry := focus.NewModelRegistry(true)
registry.Register(Category{}, Product{}, Promo{}, Store{})
```

или при использовании `di.Container` прописать объект в контейнере.

```go
{
    Name: "focus.models.registry.models",
    Build: func(ctn di.Container) (interface{}, error) {
        return []any{
            Category{},
            Product{},
            Promo{},
            Store{},
        }, nil
    },
}
```
