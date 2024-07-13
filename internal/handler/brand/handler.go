package brand

import (
	"github.com/SchoolAF/exodus/model"
	"github.com/SchoolAF/exodus/repository/db"
	"github.com/SchoolAF/exodus/repository/profile"
	"api/internal/package/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func AddBrand(c *fiber.Ctx) error {
	var requestData model.Brand
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

	if err := db.CountBrand(requestData.BrandLower); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Data already exists in the database",
		})
	} else if err != nil && err.Error() == "data already exists" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to count data in the database",
		})
	}

	if err := db.InsertBrand(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert data into the database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data inserted successfully",
		"data":    requestData,
	})
}

func GetAllBrand(c *fiber.Ctx) error {
	filter := bson.M{}

	requestData, err := db.GetBrandFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

func GetOneBrand(c *fiber.Ctx) error {
	brand := c.Params("brand")
	filter := bson.M{"brand_lower": brand}
	requestData, err := db.GetOneBrandFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}
