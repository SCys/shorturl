package main

import (
	"flag"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"
	"iscys.com/shorturl/core"
)

var (
	db *redis.Client

	redisDSN   = flag.String("redis", "redis://localhost:3200", "--redis redis://localhost:3200/0")
	listenAddr = flag.String("listen", "127.0.0.1:3100", "--listen 127.0.0.1:3100")
	hostName   = flag.String("host", "https://mz.ci", "--host https://mz.ci")
)

func main() {
	flag.Parse()

	{
		options, err := redis.ParseURL(*redisDSN)
		if err != nil {
			core.F("invalid redis server options", err)
		}

		db = redis.NewClient(options)

		if sts, err := db.Ping().Result(); err != nil {
			core.F("ping redis server failed", err)
		} else {
			core.I("ping redis:%s", sts)
		}

		core.I("redis server is connected")
		if err != nil {
			core.F("open database failed:%s", err.Error())
		}
	}

	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")
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
		c.Render("index", core.H{"Host": *hostName})
		return nil
	})

	app.Get("/r/:name", apiRedirect)
	app.Post("/r/:url", apiUpsertQuick)
	app.Post("/api/shorturl/upsert", apiUpsert)
	app.Post("/api/shorturl/remove", apiRemove)

	core.I("service is started:%s", *listenAddr)
	_ = app.Listen(*listenAddr)
}
