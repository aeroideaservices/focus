package handlers

import (
	"github.com/aeroideaservices/focus/services/errors"
	middleware "github.com/aeroideaservices/focus/services/gin-middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jemzee04/focus/page/plugin/actions"
	"github.com/jemzee04/focus/page/rest/services"
	"net/http"
)

type CardHandler struct {
	cardUseCase  *actions.CardUseCase
	errorHandler *middleware.ErrorHandler
	validator    services.Validator
}

func NewCardHandler(
	cardUseCase *actions.CardUseCase, errorHandler *middleware.ErrorHandler, validator services.Validator,
) *CardHandler {
	return &CardHandler{
		cardUseCase:  cardUseCase,
		errorHandler: errorHandler,
		validator:    validator,
	}
}

func (h CardHandler) GetList(c *gin.Context) {
	searchValue := c.Query("searchValue")
	name := c.Query("name")
	pages, err := h.cardUseCase.GetListWithSearch(c, searchValue, name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pages)
}

func (h CardHandler) Create(c *gin.Context) {
	request := &actions.CreateCardRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	err := h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	resp, err := h.cardUseCase.Create(c, request)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == UniqueViolationErr {
				c.JSON(
					http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"error":   "Duplicate key value violates unique constraint",
						"message": "Карточка с таким кодом уже существует",
						"debug":   err.Error(),
					},
				)
				return
			}
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h CardHandler) GetById(c *gin.Context) {
	dto := actions.GetCardRequest{
		ID: c.Param("card-id"),
	}

	err := h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	card, err := h.cardUseCase.GetById(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if card.ID == uuid.Nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error while getting card"))
		return
	}

	c.JSON(http.StatusOK, card)
}

func (h CardHandler) Update(c *gin.Context) {
	request := &actions.UpdateCardRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	request.ID = cardId

	err = h.cardUseCase.Update(c, request)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == UniqueViolationErr {
				c.JSON(
					http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"error":   "Duplicate key value violates unique constraint",
						"message": "Карточка с таким кодом уже существует",
						"debug":   err.Error(),
					},
				)
				return
			}
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h CardHandler) Delete(c *gin.Context) {
	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err = h.cardUseCase.Delete(c, cardId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h CardHandler) PatchUser(c *gin.Context) {
	request := &actions.PatchUserRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	request.CardId = cardId

	err = h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.cardUseCase.PatchUser(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h CardHandler) PatchPreviewText(c *gin.Context) {
	request := &actions.PatchPreviewTextRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	request.CardId = cardId

	err = h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.cardUseCase.PatchPreviewText(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h CardHandler) PatchDetailText(c *gin.Context) {
	request := &actions.PatchDetailTextRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	request.CardId = cardId

	err = h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.cardUseCase.PatchDetailText(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h CardHandler) PatchLearnMoreUrl(c *gin.Context) {
	request := &actions.PatchLearnMoreUrlRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	request.CardId = cardId

	err = h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.cardUseCase.PatchLearnMoreUrl(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h CardHandler) PatchTags(c *gin.Context) {
	var request []actions.TagDtoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err = h.cardUseCase.PatchTags(c, cardId, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h CardHandler) LinkTags(c *gin.Context) {
	var tagStrIds []string
	if err := c.ShouldBindJSON(&tagStrIds); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	tagIds, err := services.GetIdsFromStrings(tagStrIds)
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err, isUniqViolErr := h.cardUseCase.LinkTags(c, cardId, tagIds)
	if err != nil {
		if isUniqViolErr {
			c.JSON(
				http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"error":   "Bad Request",
					"message": "Тег уже привязан к карточке",
					"debug":   err.Error(),
				},
			)
			return
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, "success")

}

func (h CardHandler) UnlinkTags(c *gin.Context) {
	var tagStrIds []string
	if err := c.ShouldBindJSON(&tagStrIds); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	tagIds, err := services.GetIdsFromStrings(tagStrIds)
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err = h.cardUseCase.UnlinkTags(c, cardId, tagIds)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")

}
