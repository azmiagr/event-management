package rest

import (
	"event-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) HandleNotification(c *gin.Context) {
	var notificationPayload map[string]interface{}
	err := c.ShouldBindJSON(&notificationPayload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid notification payload", err)
		return
	}

	err = r.service.PaymentService.HandleNotification(notificationPayload, r.service.RegistrationService)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to handle payment callback", err)
		return
	}

	response.Success(c, http.StatusOK, "success", nil)
}
