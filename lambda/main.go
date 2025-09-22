package main

import (
	"fmt"
	"lambda-func/app"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type GitProjects struct {
	ProjectID string `json:"projectID"`
}

// take in a payload and something with it
func HandleRequest(project GitProjects) (string, error) {
	if project.ProjectID == "" {
		return "", fmt.Errorf("There is no project")
	}

	return fmt.Sprintf("Succesfully called by - %s", project.ProjectID), nil
}

func main() {
	myAPP := app.NewApp()
	handler := myAPP.ApiHandler

	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Resource {
		case "/project":
			switch request.HTTPMethod {
			case "POST":
				return handler.CreateProject(request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusMethodNotAllowed,
					Body:       "Method Not Allowed",
				}, nil
			}
		case "/project/{projectID}":
			switch request.HTTPMethod {
			case "GET":
				return handler.GetProject(request)
			case "PUT":
				return handler.EditProject(request)
			case "DELETE":
				return handler.DeleteProject(request)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusMethodNotAllowed,
					Body:       "Method Not Allowed",
				}, nil
			}
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       "Not Found",
			}, nil
		}
	})
}
