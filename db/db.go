package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func InitDB() *gorm.DB {
	godotenv.Load()
	user := getEnv("DB_USERNAME", "")
	pass := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "")
	port := getEnv("DB_PORT", "")
	dbName := getEnv("DB_NAME", "")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	return db
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
