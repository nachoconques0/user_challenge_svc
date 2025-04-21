package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/controller/http/user"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/mocks"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/model"
)

func TestController_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserService(ctrl)
	handler := user.NewController(mockService)

	input := model.CreateUserInput{
		FirstName: "Nacho",
		LastName:  "Calcagno",
		Nickname:  "bandido",
		Password:  "111123123",
		Email:     "nacho@bandidoclub.com",
		Country:   "VE",
	}

	expected := &model.UserOutput{
		ID:        uuid.New().String(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Nickname:  input.Nickname,
		Email:     input.Email,
		Country:   input.Country,
	}

	mockService.EXPECT().
		Create(gomock.Any(), &input).
		Return(expected, nil)

	body, _ := json.Marshal(input)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Create(ctx)

	assert.Equal(t, http.StatusCreated, w.Code)
	var res model.UserOutput
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, input.Nickname, res.Nickname)
}

func TestController_Create_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := user.NewController(mockService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`bad-json`)))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Create(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestController_Find_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := user.NewController(mockService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/users?limit=abc&page=-1", nil)
	ctx.Request = req

	handler.Find(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestController_Update_EmptyNickname(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := user.NewController(mockService)

	id := uuid.New()

	body := []byte(`{"nickname":"   "}`)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodPatch, "/users/"+id.String(), bytes.NewReader(body))
	ctx.Params = gin.Params{{Key: "id", Value: id.String()}}
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Update(ctx)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestController_Update_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := user.NewController(mockService)

	body := []byte(`{"nickname":"newnick"}`)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodPatch, "/users/invalid-uuid", bytes.NewReader(body))
	ctx.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Update(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestController_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := user.NewController(mockService)

	id := uuid.New()
	mockService.EXPECT().
		Delete(gomock.Any(), id).
		Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
	ctx.Params = gin.Params{{Key: "id", Value: id.String()}}

	handler.Delete(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestController_Delete_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := user.NewController(mockService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodDelete, "/users/not-a-uuid", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	handler.Delete(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
