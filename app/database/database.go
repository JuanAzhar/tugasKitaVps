package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	godotenv.Load(".env")
	host := os.Getenv("DBHOST")
	user := os.Getenv("DBUSER")
	// port := os.Getenv("DBPORT")
	pass := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")

	dbURL := "postgres://"+ user + ":" + pass + "@" + host + "/" + name

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	return db
}
