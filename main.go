package main

import (
    "fmt"
    "log"
    "os"

    "ekrp/config"
    "ekrp/routes"
)

func main() {

	err := config.LoadEnv()
	if err != nil {
		log.Fatal("❌ Failed to load .env:", err)
	}

	err = config.InitPostgres()
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	app := config.NewApp()

	routes.RegisterRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running at http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
}
