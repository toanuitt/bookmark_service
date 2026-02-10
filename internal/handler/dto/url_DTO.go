package dto

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
