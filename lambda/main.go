package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type GitProjects struct {
	ProjectName string `json:"projectname"`
}

// take in a payload and something with it
func HandleRequest(project GitProjects) (string, error) {
	if project.ProjectName == "" {
		return "", fmt.Errorf("There is no project")
	}

	return fmt.Sprintf("Succesfully called by - %s", project.ProjectName), nil
}

func main() {
	lambda.Start(HandleRequest)
}
