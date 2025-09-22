package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.ProjectStore
}

func NewApiHandler(dbStore database.ProjectStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api *ApiHandler) DoesProjectExist(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerProject types.RegisterProject

	err := json.Unmarshal([]byte(request.Body), &registerProject)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerProject.Repo == "" || registerProject.Title == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("register project fields cannot be empty")
	}

	// OPTIONAL: If you want to enforce unique titles/repos, check here.
	// exists, err := api.dbStore.DoesProjectExistByTitleOrRepo(registerProject.Title, registerProject.Repo)
	// if exists {
	//     return events.APIGatewayProxyResponse{ ... }
	// }

	project := types.Project{
		ProjectID:   "",
		Title:       registerProject.Title,
		Repo:        registerProject.Repo,
		Description: registerProject.Description,
	}

	err = api.dbStore.InsertProject(project)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting project into the database %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "Succes",
		StatusCode: http.StatusAccepted,
	}, nil
}
