package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
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

func (api *ApiHandler) CreateProject(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input types.RegisterProject

	err := json.Unmarshal([]byte(request.Body), &input)

	// tjekker om er der andre fejl
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request" + err.Error(),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	//tjeker om felterne er tomme
	if input.Description == "" || input.Repo == "" || input.Title == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request" + err.Error(),
			StatusCode: http.StatusBadGateway,
		}, err
	}

	//når vi har tjekket om felterne er udfyldte, så laver vi et project med felterne
	project := types.Project{
		ProjectID:   uuid.New().String(),
		Title:       input.Title,
		Description: input.Description,
		Repo:        input.Repo,
	}

	//tjekker om der er fejl på database laget
	err = api.dbStore.InsertProject(project)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error creating project",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	body, _ := json.Marshal(project)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusAccepted,
	}, nil
}

func (api *ApiHandler) GetProject(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projectID := request.PathParameters["projectID"]

	//tjekker om der er et projekt med det angivede projektID
	if projectID == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Missing ProjectID",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	//Tjekker om projektet er i databasen
	project, err := api.dbStore.GetProject(projectID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Project not found",
			StatusCode: http.StatusNotFound,
		}, err
	}

	body, _ := json.Marshal(project)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

func (api *ApiHandler) EditProject(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projectID := request.PathParameters["projectID"]

	if projectID == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Missing ProjectID",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	// Takes the request and converts it to a Go struct using Unmarshal so we can work with it
	var input types.RegisterProject
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request" + err.Error(),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	//tjeker om felterne er tomme
	if input.Description == "" || input.Repo == "" || input.Title == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request" + err.Error(),
			StatusCode: http.StatusBadGateway,
		}, err
	}

	updatedProject := types.Project{
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		Repo:        input.Repo,
	}

	// Checks for any errors returned from the database when updating the project
	project, err := api.dbStore.EditProject(projectID, updatedProject)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error updating the project",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	body, _ := json.Marshal(project)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

func (api *ApiHandler) DeleteProject(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projectID := request.PathParameters["projecktID"]

	if projectID == "" {
		return events.APIGatewayProxyResponse{
			Body:       "There is no project with such ID",
			StatusCode: http.StatusBadGateway,
		}, nil
	}

	err := api.dbStore.DeleteProject(projectID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error deleting the project",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       "Project got deleted",
		StatusCode: http.StatusOK,
	}, nil
}
