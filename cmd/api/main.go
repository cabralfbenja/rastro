package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"rastro/cmd/api/handlers"
	"rastro/cmd/api/middlewares"
	"rastro/common"
	"rastro/internal/mailer"

	"github.com/labstack/echo/v4"
)

type Application struct {
	Server     *echo.Echo
	handler    handlers.Handler
	Middleware middlewares.AppMiddleware
}

func main() {
	e := echo.New()

	err := godotenv.Load(".env")

	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	db, err := common.NewMySql()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	appMailer := mailer.NewMailer(e.Logger)

	h := handlers.Handler{DB: db, Mailer: appMailer, Logger: e.Logger}

	appMiddleware := middlewares.AppMiddleware{
		DB:     db,
		Logger: e.Logger,
	}

	app := Application{
		Server:     e,
		handler:    h,
		Middleware: appMiddleware,
	}

	e.Use(middleware.Logger())

	app.routes()

	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	e.Logger.Fatal(e.Start(appAddress))
}
