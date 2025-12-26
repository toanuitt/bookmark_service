package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/internal/service"
)

// passwordHandler is the concrete implementation of the PassWord handler interface.
type passwordHandler struct {
	svc service.Password
}

// PassWord defines the interface for password generation HTTP handlers.
type PassWord interface {
	GenPass(c *gin.Context)
}

// NewPassword creates and returns a new PassWord handler instance.
// It takes a service.Password dependency to generate secure passwords.
func NewPassword(svc service.Password) PassWord {
	return &passwordHandler{svc: svc}
}

// @Summary Generate a random password
// @Description Generates a cryptographically secure random password
// @Tags password
// @Produce plain
// @Success 200 {string} string "Generated password"
// @Failure 500 {string} string "Error message"
// @Router /gen-pass [get]
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, "err")
		return
	}
	c.String(http.StatusOK, pass)
}
