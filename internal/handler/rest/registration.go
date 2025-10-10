package rest

import (
	"event-management/entity"
	"event-management/model"
	"event-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) RegisterEvent(c *gin.Context) {
	param := model.RegisterEvent{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	user := c.MustGet("user").(*entity.User)
	resp, err := r.service.RegistrationService.RegisterEvent(user.UserID, param.EventID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to register event", err)
		return
	}

	response.Success(c, http.StatusOK, "success register event", resp)

}
