// Package dynamodbquery executes a query against an Amazon DynamoDB
package dynamodbquery

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// Constants used by the code to represent the input and outputs of the JSON structure
const (
	ivAWSAccessKeyID                 = "AWSAccessKeyID"
	ivAWSSecretAccessKey             = "AWSSecretAccessKey"
	ivAWSDefaultRegion               = "AWSDefaultRegion"
	ivDynamoDBTableName              = "DynamoDBTableName"
	ivDynamoDBKeyConditionExpression = "DynamoDBKeyConditionExpression"
	ivDynamoDBExpressionAttributes   = "DynamoDBExpressionAttributes"
	ivDynamoDBFilterExpression       = "DynamoDBFilterExpression"

	ovResult           = "result"
	ovScannedCount     = "scannedCount"
	ovConsumedCapacity = "consumedCapacity"
)

// log is the default package logger
var log = logger.GetLogger("activity-dynamodbquery")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// ExpressionAttribute is a structure representing the JSON payload for the expression syntax
type ExpressionAttribute struct {
	Name  string
	Value string
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// Get the inputs
	awsAccessKeyID := context.GetInput(ivAWSAccessKeyID).(string)
	awsSecretAccessKey := context.GetInput(ivAWSSecretAccessKey).(string)
	awsDefaultRegion := context.GetInput(ivAWSDefaultRegion).(string)
	dynamoDBTableName := context.GetInput(ivDynamoDBTableName).(string)
	dynamoDBKeyConditionExpression := context.GetInput(ivDynamoDBKeyConditionExpression).(string)
	dynamoDBExpressionAttributes := context.GetInput(ivDynamoDBExpressionAttributes)
	dynamoDBFilterExpression := context.GetInput(ivDynamoDBFilterExpression).(string)

	// Create new credentials using the accessKey and secretKey
	awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	// Create a new session to AWS
	awsSession := session.Must(session.NewSession(&aws.Config{
		Credentials: awsCredentials,
		Region:      aws.String(awsDefaultRegion),
	}))

	// Create a new login to the DynamoDB service
	dynamoService := dynamodb.New(awsSession)

	// Construct the expression attributes from the JSON payload
	var expressionAttributes []ExpressionAttribute
	json.Unmarshal([]byte(dynamoDBExpressionAttributes.(string)), &expressionAttributes)

	expressionAttributeMap := make(map[string]*dynamodb.AttributeValue)
	for _, attribute := range expressionAttributes {
		expressionAttributeMap[attribute.Name] = &dynamodb.AttributeValue{S: aws.String(attribute.Value)}
	}

	// Construct the DynamoDB query
	var queryInput = &dynamodb.QueryInput{}
	if dynamoDBFilterExpression == "" {
		queryInput = &dynamodb.QueryInput{
			TableName:                 aws.String(dynamoDBTableName),
			KeyConditionExpression:    aws.String(dynamoDBKeyConditionExpression),
			ExpressionAttributeValues: expressionAttributeMap,
		}
	} else {
		queryInput = &dynamodb.QueryInput{
			TableName:                 aws.String(dynamoDBTableName),
			KeyConditionExpression:    aws.String(dynamoDBKeyConditionExpression),
			ExpressionAttributeValues: expressionAttributeMap,
			FilterExpression:          aws.String(dynamoDBFilterExpression),
		}
	}

	// Prepare and execute the DynamoDB query
	var queryOutput, err1 = dynamoService.Query(queryInput)
	if err1 != nil {
		log.Errorf("Error while executing query [%s]", err1)
	} else {
		result := make([]map[string]interface{}, len(queryOutput.Items))
		// Loop over the result items and build a new map structure from it
		for index, element := range queryOutput.Items {
			dat := make(map[string]interface{})
			for key, value := range element {
				if value.N != nil {
					actual := *value.N
					dat[key] = actual
				}
				if value.S != nil {
					actual := *value.S
					dat[key] = actual
				}
			}
			result[index] = dat
		}
		// Set the output value in the context
		sc := *queryOutput.ScannedCount
		context.SetOutput(ovScannedCount, sc)
		// TODO: Add consumed capacity
		// Create a JSON representation from the result
		jsonString, _ := json.Marshal(result)
		var resultinterface interface{}
		d := json.NewDecoder(bytes.NewReader(jsonString))
		d.UseNumber()
		err = d.Decode(&resultinterface)
		f := map[string]interface{}{"results": resultinterface}
		context.SetOutput(ovResult, f)
	}
	// Complete the activity
	return true, nil
}
