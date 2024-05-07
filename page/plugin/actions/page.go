package actions

import (
	"context"
	"github.com/aeroideaservices/focus/page/plugin/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PageUseCase struct {
	pageRepository    PageRepository
	galleryRepository GalleryRepository
	galleryUseCase    GalleryUseCase
	copierService     CopierInterface
	logger            *zap.SugaredLogger
}

func NewPageUseCase(
	pageRepository PageRepository, galleryRepository GalleryRepository, galleryUseCase GalleryUseCase,
	copierService CopierInterface, logger *zap.SugaredLogger,
) *PageUseCase {
	return &PageUseCase{
		pageRepository:    pageRepository,
		galleryRepository: galleryRepository,
		galleryUseCase:    galleryUseCase,
		copierService:     copierService,
		logger:            logger,
	}
}

func (uc PageUseCase) GetById(ctx context.Context, dto GetPageDto) (*PageDto, error) {
	uc.logger.Debug("Getting page")

	page, err := uc.pageRepository.GetById(ctx, uuid.Must(uuid.Parse(dto.ID)))
	if err != nil {
		return nil, err
	}

	pageDto, err := uc.GetPageDto(*page)
	if err != nil {
		return nil, err
	}

	uc.logger.Debug("Got page")
	return pageDto, nil
}

func (uc PageUseCase) Create(ctx context.Context, dto *CreatePageRequest) (*CreatePageResponse, error) {
	uc.logger.Debug("Creating page")
	pageEntity, err := uc.GetPageFromDto(*dto)
	if err != nil {
		return nil, err
	}
	//pageEntity.ID = uuid.New()
	pageId, err := uc.pageRepository.Create(ctx, pageEntity)
	if err != nil {
		return nil, err
	}
	result := &CreatePageResponse{
		ID: *pageId,
	}
	uc.logger.Debug("Created page")
	return result, nil
}

func (uc PageUseCase) Delete(ctx context.Context, pageId uuid.UUID) error {
	uc.logger.Debug("Deleting page")
	_, err := uc.pageRepository.GetByIdWithoutAssociate(ctx, pageId)
	if err != nil {
		return err
	}
	err = uc.pageRepository.Delete(ctx, pageId)
	if err != nil {
		return err
	}
	uc.logger.Debug("Deleted page")
	return nil
}

func (uc PageUseCase) PatchGalleryPosition(ctx context.Context, dto *PatchGalleryPosition) error {
	_, err := uc.pageRepository.GetByIdWithoutAssociate(ctx, dto.PageID)
	if err != nil {
		return err
	}

	_, err = uc.galleryRepository.GetByIdWithoutAssociate(ctx, dto.GalleryID)
	if err != nil {
		return err
	}

	err = uc.pageRepository.PatchGalleryPosition(ctx, dto)
	if err != nil {
		return err
	}
	uc.logger.Debug("Patched page")
	return nil
}

func (uc PageUseCase) GetList(ctx context.Context) (*GetPageListResponse, error) {
	uc.logger.Debug("Getting pages list")
	pages, err := uc.pageRepository.GetList(ctx, []string{"id", "name", "code", "sort", "is_published"}, true)
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, err
	}

	res := &GetPageListResponse{PageList: pages}

	uc.logger.Debug("Got pages list")
	return res, nil
}

func (uc PageUseCase) LinkGalleries(ctx context.Context, pageId uuid.UUID, galleryIds []uuid.UUID) (error, bool) {
	_, err := uc.pageRepository.GetByIdWithoutAssociate(ctx, pageId)
	if err != nil {
		return err, false
	}

	for _, id := range galleryIds {
		_, err = uc.galleryRepository.GetByIdWithoutAssociate(ctx, id)
		if err != nil {
			return err, false
		}
	}

	lastPosition, err := uc.pageRepository.GetLastPosition(ctx, pageId)
	if err != nil {
		uc.logger.Error(err.Error())
		return err, false
	}

	newPosition := lastPosition + 1
	pagesGalleries := uc.getPagesGalleriesToCreate(pageId, newPosition, galleryIds)

	err, isUniqViolErr := uc.pageRepository.CreatePagesGalleries(ctx, pagesGalleries)
	if err != nil {
		uc.logger.Error(err.Error())
		return err, isUniqViolErr
	}

	return nil, false
}

func (uc PageUseCase) UnlinkGalleries(ctx context.Context, pageId uuid.UUID, galleryIds []uuid.UUID) error {

	_, err := uc.pageRepository.GetByIdWithoutAssociate(ctx, pageId)
	if err != nil {
		return err
	}

	for _, id := range galleryIds {
		_, err = uc.galleryRepository.GetByIdWithoutAssociate(ctx, id)
		if err != nil {
			return err
		}
	}

	//pagesGalleries := uc.getPagesGalleriesToDelete(pageId, galleryIds)

	err = uc.pageRepository.DeletePagesGalleries(ctx, pageId, galleryIds)
	if err != nil {
		uc.logger.Error(err.Error())
		return err
	}

	return nil
}

func (uc PageUseCase) PatchProperties(ctx context.Context, dto *PatchPageRequest) error {
	uc.logger.Debug("Patching page")

	_, err := uc.pageRepository.GetByIdWithoutAssociate(ctx, dto.ID)
	if err != nil {
		return err
	}

	pageEntity, err := uc.GetPageFromPatchDto(*dto)
	if err != nil {
		return err
	}
	//pageEntity.ID = uuid.New()
	err = uc.pageRepository.PatchProperties(ctx, pageEntity)
	if err != nil {
		return err
	}

	if dto.IsPublished != nil {
		uc.logger.Debug("Changing publish")
		err = uc.pageRepository.UpdatePublish(ctx, dto.ID, dto.IsPublished)
		if err != nil {
			return err
		}
		uc.logger.Debug("Changed publish")
	}
	uc.logger.Debug("Patched page")
	return nil
}

func (uc PageUseCase) GetPageFromDto(pageDto CreatePageRequest) (page *entity.Page, err error) {
	page = &entity.Page{}
	if err != nil {
		return nil, err
	}
	newPageId := uuid.New()

	err = uc.copierService.Copy(page, &pageDto)
	if err != nil {
		return nil, err
	}

	galleries, err := uc.galleryUseCase.GetGalleryFromDtos(pageDto.Galleries, newPageId)

	page.PagesGalleries = galleries
	page.ID = newPageId

	return
}

func (uc PageUseCase) GetPageDto(page entity.Page) (pageDto *PageDto, err error) {
	if err != nil {
		return nil, err
	}
	pageDto = &PageDto{}

	err = uc.copierService.Copy(pageDto, &page)
	if err != nil {
		return nil, err
	}

	galleryDtos, err := uc.galleryUseCase.GetGalleryDtos(page.PagesGalleries)
	if err != nil {
		return
	}
	pageDto.Galleries = galleryDtos
	return
}

func (uc PageUseCase) GetPageFromPatchDto(pageDto PatchPageRequest) (page *entity.Page, err error) {
	page = &entity.Page{}

	err = uc.copierService.Copy(page, &pageDto)
	if err != nil {
		return nil, err
	}

	return
}

func (uc PageUseCase) getPagesGalleriesToCreate(
	pageId uuid.UUID, startPosition int, galleryIds []uuid.UUID,
) []entity.PagesGalleries {
	var pagesGalleries []entity.PagesGalleries
	for i, _ := range galleryIds {
		pagesGalleries = append(
			pagesGalleries, entity.PagesGalleries{
				PagesID:   &pageId,
				GalleryID: &galleryIds[i],
				Position:  startPosition + i,
			},
		)
	}
	return pagesGalleries
}

func (uc PageUseCase) getPagesGalleriesToDelete(
	pageId uuid.UUID, galleryIds []uuid.UUID,
) []entity.PagesGalleries {
	var pagesGalleries []entity.PagesGalleries
	for i, _ := range galleryIds {
		pagesGalleries = append(
			pagesGalleries, entity.PagesGalleries{
				PagesID:   &pageId,
				GalleryID: &galleryIds[i],
			},
		)
	}
	return pagesGalleries
}
