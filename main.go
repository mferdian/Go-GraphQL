package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mferdian/Go-GraphQL/cmd"
	"github.com/mferdian/Go-GraphQL/config/database"
	"github.com/mferdian/Go-GraphQL/config/jwt"
	"github.com/mferdian/Go-GraphQL/domain/user"
	"github.com/mferdian/Go-GraphQL/logging"
	"github.com/mferdian/Go-GraphQL/middleware"
	"github.com/mferdian/Go-GraphQL/routes"
)

func main() {
	// ==== Set up logger ====
	logging.SetUpLogger()
	logging.Log.Info("Logger initialized")

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// DB
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// Seeder command
	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	var (
		jwtService = jwt.NewJWTService()

		userRepo       = user.NewUserRepository(db)
		userService    = user.NewUserService(userRepo, jwtService)
		userController = user.NewUserController(userService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.PublicRoutes(server, userController)
	routes.AdminRoutes(server, userController, jwtService)
	routes.UserRoutes(server, userController, jwtService)

	server.Static("/assets", "./assets")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
