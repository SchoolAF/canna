package account

import (
	"github.com/SchoolAF/exodus/model"
	"github.com/SchoolAF/exodus/repository/db"
	"github.com/SchoolAF/exodus/repository/profile"
	"api/internal/package/passkey"
	jwtoken "api/internal/package/token"
	"api/internal/package/validator"
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strconv"
)

func AddAccount(c *fiber.Ctx) error {
	var requestData model.Accounts
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	requestData.Passkey = passkey.Generate()
	requestData.Role = "user"

	if err := validator.ValidateData(requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := db.CountTeleID(requestData.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Data already exists in the database",
		})
	} else if err := db.CountUsername(requestData.Username); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Data already exists in the database",
		})
	} else if err != nil && err.Error() == "data already exists" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to count data in the database",
		})
	}

	if err := db.CreateAccount(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert data into the database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data inserted successfully",
	})
}

func GetAllAccount(c *fiber.Ctx) error {
	filter := bson.M{}

	requestData, err := db.GetAccountFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

func GetOneAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"id": id}
	requestData, err := db.GetOneAccountFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

func Authorize(c *fiber.Ctx) error {
	var requestData model.LoginRequest
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	filter := bson.M{"username": requestData.Username}

	accountData, err := db.GetOneAccountFilter(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve account",
		})
	}

	var responseStatus int
	var responseMessage string
	var token string

	if requestData.Username == accountData.Username && requestData.Passkey == accountData.Passkey {
		responseStatus = fiber.StatusOK
		responseMessage = "Login successful"

		token, err = jwtoken.GenerateJWT(accountData.ID, accountData.Role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}
	} else {
		responseStatus = fiber.StatusUnauthorized
		responseMessage = "Invalid credentials"
	}

	// Generate and update the passkey regardless of the credential check
	accountData.Passkey = passkey.Generate()
	if err := passkey.Update(accountData, filter); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := fiber.Map{
		"message": responseMessage,
	}

	if token != "" {
		response["token"] = token
	}

	return c.Status(responseStatus).JSON(response)
}

func RequestPasskey(c *fiber.Ctx) error {
	username := c.Params("username")
	filter := bson.M{"username": username}
	requestData, err := db.GetOneAccountFilter(filter)
	if err != nil {
		return err
	}

	// Convert chat_id from string to integer
	chatID, err := strconv.Atoi(requestData.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid chat_id",
		})
	}

	// Prepare the data for the POST request
	postBody, _ := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    "Your Halcyon Account Passkey is " + strconv.Itoa(requestData.Passkey),
	})

	// Send the POST request
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(config.BotAPI, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
		return err
	}
	defer resp.Body.Close()

	return c.JSON(fiber.Map{
		"message": "Request sent!",
	})
}

func GetMyID(c *fiber.Ctx) error {
	userID, err := profile.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	myRole, err := profile.GetMyRole(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user role",
		})
	}

	isMaintainer := false
	if myRole == "maintainer" || myRole == "admin" {
		isMaintainer = true
	}

	return c.JSON(fiber.Map{
		"userid":       userID,
		"isMaintainer": isMaintainer,
	})
}
