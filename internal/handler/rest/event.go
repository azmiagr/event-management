package rest

import (
	"event-management/entity"
	"event-management/model"
	"event-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateEvent(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	param := model.CreateEventParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	resp, err := r.service.EventService.CreateEvent(user.UserID, &param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create event", err)
		return
	}

	response.Success(c, http.StatusCreated, "success create event", resp)
}
