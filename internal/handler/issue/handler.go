package issue

import (
	"api/internal/package/validator"
	"fmt"
	"github.com/SchoolAF/exodus/model"
	"github.com/SchoolAF/exodus/repository/db"
	"github.com/SchoolAF/exodus/repository/profile"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"strconv"
	"time"
)

func GetDeviceName(codename string) (string, error) {
	filter := bson.M{"codename": codename}
	account, err := db.GetOneDeviceFilter(filter)
	if err != nil {
		return "", err
	}
	return account.Marketname, nil
}

func insertionSort(issues []model.IssueData) {
	for i := 1; i < len(issues); i++ {
		key := issues[i]
		j := i - 1

		for j >= 0 && issues[j].Date.Before(key.Date) {
			issues[j+1] = issues[j]
			j = j - 1
		}
		issues[j+1] = key
	}
}

// AddIssue godoc
// @Summary Submit Project Issue
// @Description Submit Project Issue Data
// @Tags Issue
// @Accept json
// @Produce json
// @Param request body IssueDataPOST true "Payload Body [RAW]"
// @Success 200 {object} IssueDataPOST
// @Failure 400
// @Failure 500
// @Router /api/issue [post]
// @Security BearerAuth
func AddIssue(c *fiber.Ctx) error {
	// Parse the request body
	var requestData model.IssueData
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if err := db.CountDevice(requestData.Device); err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Device not found",
		})
	}

	if err := db.CountVersion(requestData.Version); err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Version not found",
		})
	}

	userID, err := profile.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	// Set the extracted user_id to requestData.UserID
	requestData.UserID = userID

	// Generate a random issue ID
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(900000) + 100000
	requestData.IssueID = strconv.Itoa(randomID)
	requestData.Date = time.Now()
	requestData.Edited = false
	requestData.Status = "open"

	if !requestData.Notify {
		requestData.Notify = false
	}

	// Validate the request data
	if err := validator.ValidateData(requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Insert the issue data into the database
	if err := db.InsertIssue(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert data into the database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data inserted successfully",
	})
}

// GetAllIssue godoc
// @Summary Get All Issue Data
// @Description Fetch all submitted issues data by User
// @Tags Issue
// @Accept json
// @Produce json
// @Success 200 {object} IssueData
// @Router /api/issue [get]
func GetAllIssue(c *fiber.Ctx) error {
	filter := bson.M{}

	// Fetch all issues
	issues, err := db.GetIssueFilter(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch issues",
		})
	}

	// Iterate over each issue and fetch the username and device name
	for i := range issues {
		username, err := profile.GetUsernameByID(issues[i].UserID)
		if err != nil {
			issues[i].AuthorName = "Unknown" // or handle the error as needed
			continue
		}
		issues[i].AuthorName = username

		// Fetch device name
		devicename, err := GetDeviceName(issues[i].Device)
		if err != nil {
			issues[i].DeviceParsed = "Unknown" // or handle the case where device name is empty
		} else {
			issues[i].DeviceParsed = fmt.Sprintf("%s (%s)", devicename, issues[i].Device)
		}
	}

	// Apply insertion sort on the issues slice
	insertionSort(issues)

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    issues,
	})
}

func GetUserIssue(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	filter := bson.M{"userid": user_id}

	// Fetch all issues for the user
	issues, err := db.GetIssueFilter(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch issues",
		})
	}

	// Iterate over each issue and fetch the username and device name
	for i := range issues {
		username, err := profile.GetUsernameByID(issues[i].UserID)
		if err != nil {
			issues[i].AuthorName = "Unknown" // or handle the error as needed
			continue
		}
		issues[i].AuthorName = username

		// Fetch device name
		devicename, err := GetDeviceName(issues[i].Device)
		if err != nil {
			issues[i].DeviceParsed = "Unknown" // or handle the case where device name is empty
		} else {
			issues[i].DeviceParsed = fmt.Sprintf("%s (%s)", devicename, issues[i].Device)
		}
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    issues,
	})
}

// GetOneIssue godoc
// @Summary Get Specific Issue Post.
// @Description Fetch one specified issue data subm,itted by User.
// @Tags Issue
// @Accept json
// @Produce json
// @Param id path string true "Insert Post ID"
// @Success 200 {object} IssueData
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/issue/post/{id} [get]
func GetOneIssue(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"issueid": id}
	requestData, err := db.GetOneIssueFilter(filter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to fetch issue data",
		})
	}

	// Validate the request data
	if err := validator.ValidateData(requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Fetch username
	username, err := profile.GetUsernameByID(requestData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch author name",
		})
	}
	requestData.AuthorName = username

	// Fetch device name
	devicename, err := GetDeviceName(requestData.Device)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch device name",
		})
	}
	requestData.DeviceParsed = fmt.Sprintf("%s (%s)", devicename, requestData.Device)

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

// DeleteIssue godoc
// @Summary Delete Issue.
// @Description Delete Submitted Issue Data.
// @Tags Issue
// @Accept json
// @Produce json
// @Param id path string true "Insert Post ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/issue/post/{id} [delete]
// @Security BearerAuth
func DeleteIssue(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"issueid": id}

	var err error

	// Fetch issue data
	issueData, err := db.GetOneIssueFilter(filter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to fetch issue data",
		})
	}

	// Get user ID
	userID, err := profile.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	// Check if user is authorized
	if userID != issueData.UserID {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Delete issue
	_, err = db.DeleteIssue(filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "No document found")
		}
		return err
	}

	// Success response
	return c.JSON(fiber.Map{
		"message": "Issue deleted successfully",
	})
}

// UpdateIssue godoc
// @Summary Update issue
// @Description Update submitted issue data post.
// @Tags Issue
// @Accept json
// @Produce json
// @Param id path string true "Insert Post ID"
// @Param request body IssueDataPOST true "Payload Body [RAW]"
// @Success 200 {object} IssueDataPOST
// @Failure 400
// @Failure 500
// @Router /api/issue/post/{id} [put]
// @Security BearerAuth
func UpdateIssue(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"issueid": id}

	// Fetch issue data
	issueData, err := db.GetOneIssueFilter(filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Issue not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch issue data",
		})
	}

	// Get user ID
	userID, err := profile.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	// Check if user is authorized
	if userID != issueData.UserID {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Parse request body
	var requestData model.IssueData
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Update fields
	requestData.IssueID = issueData.IssueID
	requestData.UserID = userID
	requestData.Date = issueData.Date
	requestData.Status = issueData.Status
	requestData.Edited = true

	// Update issue in database
	if err := db.UpdateIssue(filter, requestData); err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "No document found")
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Issue updated successfully",
	})
}
