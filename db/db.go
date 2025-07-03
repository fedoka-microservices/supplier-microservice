package db

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func InitDB(user string, pass string, host string, port string, dbName string) *gorm.DB {
	godotenv.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, dbName)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}
		log.Printf("⏳ waiting for connection (%d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return nil
}
