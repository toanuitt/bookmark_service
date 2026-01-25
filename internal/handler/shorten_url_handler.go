package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/internal/service"
)

// ShortenURLhandler defines the interface for handling URL shortening HTTP requests.
type ShortenURLhandler interface {
	ShortenURL(c *gin.Context)
}

// shortenUrl is the concrete implementation of the ShortenURLhandler interface.
type shortenUrl struct {
	svc service.ShortenURLservice
}

// ShortenURLRequest represents the JSON request body for the URL shortening endpoint.
type ShortenURLRequest struct {
	URL      string `json:"url" binding:"required,url" example:"https://example.com"`
	ExpireIn int64  `json:"exp" binding:"required,gt=0,lte=3600" example:"3600"`
}

// ShortenURLResponse represents the JSON response body for the URL shortening endpoint.
type ShortenURLResponse struct {
	Message string `json:"message" example:"Shorten URL generated successfully!"`
	Code    string `json:"code" example:"abc123"`
}

// NewShortenURL creates and returns a new ShortenURLhandler instance.
func NewShortenURL(shortenUrlSvc service.ShortenURLservice) ShortenURLhandler {
	return &shortenUrl{svc: shortenUrlSvc}
}

// @Summary Shorten URL
// @Description Shortens a given URL and returns a shortened URL code.
// @Tags URL Shortener
// @Accept json
// @Produce json
// @Param request body handler.ShortenURLRequest true "URL to shorten"
// @Success 200 {object} handler.ShortenURLResponse
// @Failure 400 {object} map[string]string "Bad Request - invalid URL or validation error"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /v1/links/shorten [post]
func (h *shortenUrl) ShortenURL(c *gin.Context) {
	req := &ShortenURLRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
		return
	}

	code, err := h.svc.ShortlengthURL(c, req.URL, req.ExpireIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, ShortenURLResponse{
		Message: "Shorten URL generated successfully!",
		Code:    code,
	})
}
