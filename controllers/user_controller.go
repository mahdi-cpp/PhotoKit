package controllers

import (
	"github.com/mahdi-cpp/PhotoKit/models"
	"github.com/mahdi-cpp/PhotoKit/repositories"
	"github.com/mahdi-cpp/PhotoKit/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userRepo *repositories.UserRepository
}

func NewUserController(userRepo *repositories.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.CreateUserRequest true "Create user"
// @Success 201 {object} models.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user := models.User{
		Username:    req.Username,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Bio:         req.Bio,
		IsOnline:    true,
		LastSeen:    time.Now(),
	}

	if err := uc.userRepo.CreateUser(&user); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, user)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := uc.userRepo.GetUserByID(id)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "Update user"
// @Success 200 {object} models.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateUser := models.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
		IsOnline:  req.IsOnline,
		LastSeen:  time.Now(),
	}

	if err := uc.userRepo.UpdateUser(id, &updateUser); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	user, err := uc.userRepo.GetUserByID(id)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := uc.userRepo.DeleteUser(id); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all users
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.User
// @Failure 500 {object} utils.ErrorResponse
// @Router /users [get]
func (uc *UserController) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := uc.userRepo.ListUsers(limit, offset)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.SendSuccess(c, http.StatusOK, users)
}

// UpdateOnlineStatus godoc
// @Summary Update user online status
// @Description Update a user's online status
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param isOnline query bool true "Online status"
// @Success 200 {object} models.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/{id}/online [put]
func (uc *UserController) UpdateOnlineStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	isOnline, err := strconv.ParseBool(c.Query("isOnline"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid online status")
		return
	}

	if err := uc.userRepo.UpdateUserOnlineStatus(id, isOnline); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to update online status")
		return
	}

	user, err := uc.userRepo.GetUserByID(id)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, user)
}
