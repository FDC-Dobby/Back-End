package router

import (
	"github.com/HoseonYim/isfree-backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func auth(router *fiber.App) {
	auth := router.Group("/auth")

	auth.Post("/login", func(c *fiber.Ctx) error {
		return controllers.Login(c)
	})
	auth.Post("/register", func(c *fiber.Ctx) error {
		return controllers.Register(c)
	})
	auth.Get("/refresh", func(c *fiber.Ctx) error {
		return controllers.Refresh(c)
	})
	auth.Get("/logout", func(c *fiber.Ctx) error {
		return controllers.Logout(c)
	})
	//auth.Get("/verify", func(c *fiber.Ctx) error {
	//	return controllers.Verify(c)
	//})
	//auth.Post("/changePass", func(c *fiber.Ctx) error {
	//	return controllers.RequestChangePassword(c)
	//})
	//auth.Get("/changePass", func(c *fiber.Ctx) error {
	//	return controllers.VerifyChangePassword(c)
	//})
}
