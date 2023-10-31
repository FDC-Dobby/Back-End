package router

import (
	"github.com/HoseonYim/isfree-backend/controllers"
	"github.com/HoseonYim/isfree-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func loc(router *fiber.App) {
	loc := router.Group("/loc")
	loc.Use(middleware.JWTMiddlware())

	loc.Post("/postLoc", func(c *fiber.Ctx) error {
		return controllers.PostLoc(c)
	})
	loc.Get("/getAllLoc", func(c *fiber.Ctx) error {
		return controllers.GetAllLoc(c)
	})
}
