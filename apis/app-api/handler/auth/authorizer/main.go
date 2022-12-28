package main

import (
	"context"
	"errors"
	"fmt"

	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/loadfms/serverless-golang-monorepo/apis/app-api/services"
)

const TokenTypeAccessToken = "access_token"

type Response events.APIGatewayProxyResponse

var authService *services.Auth

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	bearerToken := event.AuthorizationToken
	tokenSlice := strings.Split(bearerToken, "Bearer ")
	var token string = tokenSlice[0]

	if len(tokenSlice) > 1 {
		token = tokenSlice[len(tokenSlice)-1]
	}

	context, err := authService.ValidateJWT(token, TokenTypeAccessToken)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	scope := map[string]interface{}{
		"pk":         context.PK,
		"token_type": context.TokenType,
	}

	wildCardArn := strings.SplitAfter(event.MethodArn, "/")[0] + "*"

	return generatePolicy("user", "Allow", wildCardArn, scope), nil
}

func generatePolicy(principalId, effect, resource string, context map[string]interface{}) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	authResponse.Context = context
	fmt.Sprintf("Auth response: %v", authResponse)

	return authResponse
}

func main() {
	authService, _ = services.NewAuthService()
	lambda.Start(handler)
}
