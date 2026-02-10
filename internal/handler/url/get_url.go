package url

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/toanuitt/bookmark_service/internal/service"
	"github.com/toanuitt/bookmark_service/pkg/response"
)

// GetUrl shortens a given URL and returns a shortened URL code.
// @Summary Get URL
// @Description Get URL by code
// @Tags URL Shortener
// @Accept json
// @Produce json
// @Param code path string true "Url code" Format(string)
// @Success 302
// @Failure 400  "Bad Request - invalid URL or validation error"
// @Failure 404  "URL not found"
// @Failure 500  "Internal Server Error"
// @Router /v1/links/redirect/{code} [get]
func (h *Url) GetURL(c *gin.Context) {
	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid url code",
		})
		return
	}

	url, err := h.svc.GetURL(c, code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "url not found",
			})
			return
		}

		log.Error().Err(err).Msg("Service return error on GetURL")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.Redirect(http.StatusFound, url)
}
