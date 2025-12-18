package handler

import (
	"strconv"

	_ "searchav/internal/dto"
	"searchav/internal/service"
)

// DetailHandler handles video detail requests
type DetailHandler struct {
	service *service.DetailService
}

// NewDetailHandler creates a new detail handler
func NewDetailHandler(service *service.DetailService) *DetailHandler {
	return &DetailHandler{
		service: service,
	}
}

// GetDetail handles video detail requests
// @Summary Get video detail
// @Description Get video details and play URLs from a specific source
// @Tags detail
// @Accept json
// @Produce json
// @Param source query string true "Source code"
// @Param id query int true "Video ID"
// @Success 200 {object} dto.DetailResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /detail [get]
func (h *DetailHandler) GetDetail(ctx *Context) error {
	sourceCode := ctx.Query("source")
	if sourceCode == "" {
		return ctx.BadRequest("missing source parameter")
	}

	idStr := ctx.Query("id")
	if idStr == "" {
		return ctx.BadRequest("missing id parameter")
	}

	vodID, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.BadRequest("invalid id parameter")
	}

	detail, err := h.service.GetDetail(ctx.Context(), sourceCode, vodID)
	if err != nil {
		return ctx.InternalError(err)
	}

	return ctx.SuccessWithData(detail)
}
