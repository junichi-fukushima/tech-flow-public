package http

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var defaultHeaders = map[string]string{
	"Content-Type": "application/json; charset=utf-8",
	// for CORS
	"Access-Control-Allow-Origin":      "https://techflow.tokyo",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Headers":     "*",
}

func CreateSuccessResponse(body string, customHeaders map[string]string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: http.StatusOK,
		Headers:    createHeaders(customHeaders),
	}
}

func CreateNotFoundResponse(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: http.StatusNotFound,
		Headers:    defaultHeaders,
	}
}

func CreateErrorResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       err.Error(),
		StatusCode: http.StatusInternalServerError,
		Headers:    defaultHeaders,
	}
}

func createHeaders(headers map[string]string) map[string]string {
	if headers == nil {
		return defaultHeaders
	}

	for k, v := range defaultHeaders {
		headers[k] = v
	}
	return headers
}
