package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/internal/handler/utils"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/internal/service"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
	"github.com/toanuitt/bookmark_service/pkg/response"
)

var (
	// MissingInputErr indicates that required fields are missing from the request
	MissingInputErr = errors.New("missing required fields")
)

// Response message constants
const (
	// Success messages
	registerSuccessMessage       = "Register an user successfully!"
	loginSuccessMessage          = "Logged in successfully!"
	updateSelfInfoSuccessMessage = "Edit current user successfully!"

	// Error messages - Authentication
	invalidTokenError = "Invalid Token"

	// Error messages - Registration
	usernameTakenError = "username or email is already taken"

	// Error messages - Login
	invalidCredentialsError = "invalid username or password"

	// Error messages - Profile Update
	noUpdateDataError = "at least one field must be provided for update"
	emailTakenError   = "email is already taken"
)

// Userhandler defines the interface for user-related HTTP handlers
type Userhandler interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

// user implements the Userhandler interface
type user struct {
	svc service.Userservice
}

// NewUser creates a new user handler instance
func NewUser(svc service.Userservice) Userhandler {
	return &user{svc: svc}
}

// registerInputBody represents the request payload for user registration
type registerInputBody struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,gte=8"`
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

// registerUserData represents the user data returned after successful registration
type registerUserData struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	UpdatedAt   string `json:"updated_at"`
}

// registerResBody represents the response body for user registration
type registerResBody struct {
	Data    *registerUserData `json:"data"`
	Message string            `json:"message"`
}

// RegisterUser handles user registration.
//
// @Summary Register a new user
// @Description Create a new user account using username, password, display name, and email.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body registerInputBody true "User registration payload"
// @Success 200 {object} registerResBody
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/users/register [post]
func (h *user) RegisterUser(c *gin.Context) {
	input, err := utils.BindInputFromRequest[registerInputBody](c)
	if err != nil {
		return
	}

	res, err := h.svc.Register(c, input.Username, input.Password, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplicationType):
		c.JSON(http.StatusBadRequest, gin.H{"error": usernameTakenError})
		return
	case errors.Is(err, nil):
	default:
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, &registerResBody{
		Message: registerSuccessMessage,
		Data: &registerUserData{
			ID:          res.ID,
			Username:    res.Username,
			DisplayName: res.DisplayName,
			Email:       res.Email,
			UpdatedAt:   res.UpdatedAt.String(),
		},
	})
}

// loginInputBody represents the request payload for user login
type loginInputBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=8"`
}

// loginResBody represents the response body for user login
type loginResBody struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Login authenticates a user and returns a JWT token.
//
// @Summary Login a user
// @Description Authenticate a user with username and password, returns a JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param body body loginInputBody true "User login credentials"
// @Success 200 {object} loginResBody
// @Failure 400 {object} response.Response "Invalid username or password"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/users/login [post]
func (h *user) Login(c *gin.Context) {
	input, err := utils.BindInputFromRequest[loginInputBody](c)
	if err != nil {
		return
	}

	// call service
	token, err := h.svc.Login(c, input.Username, input.Password)
	switch {
	case errors.Is(err, service.ErrClientErr):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case errors.Is(err, dbutils.ErrNotFoundType):
		c.JSON(http.StatusNotFound, gin.H{"error": invalidCredentialsError})
		return
	case errors.Is(err, nil):
	default:
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	// return token
	c.JSON(http.StatusOK, &loginResBody{
		Message: loginSuccessMessage,
		Data:    token,
	})
}

// profileResBody represents the response body for user profile
type profileResBody struct {
	Data *model.User `json:"data"`
}

// GetProfile returns the current user's profile based on JWT token.
//
// @Summary Get current user profile
// @Description Get profile of the currently authenticated user using JWT in Authorization header
// @Tags Users
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Success 200 {object} profileResBody
// @Failure 401 {object} response.Response "Invalid or missing token"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/self/info [get]
func (h *user) GetProfile(c *gin.Context) {
	userIDvalue, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidTokenError})
		return
	}

	userID, ok := userIDvalue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidTokenError})
		return
	}

	res, err := h.svc.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, &profileResBody{
		Data: res,
	})
}

// updateProfileInputBody represents the request payload for profile updates
type updateProfileInputBody struct {
	DisplayName string `json:"display_name" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,email"`
}

// UpdateProfile handles updating the current user's profile.
//
// @Summary Update current user profile
// @Description Update display name and/or email of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body updateProfileInputBody true "Update profile request"
// @Success 200 {object} response.Response "Profile updated successfully"
// @Failure 400 {object} response.Response "No data provided for update"
// @Failure 401 {object} response.Response "Invalid or missing token"
// @Failure 404 {object} response.Response "User does not exist"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/self/info [put]
func (h *user) UpdateProfile(c *gin.Context) {
	userIDvalue, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidTokenError})
		return
	}

	userID, ok := userIDvalue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidTokenError})
		return
	}

	// Bind and validate input
	input, err := utils.BindInputFromRequest[updateProfileInputBody](c)
	if err != nil {
		return
	}

	// Call service to update user
	err = h.svc.UpdateUser(c, userID, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, service.ErrClientNoUpdate):
		c.JSON(http.StatusBadRequest, gin.H{"error": noUpdateDataError})
		return
	case errors.Is(err, dbutils.ErrDuplicationType):
		c.JSON(http.StatusBadRequest, gin.H{"error": emailTakenError})
		return
	case errors.Is(err, nil):
		c.JSON(http.StatusOK, &response.Response{
			Message: updateSelfInfoSuccessMessage,
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}
}
