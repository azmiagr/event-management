package rest

import (
	"event-management/internal/service"
	"event-management/pkg/middleware"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router     *gin.Engine
	service    *service.Service
	middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
	return &Rest{
		router:     gin.Default(),
		service:    service,
		middleware: middleware,
	}
}

func (r *Rest) MountEndpoint() {
	router := r.router.Group("/api/v1")

	auth := router.Group("/auth")
	auth.POST("/register", r.Register)
	auth.PATCH("/register", r.VerifyUser)
	auth.PATCH("/register/resend", r.ResendOtp)
	auth.POST("/login", r.Login)

	google := auth.Group("/google")
	google.GET("/login", r.GoogleLogin)
	google.GET("/callback", r.GoogleCallback)

	event := router.Group("/events")
	event.Use(r.middleware.AuthenticateUser)
	event.POST("/add-event", r.CreateEvent)

	registration := router.Group("/registrations")
	registration.Use(r.middleware.AuthenticateUser)
	registration.POST("/event", r.RegisterEvent)

	notification := router.Group("/notification")
	notification.POST("", r.HandleNotification)
}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
