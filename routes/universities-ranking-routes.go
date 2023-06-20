package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rijurajarshi/universities-ranking/controller"
)

func UniversitiesRankingRoute(router *gin.Engine) {
	router.GET("/all-universities", controller.GetAllUniversities)
	router.GET("/university-by-rank/:ranking", controller.GetUniversityByRank)
	router.POST("/add-university", controller.AddUniversity)
	router.PUT("/update-university/:ranking", controller.UpdateUniversity)
	router.DELETE("/delete-university/:ranking", controller.DeleteUniversity)
}
