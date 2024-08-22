package main

import (
	"fmt"
	"os"
	"strconv"
	"tugaskita/app/database"
	"tugaskita/app/migration"
	"tugaskita/app/route"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load(".env")
	port, _ := strconv.Atoi(os.Getenv("SERVERPORT"))
	db := database.Init()
	migration.InitMigration(db)

	e := echo.New()
	e.Use(middleware.CORS())

	route.New(e, db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
