package main

import (
	"event-management/internal/handler/rest"
	"event-management/internal/repository"
	"event-management/internal/service"
	"event-management/pkg/bcrypt"
	"event-management/pkg/config"
	"event-management/pkg/database/mariadb"
	"event-management/pkg/jwt"
	"event-management/pkg/middleware"
	"log"
)

func main() {
	config.LoadEnvironment()

	db, err := mariadb.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = mariadb.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	bcrypt := bcrypt.Init()
	jwt := jwt.Init()
	oauth := config.NewOAuthConfig()
	snapClient := config.NewMidtransSnapClient()
	coreAPIClient := config.NewMidtransCoreAPIClient()
	svc := service.NewService(repo, bcrypt, jwt, oauth, snapClient, coreAPIClient)
	middleware := middleware.Init(svc, jwt)

	r := rest.NewRest(svc, middleware)
	r.MountEndpoint()
	r.Run()
}
