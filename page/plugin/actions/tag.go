package actions

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"pages/pkg/page/plugin/entity"
)

type TagUseCase struct {
	tagRepository TagRepository
	copierService CopierInterface
	logger        *zap.SugaredLogger
}

func NewTagUseCase(tagRepository TagRepository, copierService CopierInterface, logger *zap.SugaredLogger) *TagUseCase {
	return &TagUseCase{
		tagRepository: tagRepository,
		copierService: copierService,
		logger:        logger,
	}
}

func (uc TagUseCase) GetListWithSearch(ctx context.Context, searchValue string) ([]TagSearchDto, error) {
	uc.logger.Debug("Getting list tags with search")
	tags, tagCardsIds, err := uc.tagRepository.GetListWithSearch(ctx, searchValue)
	if err != nil {
		return nil, err
	}

	tagDtos, err := uc.GetDtosFromTags(tags, tagCardsIds)
	if err != nil {
		return nil, err
	}
	uc.logger.Debug("Got list tags with search")
	return tagDtos, nil
}

func (uc TagUseCase) Create(ctx context.Context, dto *TagDto) (
	*CreateTagResponse, error,
) {
	uc.logger.Debug("Creating tag")
	tagEntity, err := uc.GetTagFromDto(*dto)

	tagEntity.ID = uuid.New()

	tagId, err := uc.tagRepository.Create(ctx, tagEntity)
	if err != nil {
		return nil, err
	}
	result := &CreateTagResponse{
		ID: *tagId,
	}
	uc.logger.Debug("Created tag")
	return result, nil
}

func (uc TagUseCase) Update(ctx context.Context, dto *TagDtoRequest) error {
	uc.logger.Debug("Updating tag")

	_, err := uc.tagRepository.GetById(ctx, dto.ID)
	if err != nil {
		return err
	}

	tagEntity, err := uc.GetTagFromRequestDto(*dto)
	if err != nil {
		return err
	}

	err = uc.tagRepository.Update(ctx, tagEntity)

	if err != nil {
		return err
	}

	uc.logger.Debug("Updated tag")
	return nil
}

func (uc TagUseCase) Delete(ctx context.Context, tagId uuid.UUID) error {
	uc.logger.Debug("Deleting card")
	_, err := uc.tagRepository.GetById(ctx, tagId)
	if err != nil {
		return err
	}
	err = uc.tagRepository.Delete(ctx, tagId)
	if err != nil {
		return err
	}
	uc.logger.Debug("Deleted card")
	return nil
}

func (uc TagUseCase) GetDtosFromTags(tags []entity.Tag, tagCardsIds map[uuid.UUID][]uuid.UUID) ([]TagSearchDto, error) {
	var dtos []TagSearchDto

	for _, tag := range tags {
		tagDto := &TagSearchDto{}
		err := uc.copierService.Copy(tagDto, &tag)
		if err != nil {
			return nil, err
		}
		tagDto.CardIds = tagCardsIds[tag.ID]
		dtos = append(dtos, *tagDto)
	}
	return dtos, nil
}

func (uc TagUseCase) GetTagFromDto(tag TagDto) (*entity.Tag, error) {
	tagEntity := &entity.Tag{}
	err := uc.copierService.Copy(tagEntity, &tag)
	if err != nil {
		return nil, err
	}
	return tagEntity, nil
}

func (uc TagUseCase) GetTagFromRequestDto(tag TagDtoRequest) (*entity.Tag, error) {
	tagEntity := &entity.Tag{}
	err := uc.copierService.Copy(tagEntity, &tag)
	if err != nil {
		return nil, err
	}
	return tagEntity, nil
}
