package main

import (
	_ "github.com/toanuitt/bookmark_service/docs"
	"github.com/toanuitt/bookmark_service/internal/api"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/pkg/common"
	"github.com/toanuitt/bookmark_service/pkg/jwtutils"
	"github.com/toanuitt/bookmark_service/pkg/logger"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
	sqldbPkg "github.com/toanuitt/bookmark_service/pkg/sqldb"
)

//	@title			BookMark_Service API
//	@version		1.2
//	@description	API documentation for bookmark service
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your Bearer token in the format: Bearer {token}
func main() {
	logger.SetLogLevel()
	cfg, err := api.NewConfig()
	common.HandleError(err)

	redisClient, err := redisPkg.NewRedisClient("")
	common.HandleError(err)

	db, err := sqldbPkg.NewClient("")
	common.HandleError(err)
	common.HandleError(db.AutoMigrate(&model.User{}))

	jwtGen, err := jwtutils.NewJWTGenerator("./private.pem")
	common.HandleError(err)
	jwtValidator, err := jwtutils.NewJWTValidator("./public.pem")
	common.HandleError(err)

	app := api.New(cfg, redisClient, db, jwtGen, jwtValidator)
	app.Start()
}
