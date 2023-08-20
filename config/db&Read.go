package config

import (
	"encoding/json"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/rijurajarshi/universities-ranking/models"
)

var DB *gorm.DB
var institutions []models.University

func ReadFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(file), &institutions)
}

func PushToDB() {
	db, err := gorm.Open("mysql", "root:GIVE_YOUR_DB_PASSWORD@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE DATABASE IF NOT EXISTS" + " universities")

	db.Exec("USE" + " universities")
	db.AutoMigrate(&models.University{})

	for _, value := range institutions {
		_ = db.Where(models.University{Ranking: value.Ranking}).Assign(models.University{Ranking: value.Ranking, Title: value.Title, Location: value.Location}).FirstOrCreate(&models.University{}).Error
	}
	DB = db
}

func ReadAndLoad() {
	ReadFile("../config/universities_ranking.json")
	PushToDB()
}
