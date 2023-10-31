package database

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var DB dynamodbiface.DynamoDBAPI

const (
	UserTable  = "isfree-user-table"
	TokenTable = "isfree-token-table"
	LocTable   = "isfree-loc-table"
)
