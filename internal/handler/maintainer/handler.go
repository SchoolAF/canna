package maintainer

import (
	"github.com/SchoolAF/exodus/model"
	"github.com/SchoolAF/exodus/repository/db"
	"api/internal/package/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func AddMaintainer(c *fiber.Ctx) error {
	var requestData model.Maintainer
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if err := validator.ValidateData(requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := db.CountMaintainer(requestData.Username); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Data already exists in the database",
		})
	}

	if err := db.InsertMaintainer(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert data into the database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data inserted successfully",
		"data":    requestData,
	})
}

func GetAllMaintainer(c *fiber.Ctx) error {
	filter := bson.M{}

	requestData, err := db.GetMaintainerFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

func GetOneMaintainer(c *fiber.Ctx) error {
	username := c.Params("username")
	filter := bson.M{"username": username}
	requestData, err := db.GetOneMaitainerFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}
