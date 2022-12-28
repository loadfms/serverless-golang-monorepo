package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/loadfms/serverless-golang-monorepo/entities"
	"github.com/loadfms/serverless-golang-monorepo/utils"
)

const USER_PK_PREFIX = "USER#"
const PROFILE_SK_PREFIX = "PROFILE#"

type UserRepository struct {
	conn *dynamodb.Client
}

func newUserRepository(dynamoConn *dynamodb.Client) *UserRepository {
	return &UserRepository{
		conn: dynamoConn,
	}
}

func (repo *UserRepository) CreateOrUpdate(ctx context.Context, payload entities.User) (userPK string, err error) {
	if payload.PK == "" {
		payload.PK = fmt.Sprintf("%s%s", USER_PK_PREFIX, utils.GenerateSHA256(payload.Email))
	}

	if payload.SK == "" {
		payload.SK = fmt.Sprintf("%s%s", PROFILE_SK_PREFIX, utils.GenerateSHA256(payload.Email))
	}

	dynamoItem, err := attributevalue.MarshalMap(payload)
	if err != nil {
		return userPK, err
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(APP_TABLE),
		Item:      dynamoItem,
	}

	_, err = repo.conn.PutItem(ctx, params)
	if err != nil {
		return userPK, err
	}

	return payload.PK, nil
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (retVal entities.User, err error) {
	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", USER_PK_PREFIX, utils.GenerateSHA256(email))},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", PROFILE_SK_PREFIX, utils.GenerateSHA256(email))},
	}

	dynamoResult, err := repo.conn.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: &APP_TABLE,
	})
	if err != nil {
		return retVal, err
	}

	err = attributevalue.UnmarshalMap(dynamoResult.Item, &retVal)
	if err != nil {
		return retVal, err
	}

	return retVal, nil
}

func (repo *UserRepository) GetByPK(ctx context.Context, pk string) (retVal entities.User, err error) {
	keyCond := expression.KeyAnd(
		expression.Key("pk").Equal(expression.Value(pk)),
		expression.Key("sk").BeginsWith(PROFILE_SK_PREFIX))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		Build()
	if err != nil {
		return retVal, err
	}

	dynamoResult, err := repo.conn.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(APP_TABLE),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return retVal, err
	}

	if len(dynamoResult.Items) == 0 {
		return retVal, nil
	}

	err = attributevalue.UnmarshalMap(dynamoResult.Items[0], &retVal)
	if err != nil {
		return retVal, err
	}

	return retVal, nil
}
