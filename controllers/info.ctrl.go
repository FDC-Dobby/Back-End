package controllers

import (
	"github.com/HoseonYim/isfree-backend/database"
	"github.com/HoseonYim/isfree-backend/model"
	"github.com/HoseonYim/isfree-backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CheckUser(c *fiber.Ctx) error {
	jwtClaims := c.Locals("jwtClaims").(utils.AuthTokenClaims)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "OK",
		"user-id": jwtClaims.UserID,
	})
}

func GetInfo(c *fiber.Ctx) error {
	jwtClaims := c.Locals("jwtClaims").(utils.AuthTokenClaims)
	user := model.User{}
	userUUID, err := uuid.Parse(jwtClaims.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(database.UserTable),
		Key: map[string]*dynamodb.AttributeValue{
			"user-id": {
				S: aws.String(userUUID.String()),
			},
		},
	}
	result, err := database.DB.GetItem(input)

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"code":         200,
		"message":      "OK",
		"user-id":      jwtClaims.UserID,
		"email":        user.Email,
		"username":     user.Username,
		"phone-number": user.PhoneNumber,
		"name":         user.Name,
		"birthday":     user.Birthday,
		"created-at":   user.CreatedAt,
	})

}
