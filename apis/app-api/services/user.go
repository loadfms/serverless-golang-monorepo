package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/loadfms/serverless-golang-monorepo/apis/app-api/types"
	"github.com/loadfms/serverless-golang-monorepo/entities"
	"github.com/loadfms/serverless-golang-monorepo/repositories"
	"github.com/loadfms/serverless-golang-monorepo/utils"
	"github.com/twinj/uuid"
)

const TokenTypeAccessToken = "access_token"

type User struct {
	userRepo    *repositories.UserRepository
	authService *Auth
}

func NewUserService(userRepo *repositories.UserRepository, authSvc *Auth) (*User, error) {
	return &User{
		userRepo:    userRepo,
		authService: authSvc,
	}, nil
}

func (s *User) Signup(ctx context.Context, request types.SignupRequest) (retVal events.APIGatewayProxyResponse, err error) {

	userAlreadyRegistered, err := s.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		return retVal, err
	}
	if userAlreadyRegistered.PK != "" {
		return utils.CustomErrorResponse("Usuário já cadastrado", http.StatusConflict)
	}

	salt, err := generateRandomSalt()
	if err != nil {
		return retVal, err
	}

	hashPassword := hashPassword(request.Password, salt)

	dynamoUser := entities.User{
		Email:        request.Email,
		Name:         request.Name,
		Password:     hashPassword,
		PasswordSalt: string(salt),
	}

	userPK, err := s.userRepo.CreateOrUpdate(ctx, dynamoUser)
	if err != nil {
		return retVal, err
	}

	token, err := s.authService.GenerateJWT(userPK, 24, TokenTypeAccessToken)
	if err != nil {
		return retVal, err
	}

	return utils.SuccessResponse(types.AuthResponse{
		AccessToken: token,
	})
}

func (s *User) Signin(ctx context.Context, request types.SigninRequest) (retVal events.APIGatewayProxyResponse, err error) {

	userAlreadyRegistered, err := s.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		return retVal, err
	}
	if userAlreadyRegistered.PK == "" || userAlreadyRegistered.PasswordSalt == "" {
		return utils.CustomErrorResponse("Usuário ou senha inválido", http.StatusUnauthorized)
	}

	hashPassword := hashPassword(request.Password, userAlreadyRegistered.PasswordSalt)

	if userAlreadyRegistered.Password != hashPassword {
		return utils.CustomErrorResponse("Usuário ou senha inválido", http.StatusUnauthorized)
	}

	token, err := s.authService.GenerateJWT(userAlreadyRegistered.PK, 24, TokenTypeAccessToken)
	if err != nil {
		return retVal, err
	}

	return utils.SuccessResponse(types.AuthResponse{
		AccessToken: token,
	})
}

func hashPassword(password string, salt string) string {
	h := sha256.New()

	h.Write([]byte(fmt.Sprintf("%s%s", password, salt)))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func generateRandomSalt() (string, error) {
	return uuid.NewV4().String(), nil
}
