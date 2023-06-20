package controller

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/rijurajarshi/universities-ranking/config"
	"github.com/rijurajarshi/universities-ranking/models"
)

var logger *log.Logger

var Local_cache = cache.New(2*time.Minute, 3*time.Minute)

func init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	log.SetOutput(file)
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	logger.Println("Application started.....Running on port 9090")

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		logger.Println("Application closed!!")
		os.Exit(0)
	}()

}

func SetCache(ranking string, univs models.University) bool {

	rankInt, err := strconv.Atoi(ranking)
	if err != nil {
		panic(err)
	}

	if rankInt <= 20 {
		Local_cache.Set(ranking, univs, cache.NoExpiration)
	} else {
		Local_cache.Set(ranking, univs, cache.DefaultExpiration)
	}
	return true
}

func GetCache(ranking string) (interface{}, bool, string) {
	var source string
	data, found := Local_cache.Get(ranking)

	if found {
		source = "Cache"
	}
	return data, found, source
}

func GetAllUniversities(c *gin.Context) {
	var universities []models.University
	cache_key := "universities"
	val, found := Local_cache.Get(cache_key)

	if found {
		logger.Println("Retrieved all universities from the cache successfully")
		c.JSON(200, gin.H{
			"source": "cache", "data from cache": val,
		})
	} else {
		result := config.DB.Find(&universities)
		if result.Error != nil {
			logger.Println("Failed to retrieve universities from the database")
			c.JSON(500, gin.H{
				"Error": "Failed to retrieve universities from the database",
			})
			return
		}
		Local_cache.Set(cache_key, universities, cache.DefaultExpiration)
		logger.Println("Retrieved all universities from the database successfully")
		c.JSON(200, gin.H{
			"source": "database", "data from database": universities,
		})
	}
}

func GetUniversityByRank(c *gin.Context) {
	var university models.University
	var source string

	data, err, source := GetCache(c.Param("ranking"))
	rankint, _ := strconv.Atoi(c.Param("ranking"))

	if !err {
		logger.Println("Cache miss")
		fmt.Println("Cache miss")

		if err := config.DB.Where(&models.University{Ranking: rankint}).First(&university).Error; err != nil {
			logger.Printf("No Record found for rank : %d", rankint)
			fmt.Printf("No Record found for rank : %d", rankint)
			c.JSON(400, gin.H{"Error": "No Record Found"})
			return
		} else {
			SetCache(strconv.Itoa(int(university.Ranking)), university)
			source = "database"
			logger.Printf("Retrieved the university info by rank: %d from the %s successfully", rankint, source)
			fmt.Printf("Retrieved the university info by rank: %d from the %s successfully", rankint, source)
			c.JSON(200, gin.H{
				"data":   university,
				"source": source,
			})
		}
	} else {
		university = data.(models.University)
		logger.Printf("Retrieved the university info by rank: %d from the %s successfully", rankint, source)
		fmt.Printf("Retrieved the university info by rank: %d from the %s successfully", rankint, source)
		c.JSON(200, gin.H{
			"data":   university,
			"source": source,
		})
	}
}

func AddUniversity(c *gin.Context) {
	var university models.University

	err := c.ShouldBindJSON(&university)
	if err != nil {
		logger.Println(err.Error())
		c.JSON(400, gin.H{
			"error": "Invalid request payload",
		})
	}
	err = config.DB.Create(&university).Error
	if err != nil {
		logger.Println(err.Error())
		c.JSON(500, gin.H{
			"error": "Failed to create a university",
		})
	}
	logger.Println("New University created")
	c.JSON(201, gin.H{
		"message": "University created successfully",
	})
}

func UpdateUniversity(c *gin.Context) {

	var university models.University
	var input models.University

	rankint, _ := strconv.Atoi(c.Param("ranking"))

	if err := config.DB.Where(models.University{Ranking: rankint}).First(&university).Error; err != nil {
		logger.Printf("No Record found for rank : %d", rankint)
		fmt.Printf("No Record found for rank : %d", rankint)
		c.JSON(400, gin.H{"Error": "No Record Found"})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Println(err.Error())
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"Error": "Not able to Bind the data",
		})
		return
	}
	result := config.DB.Model(&university).Updates(&input)
	if result.Error != nil {
		logger.Println("Failed to update")
		c.JSON(500, gin.H{
			"Error": "Failed to update",
		})
	}

	logger.Printf("Updated university info for rank : %d successfully", rankint)
	c.JSON(200, gin.H{
		"data": university,
	})
}

func DeleteUniversity(c *gin.Context) {
	var university models.University
	rankint, err := strconv.Atoi(c.Param("ranking"))

	if err != nil {
		panic(err)
	}
	if err := config.DB.Where(&models.University{Ranking: rankint}).First(&university).Error; err != nil {
		logger.Printf("No Record found for rank : %d", rankint)
		fmt.Printf("No Record found for rank : %d", rankint)
		c.JSON(400, gin.H{"Error": "No Record Found"})
		return
	}
	result := config.DB.Model(&university).Delete(&university)
	if result.Error != nil {
		logger.Println("Failed to update")
		c.JSON(500, gin.H{
			"Error": "Failed to update",
		})
	}
	logger.Printf("Deleted university info for rank : %d successfully", rankint)
	c.JSON(200, gin.H{
		"status": "Deleted record successfully",
	})

}
