package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	godotenv.Load(".env")
	host := os.Getenv("DBHOST")
	user := os.Getenv("DBUSER")
	port := os.Getenv("DBPORT")
	pass := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")

	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	return db
}
