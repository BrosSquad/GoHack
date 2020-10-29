package main

import (
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/viper"
	"gohack/api"
	"gohack/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Setting up Viper
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/gohack")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	app := fiber.New(fiber.Config{
		Prefork: viper.GetBool("http.prefork"),
	})

	db, err := database.ConnectDB(database.Config{
		Host:     viper.GetString("database.host"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DbName:   viper.GetString("database.dbname"),
		Port:     uint16(viper.GetUint32("database.port")),
		TimeZone: viper.GetString("database.timezone"),
		SslMode:  viper.GetBool("database.sslmode"),
	})

	_ = api.Container{
		DB: db,
	}

	if err != nil {
		log.Fatalf("Fatal error database file: %s \n", err)
	}

	setupMiddleware(app)

	// Serving React
	app.Static("/", "../../dist")
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("../../dist/index.html")
	})

	if err := app.Listen(viper.GetString("http.address")); err != nil {
		log.Fatalf("Error while starting application: %v", err)
	}
}

func setupMiddleware(app *fiber.App) {
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(csrf.New(csrf.Config{
		Cookie: &fiber.Cookie{
			HTTPOnly: true,
		},
	}))
}
