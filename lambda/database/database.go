package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const (
	TABLE_NAME = "holmstedProjectTable"
)

type ProjectStore interface {
	DoesProjectExist(string) (bool, error)
	InsertProject(types.Project) error
	GetProject(string) (types.Project, error)
	DeleteProject(string) error
	EditProject(string, types.Project) (types.Project, error)
}

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() *DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return &DynamoDBClient{
		databaseStore: db,
	}
}

func (d *DynamoDBClient) DoesProjectExist(projectID string) (bool, error) {
	result, err := d.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"projectID": {
				S: aws.String(projectID),
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
	if project.ProjectID == "" {
		project.ProjectID = uuid.New().String()
	}

	//assemble the item (project)
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"projectID":   {S: aws.String(project.ProjectID)},
			"title":       {S: aws.String(project.Title)},
			"description": {S: aws.String(project.Description)},
			"repo":        {S: aws.String(project.Repo)},
		},
	}

	_, err := d.databaseStore.PutItem(item)

	if err != nil {
		return err
	}
	return nil
}

func (d *DynamoDBClient) GetProject(projectID string) (types.Project, error) {
	var project types.Project
	result, err := d.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"projectID": {
				S: aws.String(projectID),
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

func (d *DynamoDBClient) DeleteProject(projectID string) error {
	item := &dynamodb.DeleteItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"projectID": {
				S: aws.String(projectID),
			},
		},
	}

	_, err := d.databaseStore.DeleteItem(item)

	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoDBClient) EditProject(projectID string, updatedProject types.Project) (types.Project, error) {
	item := &dynamodb.UpdateItemInput{
		TableName: aws.String(TABLE_NAME), //the table we want to update
		Key: map[string]*dynamodb.AttributeValue{ //Which item I want to update, identified by the "projectID"
			"projectID": {
				S: aws.String(projectID),
			},
		},
		UpdateExpression: aws.String("SET title = :title, description = :desc, repo = :repo"), //which fields that I want updated
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{ //Maps placeholders (:title, :desc, :repo) to their new actual values.
			":title": {S: aws.String(updatedProject.Title)},
			":desc":  {S: aws.String(updatedProject.Description)},
			":repo":  {S: aws.String(updatedProject.Repo)},
		},
		ReturnValues: aws.String("ALL_NEW"), //"ALL_NEW" asks DynamoDB to return the whole updated item.
	}

	result, err := d.databaseStore.UpdateItem(item) //calls the AWS SDK’s method to update a record in DynamoDB.
	if err != nil {
		return types.Project{}, err // if there is an error, we return an empty types struct, and an error
	}

	var updated types.Project
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updated) //converts this map into your Go types.Project struct.
	if err != nil {
		return types.Project{}, err
	}

	return updated, nil

}
