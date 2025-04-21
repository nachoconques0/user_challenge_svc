package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/model"
	service "github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/service/user"
)

type Controller struct {
	svc service.Service
}

// NewController returns a HTTP User controller
func NewController(s service.Service) *Controller {
	return &Controller{svc: s}
}

// Create returns a created user
func (c *Controller) Create(ctx *gin.Context) {
	var input model.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Str("userController", "Create").Msg("not valid data")
		returnsWithError(ctx, http.StatusBadRequest, "invalid input", err.Error())
		return
	}

	user, err := c.svc.Create(ctx, &input)
	if err != nil {
		log.Error().Err(err).Str("userController", "Create").Msg("could not create user")
		returnsWithError(ctx, http.StatusInternalServerError, "could not create user", err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// Find returns a list of users. It is paginated and also can be filtered by country
func (c *Controller) Find(ctx *gin.Context) {
	country := ctx.DefaultQuery("country", "")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		log.Error().Err(err).Str("userController", "Find").Msg("invalid pagination page param")
		returnsWithError(ctx, http.StatusBadRequest, "invalid page parameter")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		log.Error().Err(err).Str("userController", "Find").Msg("invalid pagination limit param")
		returnsWithError(ctx, http.StatusBadRequest, "invalid limit parameter")
		return
	}

	users, err := c.svc.Find(ctx, country, page, limit)
	if err != nil {
		log.Error().Err(err).Str("userController", "Find").Msg("could not find users")
		returnsWithError(ctx, http.StatusInternalServerError, "could not find users", err.Error())
		return
	}

	returnsWithSuccess(ctx, users)
}

// Update updates an user nickname
func (c *Controller) Update(ctx *gin.Context) {
	var input struct {
		Nickname string `json:"nickname" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Str("userController", "Update").Msg("invalid update data")
		returnsWithError(ctx, http.StatusBadRequest, "invalid update data", err.Error())
		return
	}

	if strings.TrimSpace(input.Nickname) == "" {
		log.Error().Str("userController", "Update").Msg("nickname cannot be empty or whitespace")
		returnsWithError(ctx, http.StatusUnprocessableEntity, "nickname cannot be empty or whitespace")
		return
	}

	paramID := ctx.Param("id")
	id, err := uuid.Parse(paramID)
	if err != nil {
		log.Error().Err(err).Str("userController", "Update").Msg("invalid ID formatt")
		returnsWithError(ctx, http.StatusBadRequest, "invalid ID format", err.Error())
		return
	}

	updatedUser, err := c.svc.Update(ctx, id, input.Nickname)
	if err != nil {
		log.Error().Err(err).Str("userController", "Update").Msg("could not update user")
		returnsWithError(ctx, http.StatusInternalServerError, "could not update user", err.Error())
		return
	}

	returnsWithSuccess(ctx, updatedUser)
}

// Delete soft deletes an user based on its ID
func (c *Controller) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("userController", "Delete").Msg("invalid user ID")
		returnsWithError(ctx, http.StatusBadRequest, "invalid user ID", err.Error())
		return
	}

	err = c.svc.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("userController", "Delete").Msg("could not delete user")
		returnsWithError(ctx, http.StatusInternalServerError, "could not delete user", err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func returnsWithError(ctx *gin.Context, code int, message string, details ...string) {
	res := model.ErrorResponse{Error: message}
	if len(details) > 0 {
		res.Details = details[0]
	}
	ctx.JSON(code, res)
}

func returnsWithSuccess(ctx *gin.Context, payload interface{}) {
	ctx.JSON(http.StatusOK, payload)
}
