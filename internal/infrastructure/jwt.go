package infrastructure

import (
	"github.com/toanuitt/bookmark_service/pkg/common"
	"github.com/toanuitt/bookmark_service/pkg/jwtutils"
)

func CreateJWTProvider() (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGen, err := jwtutils.NewJWTGenerator("./private.pem")
	common.HandleError(err)
	jwtValidator, err := jwtutils.NewJWTValidator("./public.pem")
	common.HandleError(err)
	return jwtGen, jwtValidator
}
