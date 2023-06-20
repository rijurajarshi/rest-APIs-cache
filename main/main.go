package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rijurajarshi/universities-ranking/config"
	"github.com/rijurajarshi/universities-ranking/routes"
)

func main() {
	router := gin.Default()
	config.ReadAndLoad()
	routes.UniversitiesRankingRoute(router)
	router.Run(":9090")
}
