package router

import (
	"github.com/HoseonYim/isfree-backend/controllers"
	"github.com/HoseonYim/isfree-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func info(router *fiber.App) {
	info := router.Group("/info")
	info.Use(middleware.JWTMiddlware())
	info.Get("/checkUser", func(c *fiber.Ctx) error {
		return controllers.CheckUser(c)
	})
	info.Get("/getInfo", func(c *fiber.Ctx) error {
		return controllers.GetInfo(c)
	})
}
