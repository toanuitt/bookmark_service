package url

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/toanuitt/bookmark_service/internal/handler/dto"
	"github.com/toanuitt/bookmark_service/internal/handler/utils"
	"github.com/toanuitt/bookmark_service/pkg/response"
)

// @Summary Shorten URL
// @Description Shortens a given URL and returns a shortened URL code.
// @Tags URL Shortener
// @Accept json
// @Produce json
// @Param request body dto.ShortenURLRequest true "URL to shorten"
// @Success 200 {object} dto.ShortenURLResponse
// @Failure 400 {object} map[string]string "Bad Request - invalid URL or validation error"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /v1/links/shorten [post]
func (h *Url) ShortenURL(c *gin.Context) {
	req, err := utils.BindInputFromRequest[dto.ShortenURLRequest](c)
	if err != nil {
		return
	}

	code, err := h.svc.ShortlengthURL(c, req.URL, req.ExpireIn)
	if err != nil {
		log.Error().Str("url", req.URL).Err(err).Msg("Service return error on ShortenUrl")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, dto.ShortenURLResponse{
		Message: "Shorten URL generated successfully!",
		Code:    code,
	})
}
