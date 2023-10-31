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
	"strconv"
	"time"
)

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json := new(LoginRequest)
	if err := c.BodyParser(json); err != nil || json.Email == "" || json.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	// Query DynamoDB for user data using the email
	input := &dynamodb.GetItemInput{
		TableName: aws.String(database.UserTable),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(json.Email),
			},
		},
	}

	result, err := database.DB.GetItem(input)
	if err != nil {
		// Handle error
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	if result.Item == nil {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "User not found",
		})
	}

	// You can use dynamodbattribute.UnmarshalMap to unmarshal DynamoDB item into a struct
	var user = model.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	// Check the password and generate tokens
	if !utils.ComparePasswords(user.Password, json.Password) {
		return c.Status(401).JSON(fiber.Map{
			"code":    401,
			"message": "Invalid Password",
		})
	}

	newUuid, _ := uuid.NewRandom()
	refreshToken := newUuid.String()

	expirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()

	// Create a DynamoDB input item
	input2 := &dynamodb.PutItemInput{
		TableName: aws.String(database.TokenTable),
		Item: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(refreshToken),
			},
			"value": {
				S: aws.String(user.ID.String()),
			},
			"TTL": {
				N: aws.String(strconv.FormatInt(expirationTime, 10)),
			},
		},
	}
	_, err = database.DB.PutItem(input2)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	accessToken, err := utils.CreateToken(user.ID.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"code":          200,
		"message":       "success",
		"access-token":  accessToken,
		"refresh-token": refreshToken,
	})
}

func Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username    string    `json:"username" validate:"required,min=1,max=20,excludesall=;"`
		Email       string    `json:"email" validate:"required,email"`
		PhoneNumber string    `json:"phone-number" validate:"required,e164"` //e164 format is used to save phone number
		Name        string    `json:"name" validate:"required,min=1,max=10,excludesall=;"`
		Password    string    `json:"-" validate:"required,min=8,max=30,excludesall=;"`
		Birthday    time.Time `json:"birthday" validate:"required"`
	}

	json := new(RegisterRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	if err := utils.ValidateStruct(json); err != nil || !utils.VerifyPassword(json.Password) {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	// Check if the user already exists in DynamoDB
	input := &dynamodb.GetItemInput{
		TableName: aws.String(database.UserTable),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(json.Email),
			},
		},
	}
	result, err := database.DB.GetItem(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	if result.Item != nil {
		return c.Status(409).JSON(fiber.Map{
			"code":    409,
			"message": "User already exists",
		})
	}

	// Hash the password and create a new user item
	hashedPassword := utils.HashPassword(json.Password)
	user := model.User{
		ID:          uuid.New(),
		Username:    json.Username,
		Email:       json.Email,
		PhoneNumber: json.PhoneNumber,
		Name:        json.Name,
		Password:    hashedPassword,
		Birthday:    json.Birthday,
		CreatedAt:   time.Now(),
	}

	// Marshal the user struct into a DynamoDB attribute map
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	// Put the item into DynamoDB
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(database.UserTable),
		Item:      av,
	}
	_, err = database.DB.PutItem(putInput)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})
}

func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Query("refresh-token", "")
	_, err := uuid.Parse(refreshToken)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(database.TokenTable), // Replace with your DynamoDB table name
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(refreshToken),
			},
		},
	}

	result, err := database.DB.GetItem(input)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}
	// Check if the item exists
	if len(result.Item) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var token = model.Token{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}

	userId := token.Value
	accessToken, err := utils.CreateToken(userId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}
	return c.JSON(fiber.Map{
		"code":         200,
		"message":      "success",
		"access-token": accessToken,
	})
}

func Logout(c *fiber.Ctx) error {
	refreshToken := c.Query("refresh-token", "")
	_, err := uuid.Parse(refreshToken)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(database.TokenTable), // Replace with your DynamoDB table name
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(refreshToken),
			},
		},
	}
	_, err = database.DB.DeleteItem(input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})
}

//func Verify(c *fiber.Ctx) error {
//	token := c.Query("token", "")
//	userUUID, err := uuid.Parse(token)
//	if err != nil {
//		return c.Status(400).JSON(fiber.Map{
//			"code":    400,
//			"message": "Invalid Query or Token",
//		})
//	}
//
//	unvalidatedUser := UnvalidatedUserStruct{}
//
//	user := UserStruct{
//		Email:    unvalidatedUser.Email,
//		Password: unvalidatedUser.Password,
//		Username: unvalidatedUser.Username,
//	}
//
//	return c.JSON(fiber.Map{
//		"code":    200,
//		"message": "success",
//	})
//}

//func RequestChangePassword(c *fiber.Ctx) error {
//	type RequestChangePasswordRequest struct {
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//
//	json := new(RequestChangePasswordRequest)
//	if err := c.BodyParser(json); err != nil {
//		return c.Status(400).JSON(fiber.Map{
//			"code":    400,
//			"message": "Invalid JSON",
//		})
//	}
//
//	// Check if the user exists in DynamoDB
//	input := &dynamodb.GetItemInput{
//		TableName: aws.String(database.UserTable),
//		Key: map[string]*dynamodb.AttributeValue{
//			"email": {
//				S: aws.String(json.Email),
//			},
//		},
//	}
//	result, err := database.DB.GetItem(input)
//	if err != nil {
//		return c.Status(500).JSON(fiber.Map{
//			"code":    500,
//			"message": "Internal Server Error",
//		})
//	}
//
//	if result.Item == nil {
//		// User does not exist, handle accordingly
//		return c.Status(200).JSON(fiber.Map{
//			"code":    200,
//			"message": "success",
//		})
//	}
//
//	// Generate a change password code and save it to DynamoDB
//	changePasswordCode := uuid.New().String()
//	hashedPassword := utils.HashPassword(json.Password)
//
//	changePasswordObject := YourChangePasswordStruct{
//		ID:       changePasswordCode,
//		UserID:   "user_id", // Replace with your logic to get the user ID
//		Password: hashedPassword,
//	}
//
//	// Marshal the change password object into a DynamoDB attribute map
//	av, err := dynamodbattribute.MarshalMap(changePasswordObject)
//	if err != nil {
//		return c.Status(500).JSON(fiber.Map{
//			"code":    500,
//			"message": "Internal Server Error",
//		})
//	}
//
//	// Put the change password item into DynamoDB
//	putInput := &dynamodb.PutItemInput{
//		TableName: aws.String("YourDynamoDBTableName"),
//		Item:      av,
//	}
//	_, err = database.DB.PutItem(putInput)
//	if err != nil {
//		return c.Status(500).JSON(fiber.Map{
//			"code":    500,
//			"message": "Internal Server Error",
//		})
//	}
//
//	// Send an email with the change password code
//
//	return c.JSON(fiber.Map{
//		"code":    200,
//		"message": "success",
//	})
//}

//func VerifyChangePassword(c *fiber.Ctx) error {
//	token := c.Query("token", "")
//	changePasswordCode, err := uuid.Parse(token)
//	if err != nil {
//		return c.Status(400).JSON(fiber.Map{
//			"code":    400,
//			"message": "Invalid Query or Token",
//		})
//	}
//
//	// Query DynamoDB for the change password code
//	input := &dynamodb.GetItemInput{
//		TableName: aws.String(database.UserTable),
//		Key: map[string]*dynamodb.AttributeValue{
//			"ID": {
//				S: aws.String(changePasswordCode.String()),
//			},
//		},
//	}
//	result, err := database.DB.GetItem(input)
//	if err != nil {
//		return c.Status(500).JSON(fiber.Map{
//			"code":    500,
//			"message": "Internal Server Error",
//		})
//	}
//
//	if result.Item == nil {
//		// Change password code not found, handle accordingly
//		return c.Status(400).JSON(fiber.Map{
//			"code":    400,
//			"message": "Invalid Query or Token",
//		})
//	}
//
//	// Unmarshal the change password item from DynamoDB
//	var changePasswordObject ChangePasswordStruct
//	err = dynamodbattribute.UnmarshalMap(result.Item, &changePasswordObject)
//	if err != nil {
//		return c.Status(500).JSON(fiber.Map{
//			"code":    500,
//			"message": "Internal Server Error",
//		})
//	}
//
//	return c.JSON(fiber.Map{
//		"code":    200,
//		"message": "success",
//	})
//}
