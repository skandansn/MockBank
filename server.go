package main

import (
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/middlewares"
	"github.com/skandansn/webDevBankBackend/models"
	"github.com/skandansn/webDevBankBackend/routes"
	"io"
	"os"
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {

	setupLogOutput()

	models.ConnectDataBase()

	server := gin.New()
	server.Use(gin.Recovery(), middlewares.Logger(), middlewares.AuthCheck())

	group := server.Group("/api")

	for _, route := range routes.Routes {
		group.Handle(route.Method, route.Path, route.Handler)
	}

	server.Run(":8080")
}
