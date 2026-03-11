package main

import (
	"log"
	"os"

	httpapi "github.com/TwoEggDu/content-control-plane/internal/api/http"
	"github.com/TwoEggDu/content-control-plane/internal/application/controlplane"
	"github.com/TwoEggDu/content-control-plane/internal/infrastructure/memory"
	"github.com/gin-gonic/gin"
)

func main() {
	store := memory.NewStore()
	service := controlplane.NewService(store)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	httpapi.RegisterRoutes(router, service)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
