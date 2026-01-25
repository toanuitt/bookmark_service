package main

import (
	_ "github.com/toanuitt/bookmark_service/docs"
	"github.com/toanuitt/bookmark_service/internal/api"
	"github.com/toanuitt/bookmark_service/pkg/logger"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
)

//	@title			BookMark_Service API
//	@version		1.0
//	@description	API documentation for bookmark service
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	logger.SetLogLevel()
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	redisClient, err := redisPkg.NewRedisClient("")
	if err != nil {
		panic(err)
	}
	app := api.New(cfg, redisClient)
	app.Start()
}
