package actions

import (
	"context"
	"github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/page/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
	actions2 "gitlab.aeroidea.ru/internal-projects/focus/forms/plugin/actions"
	"go.uber.org/zap"
)

type CardUseCase struct {
	cardRepository    CardRepository
	galleryRepository GalleryRepository
	tagRepository     TagRepository
	mediaProvider     actions.MediaProvider
	formActions       actions2.Forms
	copierService     CopierInterface
	logger            *zap.SugaredLogger
}

func NewCardUseCase(
	cardRepository CardRepository, galleryRepository GalleryRepository, tagRepository TagRepository,
	mediaProvider actions.MediaProvider,
	formActions actions2.Forms, copierService CopierInterface, logger *zap.SugaredLogger,
) *CardUseCase {
	return &CardUseCase{
		cardRepository:    cardRepository,
		galleryRepository: galleryRepository,
		tagRepository:     tagRepository,
		mediaProvider:     mediaProvider,
		formActions:       formActions,
		copierService:     copierService,
		logger:            logger,
	}
}

func (uc CardUseCase) GetListWithSearch(ctx context.Context, searchValue string, name string) (
	[]CardDtoWithoutPosition, error,
) {
	uc.logger.Debug("Getting list cards with search")
	cards, err := uc.cardRepository.GetListWithSearch(ctx, searchValue, name)
	if err != nil {
		return nil, err
	}

	cardsDto, err := uc.GetDtosFromCard(cards)
	if err != nil {
		return nil, err
	}
	uc.logger.Debug("Got list cards with search")
	return cardsDto, nil
}

func (uc CardUseCase) Create(ctx context.Context, dto *CreateCardRequest) (
	*CreateCardResponse, error,
) {
	cardEntity, err := uc.GetCardFromDto(*dto)
	if err != nil {
		return nil, err
	}
	galleriesCard := uc.GetGalleriesCardsFromIds(ctx, cardEntity.ID, dto.GalleryIds)

	cardId, err := uc.cardRepository.Create(ctx, cardEntity, galleriesCard)
	if err != nil {
		return nil, err
	}
	result := &CreateCardResponse{
		ID: *cardId,
	}
	uc.logger.Debug("Created card")
	return result, nil
}

func (uc CardUseCase) GetById(ctx context.Context, dto GetCardRequest) (*CardDtoWithoutPosition, error) {
	uc.logger.Debug("Getting card")

	card, err := uc.cardRepository.GetById(ctx, uuid.Must(uuid.Parse(dto.ID)))
	if err != nil {
		return nil, err
	}

	cardDto, err := uc.GetDtoFromCard(card)

	uc.logger.Debug("Got card")
	return cardDto, nil
}

func (uc CardUseCase) Update(ctx context.Context, dto *UpdateCardRequest) error {
	uc.logger.Debug("Updating card")

	_, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, dto.ID)
	if err != nil {
		return err
	}

	for _, galleryId := range dto.GalleryIds {
		_, err = uc.galleryRepository.GetByIdWithoutAssociate(ctx, galleryId)
		if err != nil {
			return err
		}
	}

	card, err := uc.GetCardFromUpdateDto(dto)
	if err != nil {
		return err
	}

	galleriesCard := uc.GetGalleriesCardsFromIds(ctx, dto.ID, dto.GalleryIds)

	err = uc.cardRepository.Update(ctx, card, galleriesCard)

	if err != nil {
		return err
	}

	if dto.IsPublished != nil {
		uc.logger.Debug("Changing publish")
		err = uc.cardRepository.UpdatePublish(ctx, dto.ID, dto.IsPublished)
		if err != nil {
			return err
		}
		uc.logger.Debug("Changed publish")
	}

	if dto.Type == "regular" && dto.RegularCard.Inverted != nil {
		uc.logger.Debug("Changing inverted")
		err = uc.cardRepository.UpdateInverted(ctx, dto.ID, dto.RegularCard.Inverted)
		if err != nil {
			return err
		}
		uc.logger.Debug("Changed inverted")
	}

	if err != nil {
		return err
	}
	uc.logger.Debug("Updated card")
	return nil
}

func (uc CardUseCase) Delete(ctx context.Context, cardId uuid.UUID) error {
	uc.logger.Debug("Deleting card")

	_, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, cardId)
	if err != nil {
		return err
	}

	err = uc.cardRepository.Delete(ctx, cardId)
	if err != nil {
		return err
	}

	uc.logger.Debug("Deleted card")
	return nil
}

func (uc CardUseCase) PatchUser(ctx context.Context, dto *PatchUserRequest) error {
	_, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, dto.CardId)
	if err != nil {
		return err
	}

	card := &entity.Card{
		ID: dto.CardId,
		RegularCard: &entity.RegularCard{
			UserId: &dto.UserId,
		},
	}
	err = uc.cardRepository.PatchUser(ctx, card)
	if err != nil {
		return err
	}
	uc.logger.Debug("Patched user card")
	return nil
}

func (uc CardUseCase) PatchPreviewText(ctx context.Context, dto *PatchPreviewTextRequest) error {
	_, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, dto.CardId)
	if err != nil {
		return err
	}

	card := &entity.Card{
		ID: dto.CardId,
		RegularCard: &entity.RegularCard{
			PreviewText: dto.PreviewText,
		},
	}
	err = uc.cardRepository.PatchPreviewText(ctx, card)
	if err != nil {
		return err
	}
	uc.logger.Debug("Patched preview text card")
	return nil
}

func (uc CardUseCase) PatchDetailText(ctx context.Context, dto *PatchDetailTextRequest) error {
	_, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, dto.CardId)
	if err != nil {
		return err
	}
	uc.logger.Debug("Patching detail text card")
	card := &entity.Card{
		ID: dto.CardId,
		RegularCard: &entity.RegularCard{
			DetailText: dto.DetailText,
		},
	}
	err = uc.cardRepository.PatchDetailText(ctx, card)
	if err != nil {
		return err
	}
	uc.logger.Debug("Patched detail text card")
	return nil
}

func (uc CardUseCase) PatchLearnMoreUrl(ctx context.Context, dto *PatchLearnMoreUrlRequest) error {
	cardEntity, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, dto.CardId)
	if err != nil {
		return err
	}
	uc.logger.Debug("Patching learn more url card")
	card := &entity.Card{
		ID: dto.CardId,
		RegularCard: &entity.RegularCard{
			LearnMoreUrl: dto.LearnMoreUrl,
		},
	}

	switch cardEntity.Type {
	case "regular":
		card.RegularCard = &entity.RegularCard{
			LearnMoreUrl: dto.LearnMoreUrl,
		}
	case "form":
		card.FormCard = &entity.FormCard{
			LearnMoreUrl: dto.LearnMoreUrl,
		}
	}

	err = uc.cardRepository.PatchLearnMoreUrl(ctx, card)
	if err != nil {
		return err
	}
	uc.logger.Debug("Patched learn more url card")
	return nil
}

func (uc CardUseCase) PatchTags(ctx context.Context, cardId uuid.UUID, dtos []TagDtoRequest) error {
	uc.logger.Debug("Patching card tags")

	_, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, cardId)
	if err != nil {
		return err
	}

	tags, err := uc.getTagsFromDtos(dtos)
	if err != nil {
		return err
	}
	err = uc.cardRepository.PatchTags(ctx, cardId, tags)
	if err != nil {
		return err
	}

	uc.logger.Debug("Patched card tags")
	return nil
}

func (uc CardUseCase) LinkTags(ctx context.Context, cardId uuid.UUID, tagIds []uuid.UUID) (error, bool) {
	uc.logger.Debug("Linking tags to card")

	_, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, cardId)
	if err != nil {
		return err, false
	}

	for _, tagId := range tagIds {
		_, err = uc.tagRepository.GetById(ctx, tagId)
		if err != nil {
			return err, false
		}
	}

	err, isUniqViolErr := uc.cardRepository.CreateRegularTags(ctx, cardId, tagIds)
	if err != nil {
		uc.logger.Error(err.Error())
		return err, isUniqViolErr
	}

	uc.logger.Debug("Linked tags to card")
	return nil, false
}

func (uc CardUseCase) UnlinkTags(ctx context.Context, cardId uuid.UUID, tagIds []uuid.UUID) error {
	uc.logger.Debug("Unlinking tags from card")

	_, err := uc.cardRepository.GetByIdWithoutAssociate(ctx, cardId)
	if err != nil {
		return err
	}

	for _, tagId := range tagIds {
		_, err = uc.tagRepository.GetById(ctx, tagId)
		if err != nil {
			return err
		}
	}

	err = uc.cardRepository.DeleteRegularTags(ctx, cardId, tagIds)
	if err != nil {
		uc.logger.Error(err.Error())
		return err
	}

	uc.logger.Debug("Unlinked tags from card")
	return nil
}

func (uc CardUseCase) GetDtosFromCard(cards []entity.Card) (cardDtos []CardDtoWithoutPosition, err error) {
	for _, card := range cards {
		cardDto := &CardDtoWithoutPosition{}

		switch card.Type {
		case "html":
			if card.HtmlCardId == nil {
				err = uc.copierService.Copy(cardDto, card)
				if err != nil {
					return nil, err
				}

				continue
			}
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}
			htmlCardDto := &HtmlCardDto{}
			err = uc.copierService.Copy(htmlCardDto, card.HtmlCard)
			if err != nil {
				return nil, err
			}
			cardDto.HtmlCard = htmlCardDto

		case "video":
			if card.VideoCardId == nil {

				err = uc.copierService.Copy(cardDto, card)
				if err != nil {
					return nil, err
				}

				continue
			}
			if card.VideoCard.VideoId != nil {
				card.VideoCard.Video.Filepath, err = uc.mediaProvider.GetUrlById(*card.VideoCard.VideoId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.VideoCard.VideoPreviewId != nil {
				card.VideoCard.VideoPreview.Filepath, err = uc.mediaProvider.GetUrlById(*card.VideoCard.VideoPreviewId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.VideoCard.VideoPreviewBlurId != nil {
				card.VideoCard.VideoPreviewBlur.Filepath, err = uc.mediaProvider.GetUrlById(*card.VideoCard.VideoPreviewBlurId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.VideoCard.VideoLiteId != nil {
				card.VideoCard.VideoLite.Filepath, err = uc.mediaProvider.GetUrlById(*card.VideoCard.VideoLiteId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}
			videoCardDto := &VideoCardDto{}
			err = uc.copierService.Copy(videoCardDto, card.VideoCard)
			if err != nil {
				return nil, err
			}
			cardDto.VideoCard = videoCardDto
			cardDto.VideoCard.Video = card.VideoCard.Video
			cardDto.VideoCard.VideoLite = card.VideoCard.VideoLite
			cardDto.VideoCard.VideoPreview = card.VideoCard.VideoPreview
			cardDto.VideoCard.VideoPreviewBlur = card.VideoCard.VideoPreviewBlur
		case "photo":
			if card.PhotoCardId == nil {
				err = uc.copierService.Copy(cardDto, card)
				if err != nil {
					return nil, err
				}
				continue
			}
			if card.PhotoCard.PictureId != nil {
				card.PhotoCard.Picture.Filepath, err = uc.mediaProvider.GetUrlById(*card.PhotoCard.PictureId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}
			photoCardDto := &PhotoCardDto{}
			err = uc.copierService.Copy(photoCardDto, card.PhotoCard)
			if err != nil {
				return nil, err
			}
			cardDto.PhotoCard = photoCardDto
			cardDto.PhotoCard.Picture = card.PhotoCard.Picture
		case "form":
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}
			if cardDto.FormCard.Background != nil {
				cardDto.FormCard.Background.Filepath, err = uc.mediaProvider.GetUrlById(cardDto.FormCard.Background.Id)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting background media url")
					return
				}
			}
			if cardDto.FormCard.User != nil && cardDto.FormCard.User.Picture != nil {
				cardDto.FormCard.User.Picture.Filepath, err = uc.mediaProvider.GetUrlById(cardDto.FormCard.User.Picture.Id)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			var tags []TagDto

			for _, tag := range card.FormCard.FormCardsTags {
				tagDto := &TagDto{}
				err = uc.copierService.Copy(tagDto, &tag.Tag)
				if err != nil {
					return nil, err
				}
				tags = append(tags, *tagDto)

			}
			cardDto.FormCard.Tags = tags

		case "regular":
			if card.RegularCardId == nil {
				err = uc.copierService.Copy(cardDto, card)
				if err != nil {
					return nil, err
				}
				continue
			}
			if card.RegularCard.VideoId != nil {
				card.RegularCard.Video.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.VideoId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.RegularCard.VideoPreviewId != nil {
				card.RegularCard.VideoPreview.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.VideoPreviewId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.RegularCard.VideoPreviewBlurId != nil {
				card.RegularCard.VideoPreviewBlur.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.VideoPreviewBlurId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.RegularCard.VideoLiteId != nil {
				card.RegularCard.VideoLite.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.VideoLiteId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}

			if card.RegularCard.UserId != nil && card.RegularCard.User.Picture != nil {
				card.RegularCard.User.Picture.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.User.PictureId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}

			var tags []TagDto

			for _, tag := range card.RegularCard.RegularCardsTags {
				tagDto := &TagDto{}
				err = uc.copierService.Copy(tagDto, &tag.Tag)
				if err != nil {
					return nil, err
				}
				tags = append(tags, *tagDto)

			}
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}
			regularCardDto := &RegularCardDto{}
			err = uc.copierService.Copy(regularCardDto, card.RegularCard)
			if err != nil {
				return nil, err
			}
			cardDto.RegularCard = regularCardDto
			cardDto.RegularCard.Tags = tags

		}
		cardDtos = append(cardDtos, *cardDto)
	}
	return cardDtos, nil
}

func (uc CardUseCase) GetDtoFromCard(card *entity.Card) (cardDto *CardDtoWithoutPosition, err error) {
	cardDto = &CardDtoWithoutPosition{}
	switch card.Type {
	case "html":
		if card.HtmlCardId == nil {
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}

		} else {
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}
			htmlCardDto := &HtmlCardDto{}
			err = uc.copierService.Copy(htmlCardDto, card.HtmlCard)
			if err != nil {
				return nil, err
			}
			cardDto.HtmlCard = htmlCardDto
		}
	case "video":
		if card.VideoCardId == nil {
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}

		}
		if card.VideoCard.VideoId != nil {
			card.VideoCard.Video.Filepath, err = uc.mediaProvider.GetUrlById(*card.VideoCard.VideoId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		if card.VideoCard.VideoPreviewId != nil {
			card.VideoCard.VideoPreview.Filepath, err = uc.mediaProvider.GetUrlById(*card.VideoCard.VideoPreviewId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		if card.VideoCard.VideoPreviewBlurId != nil {
			card.VideoCard.VideoPreviewBlur.Filepath, err = uc.mediaProvider.GetUrlById(*card.VideoCard.VideoPreviewBlurId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		if card.VideoCard.VideoLiteId != nil {
			card.VideoCard.VideoLite.Filepath, err = uc.mediaProvider.GetUrlById(*card.VideoCard.VideoLiteId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		err = uc.copierService.Copy(cardDto, card)
		if err != nil {
			return nil, err
		}
		videoCardDto := &VideoCardDto{}
		err = uc.copierService.Copy(videoCardDto, card.VideoCard)
		if err != nil {
			return nil, err
		}
		cardDto.VideoCard = videoCardDto
		cardDto.VideoCard.Video = card.VideoCard.Video
		cardDto.VideoCard.VideoLite = card.VideoCard.VideoLite
		cardDto.VideoCard.VideoPreview = card.VideoCard.VideoPreview
		cardDto.VideoCard.VideoPreviewBlur = card.VideoCard.VideoPreviewBlur

	case "regular":
		if card.RegularCardId == nil {
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}
		}
		if card.RegularCard.VideoId != nil {
			card.RegularCard.Video.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.VideoId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		if card.RegularCard.VideoPreviewId != nil {
			card.RegularCard.VideoPreview.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.VideoPreviewId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		if card.RegularCard.VideoPreviewBlurId != nil {
			card.RegularCard.VideoPreviewBlur.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.VideoPreviewBlurId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		if card.RegularCard.VideoLiteId != nil {
			card.RegularCard.VideoLite.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.VideoLiteId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		if card.RegularCard.UserId != nil && card.RegularCard.User.Picture != nil {
			card.RegularCard.User.Picture.Filepath, err = uc.mediaProvider.GetUrlById(*card.RegularCard.User.PictureId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}

		var tags []TagDto

		for _, tag := range card.RegularCard.RegularCardsTags {
			tagDto := &TagDto{}
			err = uc.copierService.Copy(tagDto, &tag.Tag)
			if err != nil {
				return nil, err
			}
			tags = append(tags, *tagDto)
		}
		err = uc.copierService.Copy(cardDto, card)
		if err != nil {
			return nil, err
		}
		regularCardDto := &RegularCardDto{}
		err = uc.copierService.Copy(regularCardDto, card.RegularCard)
		if err != nil {
			return nil, err
		}
		cardDto.RegularCard = regularCardDto
		cardDto.RegularCard.Tags = tags

	case "photo":
		if card.PhotoCard == nil {
			err = uc.copierService.Copy(cardDto, card)
			if err != nil {
				return nil, err
			}
		}
		if card.PhotoCard.PictureId != nil {
			card.PhotoCard.Picture.Filepath, err = uc.mediaProvider.GetUrlById(*card.PhotoCard.PictureId)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting picture media url")
				return
			}
		}
		err = uc.copierService.Copy(cardDto, card)
		if err != nil {
			return nil, err
		}
		photoCardDto := &PhotoCardDto{}
		err = uc.copierService.Copy(photoCardDto, card.PhotoCard)
		if err != nil {
			return nil, err
		}
		cardDto.PhotoCard = photoCardDto
		cardDto.PhotoCard.Picture = card.PhotoCard.Picture
	case "form":
		err = uc.copierService.Copy(cardDto, card)
		if err != nil {
			return nil, err
		}
		if cardDto.FormCard.Background != nil {
			cardDto.FormCard.Background.Filepath, err = uc.mediaProvider.GetUrlById(cardDto.FormCard.Background.Id)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting background media url")
				return
			}
		}
		if cardDto.FormCard.User != nil && cardDto.FormCard.User.Picture != nil {
			cardDto.FormCard.User.Picture.Filepath, err = uc.mediaProvider.GetUrlById(cardDto.FormCard.User.Picture.Id)
			if err != nil {
				err = errors.BadRequest.Wrap(err, "error getting video media url")
				return
			}
		}
		var tags []TagDto

		for _, tag := range card.FormCard.FormCardsTags {
			tagDto := &TagDto{}
			err = uc.copierService.Copy(tagDto, &tag.Tag)
			if err != nil {
				return nil, err
			}
			tags = append(tags, *tagDto)

		}
		cardDto.FormCard.Tags = tags
	}

	return cardDto, nil
}

func (uc CardUseCase) GetDtosFromGalleriesCard(gallery []entity.GalleriesCards) (cardDtos []CardDto, err error) {
	for _, card := range gallery {
		cardDto := &CardDto{}

		switch card.Card.Type {
		case "html":
			if card.Card.HtmlCardId == nil {
				err = uc.copierService.Copy(cardDto, card.Card)
				if err != nil {
					return nil, err
				}
				cardDto.Position = card.Position
				continue
			}
			err = uc.copierService.Copy(cardDto, card.Card)
			if err != nil {
				return nil, err
			}
			htmlCardDto := &HtmlCardDto{}
			err = uc.copierService.Copy(htmlCardDto, card.Card.HtmlCard)
			if err != nil {
				return nil, err
			}
			cardDto.Position = card.Position
		case "video":
			if card.Card.VideoCardId == nil {
				err = uc.copierService.Copy(cardDto, card.Card)
				if err != nil {
					return nil, err
				}
				cardDto.Position = card.Position
				continue
			}
			if card.Card.VideoCard.VideoId != nil {
				card.Card.VideoCard.Video.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.VideoCard.VideoId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.Card.VideoCard.VideoPreviewId != nil {
				card.Card.VideoCard.VideoPreview.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.VideoCard.VideoPreviewId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.Card.VideoCard.VideoPreviewBlurId != nil {
				card.Card.VideoCard.VideoPreviewBlur.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.VideoCard.VideoPreviewBlurId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.Card.VideoCard.VideoLiteId != nil {
				card.Card.VideoCard.VideoLite.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.VideoCard.VideoLiteId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}

			err = uc.copierService.Copy(cardDto, card.Card)
			if err != nil {
				return nil, err
			}
			videoCardDto := &VideoCardDto{}
			err = uc.copierService.Copy(videoCardDto, card.Card.VideoCard)
			if err != nil {
				return nil, err
			}
			cardDto.Position = card.Position
		case "photo":
			if card.Card.PhotoCardId == nil {
				err = uc.copierService.Copy(cardDto, card.Card)
				if err != nil {
					return nil, err
				}
				cardDto.Position = card.Position
				continue
			}
			if card.Card.PhotoCard.PictureId != nil {
				card.Card.PhotoCard.Picture.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.PhotoCard.PictureId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			err = uc.copierService.Copy(cardDto, card.Card)
			if err != nil {
				return nil, err
			}
			cardDto.PhotoCard = &PhotoCardDto{
				Picture: card.Card.PhotoCard.Picture,
			}
			if err != nil {
				return nil, err
			}
			cardDto.Position = card.Position
		case "form":
			err = uc.copierService.Copy(cardDto, card.Card)
			if err != nil {
				return nil, err
			}
			if cardDto.FormCard.Background != nil {
				cardDto.FormCard.Background.Filepath, err = uc.mediaProvider.GetUrlById(cardDto.FormCard.Background.Id)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting background media url")
					return
				}
			}
			if cardDto.FormCard.User != nil && cardDto.FormCard.User.Picture != nil {
				cardDto.FormCard.User.Picture.Filepath, err = uc.mediaProvider.GetUrlById(cardDto.FormCard.User.Picture.Id)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			var tags []TagDto

			for _, tag := range card.Card.FormCard.FormCardsTags {
				tagDto := &TagDto{}
				err = uc.copierService.Copy(tagDto, &tag.Tag)
				if err != nil {
					return nil, err
				}
				tags = append(tags, *tagDto)

			}
			cardDto.FormCard.Tags = tags
			cardDto.Position = card.Position
		case "regular":
			if card.Card.RegularCardId == nil {
				err = uc.copierService.Copy(cardDto, card)
				if err != nil {
					return nil, err
				}
				cardDto.Position = card.Position
				continue
			}
			if card.Card.RegularCard.VideoId != nil {
				card.Card.RegularCard.Video.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.RegularCard.VideoId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.Card.RegularCard.VideoPreviewId != nil {
				card.Card.RegularCard.VideoPreview.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.RegularCard.VideoPreviewId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.Card.RegularCard.VideoPreviewBlurId != nil {
				card.Card.RegularCard.VideoPreviewBlur.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.RegularCard.VideoPreviewBlurId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.Card.RegularCard.VideoLiteId != nil {
				card.Card.RegularCard.VideoLite.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.RegularCard.VideoLiteId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}
			if card.Card.RegularCard.UserId != nil && card.Card.RegularCard.User.Picture != nil {
				card.Card.RegularCard.User.Picture.Filepath, err = uc.mediaProvider.GetUrlById(*card.Card.RegularCard.User.PictureId)
				if err != nil {
					err = errors.BadRequest.Wrap(err, "error getting video media url")
					return
				}
			}

			var tags []TagDto

			for _, tag := range card.Card.RegularCard.RegularCardsTags {
				tagDto := &TagDto{}
				err = uc.copierService.Copy(tagDto, &tag.Tag)
				if err != nil {
					return nil, err
				}
				tags = append(tags, *tagDto)
			}
			err = uc.copierService.Copy(cardDto, card.Card)
			if err != nil {
				return nil, err
			}
			regularCardDto := &RegularCardDto{}
			err = uc.copierService.Copy(regularCardDto, card.Card.RegularCard)
			if err != nil {
				return nil, err
			}
			cardDto.RegularCard.Tags = tags
			cardDto.Position = card.Position
		}
		cardDtos = append(cardDtos, *cardDto)
	}
	return cardDtos, nil
}

func (uc CardUseCase) GetCardFromUpdateDto(cardRequest *UpdateCardRequest) (
	card *entity.Card, err error,
) {
	card = &entity.Card{}

	switch cardRequest.Type {
	case "html":
		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		htmlCardDto := &entity.HtmlCard{}
		err = uc.copierService.Copy(htmlCardDto, cardRequest.HtmlCard)
		if err != nil {
			return nil, err
		}
		card.HtmlCard = htmlCardDto
	case "video":
		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		videoCard := &entity.VideoCard{}
		err = uc.copierService.Copy(videoCard, cardRequest.VideoCard)
		if err != nil {
			return nil, err
		}
		card.VideoCard = videoCard
		card.VideoCard.VideoId = cardRequest.VideoCard.VideoId
		card.VideoCard.VideoLiteId = cardRequest.VideoCard.VideoLiteId
		card.VideoCard.VideoPreviewId = cardRequest.VideoCard.VideoPreviewId
		card.VideoCard.VideoPreviewBlurId = cardRequest.VideoCard.VideoPreviewBlurId
	case "regular":
		var tags []entity.RegularCardsTags
		for _, tagId := range cardRequest.RegularCard.Tags {
			tags = append(
				tags, entity.RegularCardsTags{
					TagID: tagId,
				},
			)
		}
		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		regularCard := &entity.RegularCard{}
		err = uc.copierService.Copy(regularCard, cardRequest.RegularCard)
		if err != nil {
			return nil, err
		}
		card.RegularCard = regularCard
		card.RegularCard.RegularCardsTags = tags

	case "photo":
		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		photoCard := &entity.PhotoCard{}
		err = uc.copierService.Copy(photoCard, cardRequest.PhotoCard)
		if err != nil {
			return nil, err
		}
		card.PhotoCard = photoCard
		card.PhotoCard.PictureId = cardRequest.PhotoCard.PictureId
	case "form":
		var tags []entity.FormCardsTags
		for _, tagId := range cardRequest.FormCard.Tags {
			tags = append(
				tags, entity.FormCardsTags{
					TagID: tagId,
				},
			)
		}
		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		formCard := &entity.FormCard{}
		err = uc.copierService.Copy(formCard, cardRequest.FormCard)
		if err != nil {
			return nil, err
		}
		card.FormCard = formCard
		card.FormCard.FormCardsTags = tags
	}
	return card, nil
}

func (uc CardUseCase) GetGalleriesCardsFromIds(
	ctx context.Context, cardId uuid.UUID, galleryIds []uuid.UUID,
) []entity.GalleriesCards {
	var galleriesCards []entity.GalleriesCards
	for _, id := range galleryIds {
		lastPosition, err := uc.cardRepository.GetLastPositionInGalley(ctx, id)
		if err != nil {
			lastPosition = 0
		}
		newPosition := lastPosition + 1
		galleriesCards = append(
			galleriesCards, entity.GalleriesCards{
				GalleryID: &id,
				CardID:    &cardId,
				Position:  newPosition,
			},
		)
	}
	return galleriesCards
}

func (uc CardUseCase) GetCardFromDto(cardRequest CreateCardRequest) (
	card *entity.Card, err error,
) {
	card = &entity.Card{}
	newCardId := uuid.New()
	switch cardRequest.Type {
	case "html":
		htmlCardId := uuid.New()
		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		htmlCard := &entity.HtmlCard{}
		err = uc.copierService.Copy(htmlCard, cardRequest.HtmlCard)
		if err != nil {
			return nil, err
		}
		card.ID = newCardId
		card.HtmlCard = htmlCard
		card.HtmlCardId = &htmlCardId
		card.HtmlCard.ID = htmlCardId
	case "video":
		videoCardId := uuid.New()

		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		videoCard := &entity.VideoCard{}
		err = uc.copierService.Copy(videoCard, cardRequest.VideoCard)
		if err != nil {
			return nil, err
		}
		card.ID = newCardId
		card.VideoCard = videoCard
		card.VideoCardId = &videoCardId
		card.VideoCard.ID = videoCardId
	case "regular":
		regularCardId := uuid.New()

		var tags []entity.RegularCardsTags
		for _, tagId := range cardRequest.RegularCard.Tags {
			tags = append(
				tags, entity.RegularCardsTags{
					RegularCardID: regularCardId,
					TagID:         tagId,
				},
			)
		}

		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		regularCard := &entity.RegularCard{}
		err = uc.copierService.Copy(regularCard, cardRequest.RegularCard)
		if err != nil {
			return nil, err
		}
		card.ID = newCardId
		card.RegularCard = regularCard
		card.RegularCardId = &regularCardId
		card.RegularCard.ID = regularCardId
		card.RegularCard.RegularCardsTags = tags
	case "photo":
		photoCardId := uuid.New()
		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		photoCard := &entity.PhotoCard{}
		err = uc.copierService.Copy(photoCard, cardRequest.PhotoCard)
		if err != nil {
			return nil, err
		}
		card.ID = newCardId
		card.PhotoCard = photoCard
		card.PhotoCardId = &photoCardId
		card.PhotoCard.ID = photoCardId
	case "form":
		formCardId := uuid.New()
		var tags []entity.FormCardsTags
		for _, tagId := range cardRequest.FormCard.Tags {
			tags = append(
				tags, entity.FormCardsTags{
					FormCardID: formCardId,
					TagID:      tagId,
				},
			)
		}
		err = uc.copierService.Copy(card, cardRequest)
		if err != nil {
			return nil, err
		}
		formCard := &entity.FormCard{}
		err = uc.copierService.Copy(formCard, cardRequest.FormCard)
		if err != nil {
			return nil, err
		}
		card.FormCard = formCard
		card.FormCard.FormCardsTags = tags

		card.ID = newCardId
		//card.FormCard = formCard
		card.FormCardId = &formCardId
		card.FormCard.ID = formCardId
		//card.FormCard.FormCardsTags = tags
	}

	return card, nil
}

func (uc CardUseCase) GetCardFromDtos(cardsRequest []CreateCardRequest, galleryId uuid.UUID) (
	cards []entity.GalleriesCards, err error,
) {
	for _, card := range cardsRequest {
		newCardId := uuid.New()
		galleriesCard := &entity.GalleriesCards{}

		switch card.Type {
		case "html":
			htmlCardId := uuid.New()

			err = uc.copierService.Copy(galleriesCard, card)
			if err != nil {
				return nil, err
			}
			cardEntity := &entity.Card{}
			err = uc.copierService.Copy(cardEntity, card)
			if err != nil {
				return nil, err
			}
			htmlCard := &entity.HtmlCard{}
			err = uc.copierService.Copy(htmlCard, card.HtmlCard)
			if err != nil {
				return nil, err
			}
			galleriesCard.GalleryID = &galleryId
			galleriesCard.CardID = &newCardId
			galleriesCard.Card = cardEntity
			galleriesCard.Card.ID = newCardId
			galleriesCard.Card.HtmlCard = htmlCard
			galleriesCard.Card.HtmlCardId = &htmlCardId
			galleriesCard.Card.HtmlCard.ID = htmlCardId
		case "video":
			videoCardId := uuid.New()

			err = uc.copierService.Copy(galleriesCard, card)
			if err != nil {
				return nil, err
			}
			cardEntity := &entity.Card{}
			err = uc.copierService.Copy(cardEntity, card)
			if err != nil {
				return nil, err
			}
			videoCard := &entity.VideoCard{}
			err = uc.copierService.Copy(videoCard, card.VideoCard)
			if err != nil {
				return nil, err
			}
			galleriesCard.GalleryID = &galleryId
			galleriesCard.CardID = &newCardId
			galleriesCard.Card = cardEntity
			galleriesCard.Card.ID = newCardId
			galleriesCard.Card.VideoCard = videoCard
			galleriesCard.Card.VideoCardId = &videoCardId
			galleriesCard.Card.VideoCard.ID = videoCardId
		case "regular":
			regularCardId := uuid.New()
			var tags []entity.RegularCardsTags
			for _, tagId := range card.RegularCard.Tags {
				tags = append(
					tags, entity.RegularCardsTags{
						RegularCardID: regularCardId,
						TagID:         tagId,
					},
				)
			}

			err = uc.copierService.Copy(galleriesCard, card)
			if err != nil {
				return nil, err
			}
			cardEntity := &entity.Card{}
			err = uc.copierService.Copy(cardEntity, card)
			if err != nil {
				return nil, err
			}
			regularCard := &entity.RegularCard{}
			err = uc.copierService.Copy(regularCard, card.RegularCard)
			if err != nil {
				return nil, err
			}
			regularCard.RegularCardsTags = tags

			galleriesCard.GalleryID = &galleryId
			galleriesCard.CardID = &newCardId
			galleriesCard.Card = cardEntity
			galleriesCard.Card.ID = newCardId
			galleriesCard.Card.RegularCard = regularCard
			galleriesCard.Card.RegularCardId = &regularCardId
			galleriesCard.Card.RegularCard.ID = regularCardId
		case "photo":
			photoCardId := uuid.New()
			err = uc.copierService.Copy(galleriesCard, card)
			if err != nil {
				return nil, err
			}
			cardEntity := &entity.Card{}
			err = uc.copierService.Copy(cardEntity, card)
			if err != nil {
				return nil, err
			}
			photoCard := &entity.PhotoCard{}
			err = uc.copierService.Copy(photoCard, card.PhotoCard)
			if err != nil {
				return nil, err
			}
			galleriesCard.GalleryID = &galleryId
			galleriesCard.CardID = &newCardId
			galleriesCard.Card = cardEntity
			galleriesCard.Card.ID = newCardId
			galleriesCard.Card.PhotoCard = photoCard
			galleriesCard.Card.PhotoCardId = &photoCardId
			galleriesCard.Card.PhotoCard.ID = photoCardId
		case "form":
			formCardId := uuid.New()
			var tags []entity.FormCardsTags
			for _, tagId := range card.FormCard.Tags {
				tags = append(
					tags, entity.FormCardsTags{
						FormCardID: formCardId,
						TagID:      tagId,
					},
				)
			}

			err = uc.copierService.Copy(galleriesCard, card)
			if err != nil {
				return nil, err
			}
			cardEntity := &entity.Card{}
			err = uc.copierService.Copy(cardEntity, card)
			if err != nil {
				return nil, err
			}
			formCard := &entity.FormCard{}
			err = uc.copierService.Copy(formCard, card.FormCard)
			if err != nil {
				return nil, err
			}
			formCard.FormCardsTags = tags

			galleriesCard.GalleryID = &galleryId
			galleriesCard.CardID = &newCardId
			galleriesCard.Card = cardEntity
			galleriesCard.Card.ID = newCardId
			galleriesCard.Card.FormCard = formCard
			galleriesCard.Card.FormCardId = &formCardId
			galleriesCard.Card.FormCard.ID = formCardId
		}
		cards = append(cards, *galleriesCard)
	}

	return cards, nil
}

func (uc CardUseCase) getTagsFromDtos(dtos []TagDtoRequest) ([]entity.Tag, error) {
	var tags []entity.Tag
	for _, dto := range dtos {
		tag := &entity.Tag{}
		err := uc.copierService.Copy(tag, dto)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}
