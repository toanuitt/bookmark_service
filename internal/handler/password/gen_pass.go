package password

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary Generate a random password
// @Description Generates a cryptographically secure random password
// @Tags password
// @Produce plain
// @Success 200 {string} string "Generated password"
// @Failure 500 {string} string "Error message"
// @Router /v1/gen-pass [get]
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		log.Error().Err(err).Msg("Service return error on GenPass")
		c.String(http.StatusInternalServerError, "err")
		return
	}
	c.String(http.StatusOK, pass)
}
