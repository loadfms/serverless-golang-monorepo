package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
)

var APP_TABLE = "serverless-golang-monorepo_app"

type RepoManager struct {
	dynamoConn *dynamodb.Client
}

func connectDynamo() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("sa-east-1"),
	)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	return dynamodb.NewFromConfig(cfg)
}

func NewRepoManager() *RepoManager {
	conn := connectDynamo()
	return &RepoManager{
		dynamoConn: conn,
	}
}

func (manager *RepoManager) NewUserRepository() *UserRepository {
	return newUserRepository(manager.dynamoConn)
}
