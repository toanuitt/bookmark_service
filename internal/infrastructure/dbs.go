package infrastructure

import (
	"github.com/redis/go-redis/v9"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/pkg/common"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
	"github.com/toanuitt/bookmark_service/pkg/sqldb"
	"gorm.io/gorm"
)

func CreateRedisConn() *redis.Client {
	redisClient, err := redisPkg.NewRedisClient("")
	common.HandleError(err)
	return redisClient
}

func CreateSQLDBWithMigration() *gorm.DB {
	db, err := sqldb.NewClient("")
	common.HandleError(err)
	common.HandleError(db.AutoMigrate(&model.User{}))
	return db
}
