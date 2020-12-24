package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"iscys.com/shorturl/core"
)

func apiRemove(c *fiber.Ctx) error {
	var keyOrigin, keyID string

	params, err := core.ParseJSONBytes(c.Body())
	if err != nil {
		core.E("invalid request", err)
		return core.FiberJSONError(c, 400, core.ErrInvalidParams)
	}

	keyRaw := core.String(params.GetStringBytes("params", "key"))

	{
		keyOrigin = fmt.Sprintf("u:%s", keyRaw)

		keyID, err = backend.Get(keyOrigin)
		if err == core.ErrObjectNotFound {
		} else if err != nil {
			core.E("get id failed", err)
			return core.FiberJSONError(c, 500, core.ErrServerError)
		} else {
			goto apiRemoveKeys
		}

	}

	{
		keyOrigin, err = backend.Get(fmt.Sprintf("i:%s", keyRaw))
		if err == core.ErrObjectNotFound {
			core.W("key not found:%s", keyRaw)
			return core.FiberJSONError(c, 404, core.ErrObjectNotFound)
		} else if err != nil {
			core.E("get origin url failed", err)
			return core.FiberJSONError(c, 500, core.ErrServerError)
		}

		keyID = keyRaw
	}

apiRemoveKeys:

	backend.Del(keyOrigin)
	backend.Del(fmt.Sprintf("i:%s", keyID))

	core.I("remove: %s => %s", keyID, keyOrigin)
	return core.FiberJSON(c, core.H{"id": keyID})
}
