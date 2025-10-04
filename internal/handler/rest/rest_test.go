package rest

import (
	"bytes"
	"encoding/json"
	"event-management/internal/service"
	"event-management/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock of IUserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(param *model.UserRegisterParam) error {
	args := m.Called(param)
	return args.Error(0)
}

func (m *MockUserService) Login(param *model.UserLoginParam) (*model.UserLoginResponse, error) {
	args := m.Called(param)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserLoginResponse), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)
	svc := &service.Service{
		UserService: mockUserService,
	}
	restHandler := NewRest(svc)

	router := gin.Default()
	router.POST("/api/v1/auth/register", restHandler.Register)

	registerParam := model.UserRegisterParam{
		Name:            "Test User",
		Email:           "test@example.com",
		Password:        "password",
		ConfirmPassword: "password",
	}

	mockUserService.On("Register", &registerParam).Return(nil)

	body, _ := json.Marshal(registerParam)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)
	svc := &service.Service{
		UserService: mockUserService,
	}
	restHandler := NewRest(svc)

	router := gin.Default()
	router.POST("/api/v1/auth/login", restHandler.Login)

	loginParam := model.UserLoginParam{
		Email:    "test@example.com",
		Password: "password",
	}

	loginResponse := &model.UserLoginResponse{
		Token: "some.jwt.token",
	}

	mockUserService.On("Login", &loginParam).Return(loginResponse, nil)

	body, _ := json.Marshal(loginParam)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NotNil(t, responseBody["data"])
	data := responseBody["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
}