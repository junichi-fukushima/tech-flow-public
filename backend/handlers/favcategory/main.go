package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/http"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/initializer"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/repository"
	"github.com/junichi-fukushima/tech-flow/backend/usecase"
)

var userRepo auth.UserRepository
var authRepo auth.AuthRepository
var userUC usecase.UserUsecase
var authUC usecase.AuthUsecase

type FavCategoryRequest struct {
	Liked_categories []int `json:"liked_categories"`
}

func init() {
	// user
	userRepo, err := repository.NewUserRepository()
	if err != nil {
		initializer.Logger.Error("failed to initialize FeedRepository", "error", err)
	}
	userUC = usecase.NewUserUsecase(userRepo)

	// auth
	authRepo, err = repository.NewAuthRepository()
	if err != nil {
		panic("failed init authRepo")
	}
	authUC = usecase.NewAuthUsecase(authRepo)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var favCategoryRequest FavCategoryRequest
	err := json.Unmarshal([]byte(request.Body), &favCategoryRequest)

	cookieHeader, ok := request.Headers["Cookie"]
	if !ok {
		return http.CreateErrorResponse(fmt.Errorf("missing cookie header")), fmt.Errorf("missing cookie header")
	}

	// get session from cookie
	sessionValue, err := parseSessionValue(cookieHeader)
	if err != nil {
		return http.CreateErrorResponse(err), err
	}

	// find user
	user, err := authUC.FindAnonymousUser(sessionValue)
	if err != nil {
		return http.CreateErrorResponse(err), err
	}

	// 好みのカテゴリー更新処理
	err = userUC.UpsertFavCategories(*user, favCategoryRequest.Liked_categories)
	if err != nil {
		return http.CreateErrorResponse(err), err
	}

	return http.CreateSuccessResponse("success", nil), nil
}

func main() {
	lambda.Start(handler)
}

func parseSessionValue(cookieHeader string) (string, error) {
	cookieName := "session_token"
	cookies := strings.Split(cookieHeader, "; ")
	for _, cookie := range cookies {
		if strings.HasPrefix(cookie, cookieName+"=") {
			return strings.TrimPrefix(cookie, cookieName+"="), nil
		}
	}
	return "", fmt.Errorf("session cookie not found")
}
