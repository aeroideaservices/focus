package actions

import (
	"context"
	"github.com/aeroideaservices/focus/page/plugin/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GalleryUseCase struct {
	galleryRepository GalleryRepository
	cardRepository    CardRepository
	cardUseCase       CardUseCase
	copierService     CopierInterface
	logger            *zap.SugaredLogger
}

func NewGalleryUseCase(
	galleryRepository GalleryRepository, cardRepository CardRepository, cardUseCase CardUseCase,
	copierService CopierInterface, logger *zap.SugaredLogger,
) *GalleryUseCase {
	return &GalleryUseCase{
		galleryRepository: galleryRepository,
		cardRepository:    cardRepository,
		cardUseCase:       cardUseCase,
		copierService:     copierService,
		logger:            logger,
	}
}

func (uc GalleryUseCase) GetByCode(ctx context.Context, code string) (*GalleryDto, error) {
	uc.logger.Debug("Getting gallery by code")

	gallery, err := uc.galleryRepository.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	//TODO: get cards
	uc.logger.Debug("got gallery")

	return &GalleryDto{
		ID:   gallery.ID,
		Name: gallery.Name,
		Code: gallery.Code,
	}, nil
}

func (uc GalleryUseCase) Create(ctx context.Context, dto *CreateGalleryRequest) (*CreateGalleryResponse, error) {
	uc.logger.Debug("Creating gallery")
	galleryEntity, err := uc.GetGalleryFromDto(dto)
	if err != nil {
		return nil, err
	}
	//pageEntity.ID = uuid.New()
	pageId, err := uc.galleryRepository.Create(ctx, galleryEntity)
	if err != nil {
		return nil, err
	}
	result := &CreateGalleryResponse{
		ID: *pageId,
	}
	uc.logger.Debug("Created gallery")
	return result, nil
}

func (uc GalleryUseCase) GetById(ctx context.Context, dto GetGalleryRequest) (*GalleryDto, error) {
	uc.logger.Debug("Getting gallery by id")

	gallery, err := uc.galleryRepository.GetById(ctx, uuid.Must(uuid.Parse(dto.ID)))
	if err != nil {
		return nil, err
	}

	galleryDto, err := uc.GetDtoFromGallery(gallery)
	if err != nil {
		return nil, err
	}

	uc.logger.Debug("Got gallery by id")
	return galleryDto, nil
}

func (uc GalleryUseCase) GetListWithSearch(ctx context.Context, searchValue string) ([]GalleryDto, error) {
	uc.logger.Debug("Getting list gallery with search")
	galleries, err := uc.galleryRepository.GetListWithSearch(ctx, searchValue)
	if err != nil {
		return nil, err
	}

	galleriesDto, err := uc.GetDtosFromGallery(galleries)
	if err != nil {
		return nil, err
	}
	uc.logger.Debug("Got list gallery with search")
	return galleriesDto, nil
}

func (uc GalleryUseCase) Update(ctx context.Context, dto *UpdateGalleryRequest) error {
	uc.logger.Debug("Updating gallery")

	_, err := uc.galleryRepository.GetByIdWithoutAssociate(ctx, dto.ID)
	if err != nil {
		return err
	}

	gallery, err := uc.GetGalleryFromUpdateDto(dto)
	if err != nil {
		return err
	}

	err = uc.galleryRepository.Update(ctx, gallery)
	if err != nil {
		return err
	}

	if dto.Hidden != nil {
		uc.logger.Debug("Changing hidden")
		err = uc.galleryRepository.UpdateHidden(ctx, dto.ID, dto.Hidden)
		if err != nil {
			return err
		}
		uc.logger.Debug("Changed hidden")
	}
	if dto.IsPublished != nil {
		uc.logger.Debug("Changing publish")
		err = uc.galleryRepository.UpdatePublish(ctx, dto.ID, dto.IsPublished)
		if err != nil {
			return err
		}
		uc.logger.Debug("Changed publish")
	}

	if err != nil {
		return err
	}
	uc.logger.Debug("Updated gallery")
	return nil
}

func (uc GalleryUseCase) PatchName(ctx context.Context, dto *PatchGalleryNameRequest) error {
	uc.logger.Debug("Patching gallery name")

	gallery := &entity.Gallery{ID: dto.ID, Name: dto.Name}

	err := uc.galleryRepository.PatchName(ctx, gallery)

	if err != nil {
		return err
	}
	uc.logger.Debug("Patched gallery name")
	return nil
}

func (uc GalleryUseCase) PatchCardPosition(ctx context.Context, dto *PatchCardPosition) error {
	uc.logger.Debug("Patching card position")

	_, err := uc.galleryRepository.GetByIdWithoutAssociate(ctx, dto.GalleryID)
	if err != nil {
		return err
	}
	_, err = uc.cardRepository.GetByIdWithoutAssociate(ctx, dto.CardID)
	if err != nil {
		return err
	}

	err = uc.galleryRepository.PatchCardPosition(ctx, dto)
	if err != nil {
		return err
	}
	uc.logger.Debug("Patched card position")
	return nil
}

func (uc GalleryUseCase) LinkCards(ctx context.Context, galleryId uuid.UUID, cardIds []uuid.UUID) (error, bool) {
	uc.logger.Debug("Linking cards to gallery")

	_, err := uc.galleryRepository.GetByIdWithoutAssociate(ctx, galleryId)
	if err != nil {
		return err, false
	}
	for _, id := range cardIds {
		_, err = uc.cardRepository.GetByIdWithoutAssociate(ctx, id)
		if err != nil {
			return err, false
		}
	}

	lastPosition, err := uc.galleryRepository.GetLastPosition(ctx, galleryId)
	if err != nil {
		uc.logger.Error(err.Error())
		return err, false
	}

	newPosition := lastPosition + 1
	galleriesCards := uc.getGalleriesCardsToCreate(galleryId, newPosition, cardIds)

	err, isUniqViolErr := uc.galleryRepository.CreateGalleriesCards(ctx, galleriesCards)
	if err != nil {
		uc.logger.Error(err.Error())
		return err, isUniqViolErr
	}

	uc.logger.Debug("Linked cards to gallery")
	return nil, false
}

func (uc GalleryUseCase) UnlinkCards(ctx context.Context, galleryId uuid.UUID, cardIds []uuid.UUID) error {
	uc.logger.Debug("Unlinking cards to gallery")

	_, err := uc.galleryRepository.GetByIdWithoutAssociate(ctx, galleryId)
	if err != nil {
		return err
	}

	for _, id := range cardIds {
		_, err = uc.cardRepository.GetByIdWithoutAssociate(ctx, id)
		if err != nil {
			return err
		}
	}

	err = uc.galleryRepository.DeleteGalleriesCards(ctx, galleryId, cardIds)
	if err != nil {
		uc.logger.Error(err.Error())
		return err
	}

	uc.logger.Debug("Unlinked cards to gallery")
	return nil
}

func (uc GalleryUseCase) DeleteList(ctx context.Context, galleryIds []uuid.UUID) error {
	uc.logger.Debug("Deleting galleries")

	//galleryIds, err := services.GetIdsFromStrings(galleryStrIds)
	//if err != nil {
	//	return err
	//}

	for _, id := range galleryIds {
		_, err := uc.galleryRepository.GetByIdWithoutAssociate(ctx, id)
		if err != nil {
			return err
		}
	}

	err := uc.galleryRepository.DeleteList(ctx, galleryIds)
	if err != nil {
		return err
	}
	uc.logger.Debug("Deleted galleries")
	return nil
}

func (uc GalleryUseCase) GetGalleryDtos(galleries []entity.PagesGalleries) (
	galleryDtos []GalleryInPageDto, err error,
) {
	for i, gallery := range galleries {
		galleryDto := &GalleryInPageDto{}
		err = uc.copierService.Copy(galleryDto, gallery.Gallery)
		if err != nil {
			return nil, err
		}
		galleryDto.Position = gallery.Position
		galleryDtos = append(galleryDtos, *galleryDto)

		galleryDtos[i].Cards, err = uc.cardUseCase.GetDtosFromGalleriesCard(gallery.Gallery.GalleriesCards)
		if err != nil {
			return nil, err
		}
	}
	return galleryDtos, err
}

func (uc GalleryUseCase) GetDtosFromGallery(galleries []entity.Gallery) (
	galleryDtos []GalleryDto, err error,
) {
	for i, gallery := range galleries {
		galleryDto := &GalleryDto{}
		err = uc.copierService.Copy(galleryDto, gallery)
		if err != nil {
			return nil, err
		}
		galleryDtos = append(galleryDtos, *galleryDto)

		galleryDtos[i].Cards, err = uc.cardUseCase.GetDtosFromGalleriesCard(gallery.GalleriesCards)
	}
	return galleryDtos, err
}

func (uc GalleryUseCase) GetDtoFromGallery(gallery *entity.Gallery) (
	galleryDto *GalleryDto, err error,
) {
	galleryDto = &GalleryDto{}

	err = uc.copierService.Copy(galleryDto, gallery)
	if err != nil {
		return nil, err
	}

	galleryDto.Cards, err = uc.cardUseCase.GetDtosFromGalleriesCard(gallery.GalleriesCards)

	return galleryDto, err
}

func (uc GalleryUseCase) GetGalleryFromDtos(galleryDtos []CreateGalleryInPageRequest, pageId uuid.UUID) (
	galleries []entity.PagesGalleries, err error,
) {
	for i, gallery := range galleryDtos {

		newGalleryId := uuid.New()

		pagesGalleries := &entity.PagesGalleries{}
		pagesGalleries.PagesID = &pageId
		pagesGalleries.GalleryID = &newGalleryId

		pagesGalleries.Position = gallery.Position

		galleryEntity := &entity.Gallery{}
		err = uc.copierService.Copy(galleryEntity, gallery)
		if err != nil {
			return nil, err
		}
		pagesGalleries.Gallery = *galleryEntity
		pagesGalleries.Gallery.ID = newGalleryId
		galleries = append(galleries, *pagesGalleries)

		galleries[i].Gallery.GalleriesCards, err = uc.cardUseCase.GetCardFromDtos(gallery.Cards, newGalleryId)
	}

	return galleries, err
}

func (uc GalleryUseCase) GetGalleryFromDto(galleryDto *CreateGalleryRequest) (
	gallery *entity.Gallery, err error,
) {
	gallery = &entity.Gallery{}

	newGalleryId := uuid.New()

	err = uc.copierService.Copy(gallery, galleryDto)
	if err != nil {
		return nil, err
	}
	gallery.ID = newGalleryId

	gallery.GalleriesCards, err = uc.cardUseCase.GetCardFromDtos(galleryDto.Cards, newGalleryId)

	return gallery, err
}

func (uc GalleryUseCase) GetGalleryFromUpdateDto(galleryDto *UpdateGalleryRequest) (
	gallery *entity.Gallery, err error,
) {
	gallery = &entity.Gallery{}
	err = uc.copierService.Copy(gallery, galleryDto)
	if err != nil {
		return nil, err
	}

	gallery.GalleriesCards, err = uc.cardUseCase.GetCardFromDtos(galleryDto.Cards, galleryDto.ID)
	if err != nil {
		return nil, err
	}

	return gallery, err
}

func (uc GalleryUseCase) getGalleriesCardsToCreate(
	galleryId uuid.UUID, startPosition int, cardIds []uuid.UUID,
) []entity.GalleriesCards {
	var galleriesCards []entity.GalleriesCards
	for i, _ := range cardIds {
		galleriesCards = append(
			galleriesCards, entity.GalleriesCards{
				GalleryID: &galleryId,
				CardID:    &cardIds[i],
				Position:  startPosition + i,
			},
		)
	}
	return galleriesCards
}
