package main

import (
	"fmt"

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

	urlID, err := backend.Get(keyOrigin)
	if err == core.ErrObjectNotFound {
		urlID = xid.New().String()
	} else if err != nil {
		core.E("get origin id failed", err)
		c.Status(500)
		return nil
	}

	backend.Set(keyOrigin, urlID, expireInterval)
	backend.Set(fmt.Sprintf("i:%s", urlID), urlOrigin, expireInterval)

	core.I("upsert:%s => %s expired %s", urlID, urlOrigin, expireInterval)
	c.Status(204)
	return nil
}
