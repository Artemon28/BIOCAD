package main

import (
	"BIOCAD/internal"
	"BIOCAD/internal/handlers"
	"BIOCAD/internal/repository"
	"BIOCAD/internal/services"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type flags struct {
	host       string
	port       string
	username   string
	dbname     string
	password   string
	dirAddress string
}

func parseFlags() (fl flags) {
	flag.StringVar(&fl.host, "host", "", "Database host")
	flag.StringVar(&fl.port, "port", "", "Database port")
	flag.StringVar(&fl.username, "username", "", "Database username")
	flag.StringVar(&fl.dbname, "dbname", "", "Database name")
	flag.StringVar(&fl.password, "password", "", "Database password")
	flag.StringVar(&fl.dirAddress, "dirAddress", "", "Address of the directory")
	flag.Parse()
	return
}

func main() {
	fl := parseFlags()

	db, err := internal.NewPostgresDB(internal.Config{
		Host:     fl.host,
		Port:     fl.port,
		Username: fl.username,
		DBName:   fl.dbname,
		SSLMode:  "disable",
		Password: fl.password,
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	router := gin.Default()
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	repo := repository.NewRepository(db)
	service := services.NewService(repo)
	go service.Scan(fl.dirAddress, time.Second*1)
	handler := handlers.NewHandler(service)
	router.GET("/pagination", handler.GetDevices)

	router.Run(":8080")
}
