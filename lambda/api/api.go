package api

import "lambda-func/database"

type ApiHandler struct {
	dbstore database.DynamoDBClient
}

func NewApiHandler(dbstore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbstore: dbstore,
	}
}
