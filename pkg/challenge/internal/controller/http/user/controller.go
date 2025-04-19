package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/model"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/service"
)

type Controller struct {
	svc service.Service
}

func NewController(s service.Service) *Controller {
	return &Controller{svc: s}
}

// Create returns a created user
func (c *Controller) Create(ctx *gin.Context) {
	var input model.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		returnsWithError(ctx, http.StatusBadRequest, "invalid input", err.Error())
		return
	}

	user, err := c.svc.Create(ctx, &input)
	if err != nil {
		returnsWithError(ctx, http.StatusInternalServerError, "could not create user", err.Error())
		return
	}

	returnsWithSuccess(ctx, user)
}

// Find returns a list of users. It is paginated and also can be filtered by country
func (c *Controller) Find(ctx *gin.Context) {
	country := ctx.DefaultQuery("country", "")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	users, err := c.svc.Find(ctx, country, page, limit)
	if err != nil {

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
		returnsWithError(ctx, http.StatusBadRequest, "invalid update data", err.Error())
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		returnsWithError(ctx, http.StatusBadRequest, "invlid ID", err.Error())
		return
	}

	updatedUser, err := c.svc.Update(ctx, id, input.Nickname)
	if err != nil {
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
		returnsWithError(ctx, http.StatusBadRequest, "invalid user ID", err.Error())
		return
	}

	err = c.svc.Delete(ctx, id)
	if err != nil {
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
