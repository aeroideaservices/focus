package repositories

import (
	"context"
	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
	user_entity "github.com/jemzee04/focus/page/plugin/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetById(ctx context.Context, userId uuid.UUID) (*user_entity.User, error) {
	user := &user_entity.User{}
	db := r.db.WithContext(ctx).Model(user_entity.User{}).Where("id", userId)

	//if !allowInactive {
	//	db.Scopes()
	//}
	err := db.First(user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "user with id %s not found", userId)
	}
	if err != nil {
		return nil, err
	}

	db = r.db.WithContext(ctx).Model(entity.Media{}).Table("media").Where("id", user.PictureId)
	db.First(&user.Picture)

	return user, nil
}
