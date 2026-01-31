package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/toanuitt/bookmark_service/internal/service"
	"github.com/toanuitt/bookmark_service/pkg/response"
)

var (
	MissingInputErr = errors.New("missing required fields")
)

type Userhandler interface {
	RegisterUser(c *gin.Context)
}

type user struct {
	svc service.Userservice
}

func NewUser(svc service.Userservice) Userhandler {
	return &user{svc: svc}
}

// registerInputBody represents the request payload for user registration
type registerInputBody struct {
	Username    string `json:"username"  binding:"required"`
	Password    string `json:"password"  binding:"required"`
	DisplayName string `json:"display_name"  binding:"required"`
	Email       string `json:"email"  binding:"required,email"`
}

// RegisterUser handles user registration.
//
// @Summary     Register a new user
// @Description Create a new user account using username, password, display name, and email.
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       user  body      registerInputBody  true  "User registration payload"
// @Success     200   {object}  response.Response  "Register success"
// @Failure     400   {object}  response.Response  "Invalid request or missing required fields"
// @Failure     500   {object}  response.Response  "Internal server error"
// @Router      /v1/users/register [post]
func (h *user) RegisterUser(c *gin.Context) {
	input := &registerInputBody{}
	if err := c.ShouldBindJSON(input); err != nil {
		log.Error().Str("url", c.Request.URL.String()).Err(err).Msg("Invalid request body on RegisterUser")
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	if input.Username == "" || input.Password == "" || input.DisplayName == "" || input.Email == "" {
		log.Error().Str("url", c.Request.URL.String()).Err(MissingInputErr).Msg("Missing required fields on RegisterUser")
		c.JSON(http.StatusBadRequest, response.InputFieldError(MissingInputErr))
		return
	}

	res, err := h.svc.Register(c, input.Username, input.Password, input.DisplayName, input.Email)
	if err != nil {
		log.Error().Str("url", c.Request.URL.String()).Err(err).Msg("Service return error on RegisterUser")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
