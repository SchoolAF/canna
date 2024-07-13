package comment

import (
	"github.com/SchoolAF/exodus/model"
	"github.com/SchoolAF/exodus/repository/db"
	"github.com/SchoolAF/exodus/repository/profile"
	"api/internal/package/validator"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"strconv"
	"time"
)

func AddComment(c *fiber.Ctx) error {
	id := c.Params("id")

	var requestData model.Comments
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Check if the content is "/delete" before other operations
	if requestData.Content == "/delete" {
		userID, err := profile.GetUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		post_filter := bson.M{"issueid": id}

		issueData, err := db.GetOneIssueFilter(post_filter)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to fetch issue data",
			})
		}

		if userID != issueData.UserID {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		_, err = db.DeleteIssue(post_filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "No document found")
			}
			return err
		}

		_, err = db.DeleteAllComment(post_filter)
		// Handle the error, but do not return it to the client
		if err != nil {
			if err != mongo.ErrNoDocuments {
				fmt.Println("Failed to delete comments:", err)
			}
		}

		return c.JSON(fiber.Map{
			"message": "Issue deleted successfully",
		})
	}

	// Continue with adding a comment if the content is not "/delete"
	userID, err := profile.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	requestData.IssueID = id
	requestData.UserID = userID

	// Generate a random comment ID
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(900000) + 100000
	requestData.CommentID = strconv.Itoa(randomID)
	requestData.Date = time.Now()
	requestData.Status = "open"

	// Validate the request data
	if err := validator.ValidateData(requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Insert the comment data into the database
	if err := db.InsertComment(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert data into the database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data inserted successfully",
		"comment": requestData.Content,
	})
}

func GetAllComment(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"issueid": id}

	// Fetch all issues
	comments, err := db.GetCommentFilter(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch issues",
		})
	}

	// Iterate over each issue and fetch the username and device name
	for i := range comments {
		username, err := profile.GetUsernameByID(comments[i].UserID)
		if err != nil {
			comments[i].AuthorName = "Unknown" // or handle the error as needed
			continue
		}
		comments[i].AuthorName = username
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    comments,
	})
}
