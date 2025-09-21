package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TABLE_NAME = "projectsTable"
)

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

type ProjectStore interface {
	DoesProjectExist(string) (bool, error)
	InsertProject()
	EditProject()
	DeleteProject()
	GetProject()
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return DynamoDBClient{
		databaseStore: db,
	}
}

func (d *DynamoDBClient) DoesProjectExist(projectTitle string) (bool, error) {
	result, err := d.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"projectTitle": {
				S: aws.String(projectTitle),
			},
		},
	})

	if err != nil {
		return true, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (d *DynamoDBClient) InsertProject(project types.Project) error {
	//assemble the item (project)
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"projectTitle": {
				S: aws.String(project.Title),
			},
			"projectDescription": {
				S: aws.String(project.Description),
			},
			"githubRepo": {
				S: aws.String(project.Repo),
			},
		},
	}

	_, err := d.databaseStore.PutItem(item)

	if err != nil {
		return err
	}
	return nil
}

func (d *DynamoDBClient) GetProject(projectTitle string) (types.Project, error) {
	var project types.Project
	result, err := d.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"projectTitle": {
				S: aws.String(projectTitle),
			},
		},
	})

	if err != nil {
		return project, err
	}

	if result.Item == nil {
		return project, fmt.Errorf("Project not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &project)
	if err != nil {
		return project, err
	}
	return project, nil
}

func (d *DynamoDBClient) DeleteProject(projectTitle types.Project) error {
	_, err := d.databaseStore.DeleteItem(project, &dynamodb.DeleteItemInput{
		TableName: aws.String(TABLE_NAME),
	})
}
