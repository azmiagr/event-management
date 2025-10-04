package rest

import (
	"event-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GoogleLogin(c *gin.Context) {
	url, state, err := r.service.OAuthService.GetGoogleLoginURL()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to generate google login url", err)
		return
	}

	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (r *Rest) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	savedState, err := c.Cookie("oauth_state")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid oauth state", err)
		return
	}

	token, err := r.service.OAuthService.HandleGoogleCallback(code, state, savedState)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to handle google callback", err)
		return
	}

	response.Success(c, http.StatusOK, "google callback success", token)
}
