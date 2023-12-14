package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/skandansn/webDevBankBackend/middlewares"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/routes"
	"io"
	"os"
	"github.com/joho/godotenv"
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {

	setupLogOutput()

	models.ConnectDataBase()

	server := gin.New()

	_ = godotenv.Load(".env")

	FRONT_END_URL := os.Getenv("FRONT_END_URL")

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{FRONT_END_URL}, // Add your frontend URL here
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	server.Use(gin.Recovery(), middlewares.Logger(), middlewares.AuthCheck())

	group := server.Group("/api")

	for _, route := range routes.Routes {
		group.Handle(route.Method, route.Path, route.Handler)
	}

	server.Run(":8080")
}
