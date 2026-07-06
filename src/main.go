package main

import (
	"flag"
	"log"
	"os"
	"shift-be/src/config"
	"shift-be/src/features/user"
	"shift-be/src/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file")
	}

	seedFlag := flag.Bool("seed", false, "Run database seeders")
	flag.Parse()

	db, psql := config.ConnectDB()
	defer db.Close()

	if *seedFlag {
		log.Println("Seeding database...")
		user.SeedUsers(psql)
		log.Println("Seeding completed!")
		return
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	appPath := os.Getenv("APP_PATH")
	router.SetupRouter(r, appPath, psql)

	port := os.Getenv("PORT")
	log.Printf("Server running at port %s (Hot Reload Active)...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
