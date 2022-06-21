package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ERR error

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func Init() {
	dbUsername := goDotEnvVariable("DB_USERNAME")
	dbPassword := goDotEnvVariable("DB_PASSWORD")
	dbHost := goDotEnvVariable("DB_HOST")
	dbPort := goDotEnvVariable("DB_PORT")
	dbDatabase := goDotEnvVariable("DB_DATABASE")
	dbOptions := goDotEnvVariable("DB_OPTIONS")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbUsername, dbPassword, dbHost, dbPort, dbDatabase, dbOptions)
	DB, ERR = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if ERR != nil {
		panic(ERR.Error())
	}
}
