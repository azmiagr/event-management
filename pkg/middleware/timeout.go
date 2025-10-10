package middleware

import (
	"errors"
	"event-management/pkg/response"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func (m *middleware) Timeout() gin.HandlerFunc {
	timeLimit, _ := strconv.Atoi(os.Getenv("TIME_OUT_LIMIT"))

	return timeout.New(
		timeout.WithTimeout(time.Duration(timeLimit)*time.Second),
		timeout.WithResponse(timeoutResponse),
	)
}

func timeoutResponse(c *gin.Context) {
	response.Error(c, http.StatusRequestTimeout, "the request take to much time", errors.New(""))
}
