package url

import (
	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/internal/service"
)

// urlHandler defines the interface for handling URL shortening HTTP requests.
type urlHandler interface {
	ShortenURL(c *gin.Context)
	GetURL(c *gin.Context)
}

// shortenUrl is the concrete implementation of the ShortenURLhandler interface.
type Url struct {
	svc service.ShortenURLservice
}

// NewShortenURL creates and returns a new ShortenURLhandler instance.
func NewShortenURL(shortenUrlSvc service.ShortenURLservice) urlHandler {
	return &Url{svc: shortenUrlSvc}
}
