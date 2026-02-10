package user

import (
	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/internal/service"
)

// Userhandler defines the interface for user-related HTTP handlers
type Userhandler interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

// user implements the Userhandler interface
type user struct {
	svc service.Userservice
}

// NewUser creates a new user handler instance
func NewUser(svc service.Userservice) Userhandler {
	return &user{svc: svc}
}
