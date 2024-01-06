package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sknutsen/harvestovertimeweb/handler"
	"github.com/sknutsen/harvestovertimeweb/routes"
)

var client *http.Client

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Could not load .env file. Err: %s\n", err)
	}

	client = &http.Client{Timeout: 10 * time.Second}

	h := handler.Handler{
		Client:       client,
		ClientId:     os.Getenv("HARVEST_CLIENT_ID"),
		ClientSecret: os.Getenv("HARVEST_CLIENT_SECRET"),
		Port:         os.Getenv("PORT"),
		Database:     os.Getenv("DATABASE"),
		DbHost:       os.Getenv("DB_HOST"),
		DbUser:       os.Getenv("DB_USER"),
		DbPass:       os.Getenv("DB_PASS"),
		DbPort:       os.Getenv("DB_PORT"),
		Secret:       []byte(os.Getenv("SECRET")),
	}

	e := echo.New()

	e.Use(middleware.Logger())
	// e.Use(session.Middleware(sessions.NewCookieStore(h.Secret)))

	e.Static(routes.Assets, "assets")

	e.GET(routes.Index, h.Index)
	e.GET(routes.Details, h.Details)

	e.GET(routes.Auth, h.Auth)
	e.GET(routes.SigninCallback, h.Callback)
	e.POST(routes.Hours, h.GetOvertimeHours)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", h.Port)))
}
