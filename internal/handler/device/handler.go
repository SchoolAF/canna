package device

import (
	"github.com/SchoolAF/exodus/model"
	"github.com/SchoolAF/exodus/repository/db"
	"github.com/SchoolAF/exodus/repository/profile"
	"api/internal/package/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func AddDevice(c *fiber.Ctx) error {
	var requestData model.Device
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
	_, err := profile.AuthorizeAdmin(c)
	if err != nil {
		return err
	}

	requestData.Deprecated = false

	if err := db.CountDevice(requestData.Codename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Data already exists in the database",
		})
	}

	if db.CountBrand(requestData.BrandLower) == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Brand not found!",
		})
	}

	if db.CountMaintainer(requestData.Maintainer) == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Maintainer not found!",
		})
	}

	if err := db.InsertDevice(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert data into the database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data inserted successfully",
		"data":    requestData,
	})
}

func GetAllDevice(c *fiber.Ctx) error {
	filter := bson.M{}

	requestData, err := db.GetDeviceFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

func GetOneDevice(c *fiber.Ctx) error {
	codename := c.Params("codename")
	filter := bson.M{"codename": codename}
	requestData, err := db.GetOneDeviceFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}
