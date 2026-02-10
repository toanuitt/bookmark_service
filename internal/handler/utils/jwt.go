package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("Invalid Token")
	ErrEmptyUID     = errors.New("empty uid")
)

func GetJWTClaimsFromRequest(c *gin.Context) (jwt.MapClaims, error) {
	//GetURL the user id from jwt Token
	tokenInfo, _ := c.Get("claims")
	claims, valid := tokenInfo.(jwt.MapClaims)
	if !valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func GetUIDFromRequest(c *gin.Context) (string, error) {
	claims, err := GetJWTClaimsFromRequest(c)
	if err != nil {
		return "", err
	}
	uid, ok := claims["sub"].(string)
	if !ok || uid == "" {
		return "", ErrEmptyUID
	}
	return uid, nil

}
