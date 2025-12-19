package main

import (
    "fmt"
    "log"
    "os"

    "github.com/Lutfania/ekrp/config"
    "github.com/Lutfania/ekrp/database"
    "github.com/Lutfania/ekrp/routes"
)

func main() {

    // load .env
    if err := config.LoadEnv(); err != nil {
        log.Fatal("❌ Failed to load .env:", err)
    }

    // connect PostgreSQL
    if err := config.InitPostgres(); err != nil {
        log.Fatal("❌ Failed to connect PostgreSQL:", err)
    }

    // connect MongoDB (WAJIB!)
    if err := database.InitMongo(); err != nil {
        log.Fatal("❌ Failed to connect MongoDB:", err)
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


