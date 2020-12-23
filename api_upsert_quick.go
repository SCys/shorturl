package main

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"iscys.com/shorturl/core"
)

func apiUpsertQuick(c *fiber.Ctx) error {
	urlOrigin := c.Query("url", "")

	if validURL(urlOrigin) {
		core.W("invalid url:%s", urlOrigin)
		c.Status(400)
		return nil
	}

	keyOrigin := fmt.Sprintf("u:%s", urlOrigin)

	urlID, err := db.Get(keyOrigin).Result()
	if redis.Nil == err {
		urlID = xid.New().String()
	} else if err != nil {
		core.E("get origin id failed", err)
		c.Status(500)
		return nil
	}

	pipe := db.Pipeline()

	pipe.Set(keyOrigin, urlID, expireInterval)
	pipe.Set(fmt.Sprintf("i:%s", urlID), urlOrigin, expireInterval)

	if _, err := pipe.Exec(); err != nil {
		core.E("pipe failed", err)
		c.Status(500)
		return nil
	}

	core.I("upsert:%s => %s expired %s", urlID, urlOrigin, expireInterval)
	c.Status(204)
	return nil
}
