package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"
	"iscys.com/shorturl/core"
)

var backend Backend

var (
	dsnRaw      = flag.String("dsn", "./data/main.db", "--dsn /data/shorturl/data/main.db")
	listenAddr  = flag.String("listen", "127.0.0.1:3100", "--listen 127.0.0.1:3100")
	hostName    = flag.String("host", "https://mz.ci", "--host https://mz.ci")
	backendType = flag.String("backend", "badger", "--backend badger")
)

func main() {
	flag.Parse()

	switch *backendType {
	case "redis":
		backend = &BackendRedis{}
	case "badger":
		backend = &BackendBadger{}
	default:
		core.F("invalid backend type:%s", *backendType)
	}

	backend.Init()

	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

	// DEBUG code
	// engine := html.NewFileSystem(http.Dir("./views"), ".html")
	// engine.Reload(true)
	// engine.Debug(true)

	app := fiber.New(fiber.Config{Views: engine})
	app.Use(
		core.FiberBasicInfo(),
		recover.New(),
		cors.New(cors.Config{AllowOrigins: *hostName}),
	)

	app.Get("", func(c *fiber.Ctx) error {
		c.Render("index", fiber.Map{"Host": *hostName})
		return nil
	})

	app.Get("/r/:name", apiRedirect)
	app.Post("/r/:url", apiUpsertQuick)
	app.Post("/api/shorturl/upsert", apiUpsert)
	app.Post("/api/shorturl/remove", apiRemove)

	core.I("service is started:%s", *listenAddr)
	_ = app.Listen(*listenAddr)
}
