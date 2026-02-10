package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toanuitt/bookmark_service/internal/api"
	"github.com/toanuitt/bookmark_service/pkg/common"
	"github.com/toanuitt/bookmark_service/pkg/logger"
)

func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	common.HandleError(err)
	if cfg.InstanceID == "" {
		id, err := uuid.NewV7()
		common.HandleError(err)
		cfg.InstanceID = id.String()
	}
	return cfg
}

func CreateAPI() api.Engine {
	logger.SetLogLevel()
	cfg := CreateAPIConfig()
	redisClient := CreateRedisConn()
	sqlDB := CreateSQLDBWithMigration()
	jwtGen, jwtVal := CreateJWTProvider()
	app := gin.New()
	return api.New(&api.EngineOpts{
		Engine:       app,
		Cfg:          cfg,
		Redis:        redisClient,
		SqlDB:        sqlDB,
		JWTGenerator: jwtGen,
		JWTValidator: jwtVal,
	})

}
