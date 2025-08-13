package rest

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router *gin.Engine
}

func NewRest() *Rest {
	return &Rest{
		router: gin.Default(),
	}
}

func (r *Rest) MountEndpoint() {

}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
