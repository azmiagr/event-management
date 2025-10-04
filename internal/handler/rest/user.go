package rest

import (
	"event-management/model"
	"event-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) Register(c *gin.Context) {
	param := model.UserRegisterParam{}
	err := c.ShouldBindJSON(&param)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	if param.Password != param.ConfirmPassword {
		response.Error(c, http.StatusBadRequest, "password mismatch", nil)
		return
	}

	err = r.service.UserService.Register(&param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed register user", err)
		return
	}

	response.Success(c, http.StatusCreated, "success register user", nil)
}

func (r *Rest) Login(c *gin.Context) {
	param := model.UserLoginParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid params", err)
		return
	}

	token, err := r.service.UserService.Login(&param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed login user", err)
		return
	}

	response.Success(c, http.StatusOK, "success login user", token)
}
