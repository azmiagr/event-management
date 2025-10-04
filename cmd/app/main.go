package main

import (
	"event-management/internal/handler/rest"
	"event-management/internal/repository"
	"event-management/internal/service"
	"event-management/pkg/bcrypt"
	"event-management/pkg/config"
	"event-management/pkg/database/mariadb"
	"event-management/pkg/jwt"
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
	svc := service.NewService(repo, bcrypt, jwt, oauth)

	r := rest.NewRest(svc)
	r.MountEndpoint()
	r.Run()
}
