package main

import (
	"event-management/internal/handler/rest"
	"event-management/pkg/config"
	"event-management/pkg/database/mariadb"
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

	r := rest.NewRest()
	r.MountEndpoint()
	r.Run()
}
