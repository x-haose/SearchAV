package handler

import (
	_ "searchav/internal/dto"
	"searchav/internal/service"
)

// SearchHandler handles video search requests
type SearchHandler struct {
	service *service.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(service *service.SearchService) *SearchHandler {
	return &SearchHandler{
		service: service,
	}
}

// Search handles video search requests
// @Summary Search videos
// @Description Aggregate search across multiple video sources
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search keyword"
// @Param adult query string false "Include adult sources (1=yes, 0=no, default=0)"
// @Success 200 {object} dto.SearchResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /search [get]
func (h *SearchHandler) Search(ctx *Context) error {
	keyword := ctx.Query("q")
	if keyword == "" {
		return ctx.BadRequest("missing search keyword")
	}

	// Parse adult filter parameter (default: false)
	includeAdult := ctx.Query("adult") == "1"

	ctx.Logger.Info().Str("keyword", keyword).Bool("adult", includeAdult).Msg("search request received")

	list, err := h.service.Search(ctx.Context(), keyword, includeAdult)
	if err != nil {
		ctx.Logger.Error().Err(err).Msg("search failed")
		return ctx.InternalError(err)
	}

	ctx.Logger.Info().Int("count", len(list)).Msg("search completed")

	return ctx.SuccessWithList(list)
}
