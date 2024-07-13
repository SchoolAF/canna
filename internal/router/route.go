package router

import (
	"api/config"
	"api/internal/handler/account"
	"api/internal/handler/applyform"
	"api/internal/handler/brand"
	"api/internal/handler/comment"
	"api/internal/handler/device"
	"api/internal/handler/issue"
	"api/internal/handler/maintainer"
	"api/internal/handler/version"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Define JWT middleware configuration
	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: config.JWTSecret,
		ContextKey: "user",
	})

	api.Post("/version", jwtMiddleware, version.AddVersion)
	api.Get("/version", version.GetAllVersion)
	api.Get("/version/:branch", version.GetBranch)

	api.Post("/brand", jwtMiddleware, brand.AddBrand)
	api.Get("/brand", brand.GetAllBrand)
	api.Get("/brand/:brand", brand.GetOneBrand)

	api.Post("/device", jwtMiddleware, device.AddDevice)
	api.Get("/device", device.GetAllDevice)
	api.Get("/device/:codename", device.GetOneDevice)

	api.Post("/maintainer", jwtMiddleware, maintainer.AddMaintainer)
	api.Get("/maintainer", jwtMiddleware, maintainer.GetAllMaintainer)
	api.Get("/maintainer/:username", jwtMiddleware, maintainer.GetOneMaintainer)

	api.Post("/apply/maintainer", jwtMiddleware, applyform.SubmitForm)
	api.Get("/apply/maintainer", jwtMiddleware, applyform.GetAllApplyForm)
	api.Get("/apply/maintainer/:username", jwtMiddleware, applyform.GetOneApplyForm)

	api.Post("/account", account.AddAccount)
	api.Get("/myid", jwtMiddleware, account.GetMyID)
	api.Get("/account", jwtMiddleware, account.GetAllAccount)
	api.Get("/account/:id", jwtMiddleware, account.GetOneAccount)
	api.Post("/auth", account.Authorize)
	api.Post("/passkey/:username", account.RequestPasskey)

	api.Post("/issue", jwtMiddleware, issue.AddIssue)
	api.Get("/issue", issue.GetAllIssue)
	api.Get("/issue/post/:id", issue.GetOneIssue)
	api.Delete("/issue/post/:id", jwtMiddleware, issue.DeleteIssue)
	api.Put("/issue/post/:id", jwtMiddleware, issue.UpdateIssue)
	api.Get("/issue/user/:user_id", issue.GetUserIssue)

	api.Post("/comment/:id", jwtMiddleware, comment.AddComment)
	api.Get("/comment/:id", comment.GetAllComment)

	app.Get("/docs/*", swagger.HandlerDefault)

}
