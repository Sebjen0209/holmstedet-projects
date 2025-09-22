package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HolmstedetProjectsStackProps struct {
	awscdk.StackProps
}

func NewHolmstedetProjectsStack(scope constructs.Construct, id string, props *HolmstedetProjectsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// define the dynamodb table
	table := awsdynamodb.NewTable(stack, jsii.String("projectTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("projectID"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("holmstedProjectTable"),
	})

	// define the lambda here
	lambdaFunction := awslambda.NewFunction(stack, jsii.String("projectLambdaFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
	})

	// define the apigateway
	apiGateway := awsapigateway.NewRestApi(stack, jsii.String("projectGateWay"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("POST", "GET", "PUT", "DELETE", "OPTIONS", "UPDATE"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
		EndpointConfiguration: &awsapigateway.EndpointConfiguration{
			Types: &[]awsapigateway.EndpointType{awsapigateway.EndpointType_REGIONAL},
		},
	})

	// integrate the routes + methods to this endpoint
	// Lambda Integration
	lambdaIntegration := awsapigateway.NewLambdaIntegration(lambdaFunction, nil)

	// /project resource
	projectResource := apiGateway.Root().AddResource(jsii.String("project"), nil)

	// POST /project -- create new project
	projectResource.AddMethod(jsii.String("POST"), lambdaIntegration, nil)

	// /project/{projectID} resource for GET, PUT, DELETE
	projectIdResource := projectResource.AddResource(jsii.String("{projectID}"), nil)

	// GET /project/{projectID} -- get project by ID
	projectIdResource.AddMethod(jsii.String("GET"), lambdaIntegration, nil)

	// PUT /project/{projectID} -- update project by ID
	projectIdResource.AddMethod(jsii.String("PUT"), lambdaIntegration, nil)

	// DELETE /project/{projectID} -- delete project by ID
	projectIdResource.AddMethod(jsii.String("DELETE"), lambdaIntegration, nil)

	table.GrantReadWriteData(lambdaFunction)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewHolmstedetProjectsStack(app, "HolmstedetProjectsStack", &HolmstedetProjectsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
