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
	"time"
)

//func GetLoc(c *fiber.Ctx) error {
//	type LocRequest struct {
//		ID uuid.UUID `json:"loc-id"`d
//		Name string    `json:"name" validate:"required,min=1,max=20,excludesall=;"`
//	}
//	jwtClaims := c.Locals("jwtClaims").(utils.AuthTokenClaims)
//	user := model.User{}
//	userUUID, err := uuid.Parse(jwtClaims.UserID)
//
//}

func GetAllLoc(c *fiber.Ctx) error {
	input := &dynamodb.ScanInput{
		TableName: aws.String(database.LocTable),
	}

	result, err := database.DB.Scan(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	var locations []model.Loc
	for _, item := range result.Items {
		loc := model.Loc{}
		if err := dynamodbattribute.UnmarshalMap(item, &loc); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"code":    500,
				"message": "Internal Server Error",
			})
		}
		locations = append(locations, loc)
	}

	return c.JSON(fiber.Map{
		"code":      200,
		"message":   "success",
		"locations": locations,
	})
}

func PostLoc(c *fiber.Ctx) error {
	type LocRequest struct {
		//ID        uuid.UUID `json:"loc-id"`
		Name      string  `json:"name" validate:"required,min=1,max=20,excludesall=;"`
		Latitude  float64 `json:"latitude" validate:"required"`
		Longitude float64 `json:"longitude" validate:"required"`
		Runway    bool    `json:"runway"`
		Elevator  bool    `json:"elevator"`
		Parking   bool    `json:"parking"`
		Restroom  bool    `json:"restroom"`
		Block     bool    `json:"block"`
		Guide     bool    `json:"guide"`
		Review    string  `json:"review"`
	}
	//jwtClaims := c.Locals("jwtClaims").(utils.AuthTokenClaims)
	//user := model.User{}
	//userUUID, err := uuid.Parse(jwtClaims.UserID)

	json := new(LocRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	if err := utils.ValidateStruct(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(database.LocTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(json.Name),
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
	loc := model.Loc{}
	if result.Item != nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &loc)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"code":    500,
				"message": "Internal Server Error",
			})
		}
		if json.Runway {
			loc.Runway++
		}
		if json.Elevator {
			loc.Elevator++
		}
		if json.Parking {
			loc.Parking++
		}
		if json.Restroom {
			loc.Restroom++
		}
		if json.Block {
			loc.Block++
		}
		if json.Guide {
			loc.Guide++
		}
		if len(json.Review) != 0 {
			loc.Review = append(loc.Review, json.Review)
		}

		loc = model.Loc{
			ID:        uuid.New(),
			Name:      json.Name,
			Latitude:  json.Latitude,
			Longitude: json.Longitude,
			Runway: func() int64 {
				if json.Runway {
					return 1
				}
				return 0
			}(),
			Elevator: func() int64 {
				if json.Elevator {
					return 1
				}
				return 0
			}(),
			Parking: func() int64 {
				if json.Parking {
					return 1
				}
				return 0
			}(),
			Restroom: func() int64 {
				if json.Restroom {
					return 1
				}
				return 0
			}(),
			Block: func() int64 {
				if json.Block {
					return 1
				}
				return 0
			}(),
			Guide: func() int64 {
				if json.Guide {
					return 1
				}
				return 0
			}(),
			Review:    []string{json.Review},
			CreatedAt: time.Now(),
		}
		return c.Status(409).JSON(fiber.Map{
			"code":    409,
			"message": "User already exists",
		})
	} else {
		//user doesn't exist make new one

		loc = model.Loc{
			ID:        uuid.New(),
			Name:      json.Name,
			Latitude:  json.Latitude,
			Longitude: json.Longitude,
			Runway: func() int64 {
				if json.Runway {
					return 1
				}
				return 0
			}(),
			Elevator: func() int64 {
				if json.Elevator {
					return 1
				}
				return 0
			}(),
			Parking: func() int64 {
				if json.Parking {
					return 1
				}
				return 0
			}(),
			Restroom: func() int64 {
				if json.Restroom {
					return 1
				}
				return 0
			}(),
			Block: func() int64 {
				if json.Block {
					return 1
				}
				return 0
			}(),
			Guide: func() int64 {
				if json.Guide {
					return 1
				}
				return 0
			}(),
			Review:    []string{json.Review},
			CreatedAt: time.Now(),
		}
	}

	av, err := dynamodbattribute.MarshalMap(loc)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(database.LocTable),
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
