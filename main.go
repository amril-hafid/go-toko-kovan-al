package main

import (
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/database"
	"go-toko-kovan-al/helper"
	"go-toko-kovan-al/router"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

func main() {
	conf := config.Get()
	validate := validator.New()
	store := session.New()
	database.InitDB(conf)

	engine := html.New("./public/componen/", ".html")
	engine.AddFuncMap(helper.FuncMap())
	app := fiber.New(fiber.Config{
		Views:     engine,
		BodyLimit: 50 * 1024 * 1024,
	})

	router.SetupRoutes(app, database.DB, validate, store, conf)

	app.Listen(conf.Srv.Host + ":" + conf.Srv.Port)
}
