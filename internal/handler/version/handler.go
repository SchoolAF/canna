package version

import (
	"github.com/SchoolAF/exodus/model"
	"github.com/SchoolAF/exodus/repository/db"
	"github.com/SchoolAF/exodus/repository/profile"
	"api/internal/package/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func AddVersion(c *fiber.Ctx) error {
	var requestData model.Version
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	_, err := profile.AuthorizeAdmin(c)
	if err != nil {
		return err
	}

	if err := validator.ValidateData(requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := db.CountVersion(requestData.Branch); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Data already exists in the database",
		})
	} else if err != nil && err.Error() == "data already exists" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to count data in the database",
		})
	}

	if err := db.InsertVersion(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert data into the database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data inserted successfully",
		"data":    requestData,
	})
}

func GetAllVersion(c *fiber.Ctx) error {
	filter := bson.M{}

	requestData, err := db.GetVersionFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

func GetBranch(c *fiber.Ctx) error {
	branch := c.Params("branch")
	filter := bson.M{"branch": branch}
	requestData, err := db.GetOneVersionFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}
