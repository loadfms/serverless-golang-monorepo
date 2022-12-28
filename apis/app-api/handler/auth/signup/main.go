package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/loadfms/serverless-golang-monorepo/apis/app-api/services"
	"github.com/loadfms/serverless-golang-monorepo/apis/app-api/types"
	"github.com/loadfms/serverless-golang-monorepo/repositories"
	"github.com/loadfms/serverless-golang-monorepo/utils"
)

var userRepository *repositories.UserRepository

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (retVal events.APIGatewayProxyResponse, err error) {
	var payload types.SignupRequest
	err = json.Unmarshal([]byte(req.Body), &payload)
	if err != nil {
		fmt.Println(err.Error())
		return utils.CustomErrorResponse("Não foi possível processar o pedido", http.StatusUnprocessableEntity)
	}

	payloadValidationString := payload.GetInvalidFields()
	if payloadValidationString != "" {
		return utils.CustomErrorResponse(payloadValidationString, http.StatusBadRequest)
	}

	authService, err := services.NewAuthService()
	if err != nil {
		fmt.Println(err.Error())
		return utils.DefaultErrorResponse()
	}

	svc, err := services.NewUserService(userRepository, authService)
	if err != nil {
		fmt.Println(err.Error())
		return utils.DefaultErrorResponse()
	}

	signupResponse, err := svc.Signup(ctx, payload)
	if err != nil {
		fmt.Println(err.Error())
		return utils.DefaultErrorResponse()
	}

	return signupResponse, nil
}

func main() {
	repoManager := repositories.NewRepoManager()
	userRepository = repoManager.NewUserRepository()

	lambda.Start(Handler)
}
