package main

import (
	"flag"
	"log"
	"os"
	"tms-server/config"
	"tms-server/migrations"
	"tms-server/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
}

func main() {
	migrate := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	config.ConnectDB()

	if *migrate {
		if err := migrations.Migrate(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations completed. Exiting.")
		return
	}

	r := gin.Default()
	routes.RegisterRoutes(r)

	port := os.Getenv("APP_PORT")
	r.Run(":" + port)
}
