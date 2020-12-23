package main

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"iscys.com/shorturl/core"
)

func apiRedirect(c *fiber.Ctx) error {
	keyID := c.Params("name")

	urlOrigin, err := db.Get(fmt.Sprintf("i:%s", keyID)).Result()
	if redis.Nil == err {
		core.W("id is not found")
		return core.FiberJSONError(c, 404, core.ErrObjectNotFound)
	} else if err != nil {
		core.E("get origin url failed", err)
		return core.FiberJSONError(c, 500, core.ErrServerError)
	}

	core.I("redirect: %s => %s", keyID, urlOrigin)
	return c.Redirect(urlOrigin, 302)
}
