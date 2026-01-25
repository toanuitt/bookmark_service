package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/toanuitt/bookmark_service/internal/handler"
	"github.com/toanuitt/bookmark_service/internal/repository"
	"github.com/toanuitt/bookmark_service/internal/service"
	"github.com/toanuitt/bookmark_service/pkg/stringutils"
)

// Engine defines the interface for the API engine.
// It exposes methods for starting the server and handling HTTP requests.
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// api is the concrete implementation of the Engine interface.
// It manages the Gin HTTP engine and application configuration.
type api struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

// New creates and returns a new Engine instance with initialized routes and handlers.
// It takes a Config parameter to configure the API engine.
func New(cfg *Config, redisClient *redis.Client) Engine {
	a := &api{
		app:         gin.New(),
		cfg:         cfg,
		redisClient: redisClient,
	}
	a.RegisterEP()
	return a
}

// Start starts the HTTP server on the configured port and blocks until an error occurs.
func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

// ServeHTTP implements the http.Handler interface, allowing the API engine to serve HTTP requests.
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

// RegisterEP registers all API endpoints and their corresponding handlers.
// It initializes service and handler dependencies and sets up routes for password generation,
// health checks, and Swagger documentation.
// Endpoints:
//   - GET /gen-pass: Generates a random password
//   - GET /health: Health check endpoint
//   - POST /v1/links/shorten: Shorten URL endpoint
//   - GET /swagger/*any: Swagger UI documentation
func (a *api) RegisterEP() {

	//Repository
	urlRepo := repository.NewUrlStorage(a.redisClient)
	healthCheckRepo := repository.NewHealthCheckStorage(a.redisClient)

	//service
	passSvc := service.NewPassword()
	healthSvc := service.NewHealthCheck(a.cfg.ServiceName, a.cfg.InstanceID, healthCheckRepo)
	urlshortenSvc := service.NewShortenURL(urlRepo, stringutils.NewCodeGenerator())

	//handler
	passHandler := handler.NewPassword(passSvc)
	healthHandler := handler.NewHealthCheck(healthSvc)
	urlshortenHandler := handler.NewShortenURL(urlshortenSvc)

	//Router
	a.app.GET("/gen-pass", passHandler.GenPass)
	a.app.GET("/health-check", healthHandler.CheckHealth)
	routers := a.app.Group("/v1")
	{
		routers.POST("/links/shorten", urlshortenHandler.ShortenURL)
		routers.GET("/links/redirect/:code", urlshortenHandler.GetURL)
	}
	// Register Swagger documentation endpoint
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
