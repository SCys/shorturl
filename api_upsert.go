package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"iscys.com/shorturl/core"
)

var expireInterval = 7 * 24 * time.Hour // 7d

func apiUpsert(c *fiber.Ctx) error {
	params, err := core.ParseJSONBytes(c.Body())
	if err != nil {
		core.E("invalid request", err)
		return core.FiberJSONError(c, 400, core.ErrInvalidParams)
	}

	urlOrigin := core.String(params.GetStringBytes("params", "url"))
	keyOrigin := fmt.Sprintf("u:%s", urlOrigin)

	urlID, err := db.Get(keyOrigin).Result()
	if redis.Nil == err {
		urlID = xid.New().String()
	} else if err != nil {
		core.E("get origin id failed", err)
		return core.FiberJSONError(c, 500, core.ErrServerError)
	}

	pipe := db.Pipeline()

	pipe.Set(keyOrigin, urlID, expireInterval)
	pipe.Set(fmt.Sprintf("i:%s", urlID), urlOrigin, expireInterval)

	if _, err := pipe.Exec(); err != nil {
		core.E("pipe failed", err)
		return core.FiberJSONError(c, 500, core.ErrServerError)
	}

	core.I("upsert:%s => %s expired %s", urlID, urlOrigin, expireInterval)
	return core.FiberJSON(c, core.H{"id": urlID, "expire": expireInterval.Seconds()})
}
