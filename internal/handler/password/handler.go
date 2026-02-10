package password

import (
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
