package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func GetClient(defaultTimeOutInSeconds int) *http.Client {
	client := new(http.Client)
	client.Timeout = time.Duration(int64(defaultTimeOutInSeconds) * int64(time.Second))

	return client
}

func headers() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
		"Content-Type":                     "application/json",
	}
}

func SuccessResponse(value any) (retVal events.APIGatewayProxyResponse, err error) {
	body, err := generateBody(value)
	if err != nil {
		return retVal, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers(),
		Body:       body,
	}, nil
}

func NoContentResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
		Headers:    headers(),
	}, nil
}

type defaultErrorResponse struct {
	Fields  *[]string `json:"fields,omitempty"`
	Message string    `json:"message"`
}

func DefaultErrorResponse() (retVal events.APIGatewayProxyResponse, err error) {
	body, err := generateBody(defaultErrorResponse{
		Message: "Erro interno do servidor",
	})

	if err != nil {
		return retVal, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Headers:    headers(),
		Body:       body,
	}, nil
}

func CustomErrorResponse(message string, status int) (retVal events.APIGatewayProxyResponse, err error) {
	body, err := generateBody(defaultErrorResponse{
		Message: message,
	})

	if err != nil {
		return retVal, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    headers(),
		Body:       body,
	}, nil
}

func generateBody(value any) (string, error) {
	body := ""

	if reflect.TypeOf(value).String() != "string" {
		byteOutput, err := json.Marshal(value)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		json.HTMLEscape(&buf, byteOutput)

		body = buf.String()
	} else {
		body = value.(string)
	}

	return body, nil
}
