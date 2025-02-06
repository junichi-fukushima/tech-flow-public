package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
	customHTTP "github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/http"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/repository"
	"github.com/junichi-fukushima/tech-flow/backend/usecase"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var authUC usecase.AuthUsecase
var metarankUC usecase.MetarankUsecase

func init() {
	var err error

	authRepo, err := repository.NewAuthRepository()
	if err != nil {
		panic("failed init auth repo")
	}
	authUC = usecase.NewAuthUsecase(authRepo)

	metarankRepo, err := repository.NewMetaRankRepository()
	if err != nil {
		panic("failed init metarank repo")
	}
	metarankUC = usecase.NewMetarankUsecase(metarankRepo)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token, favoriteCategory, err := generateToken(request)
	if err != nil {
		return customHTTP.CreateErrorResponse(err), nil
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	return customHTTP.CreateSuccessResponse(
		fmt.Sprintf("{\"has_favorite_categories\":%t}", favoriteCategory),
		map[string]string{
			"Set-Cookie": cookie.String(),
		}), nil
}

// generateToken はリクエストからユーザを特定し、トークンを生成する。
func generateToken(request events.APIGatewayProxyRequest) (string, bool, error) {
	if user := findUser(request); user != nil {
		return user.SessionToken, user.HasFavoriteCategories, nil
	}

	// トークンを生成
	token, err := GenerateToken()
	if err != nil {
		return "", false, errors.New("failed to generate token")
	}

	// 匿名ユーザを作成
	err = authUC.CreateAnonymousUser(token)
	if err != nil {
		return "", false, errors.New("failed to create anonymous user")
	}

	// 匿名ユーザを取得
	user, err := authUC.FindAnonymousUser(token)
	if err != nil {
		return "", false, errors.New("failed to find anonymous user")
	}

	// ユーザメタデータを取得
	ume, err := authUC.FindUserMetaData(user.ID)
	if err != nil {
		return "", false, errors.New("failed to find user metadata")
	}

	// ユーザメタデータを送信
	if err := metarankUC.SendUserFeedback(ume); err != nil {
		return "", false, err
	}

	return token, false, nil
}

func main() {
	lambda.Start(handler)
}

func findUser(request events.APIGatewayProxyRequest) *auth.User {
	cookieHeader, hasCookie := request.Headers["Cookie"]
	if hasCookie {
		// Cookie から sessionの 値を取得
		if sessionValue, err := parseSessionValue(cookieHeader); err == nil {
			// ユーザが存在するかを確認
			if user, _ := authUC.FindAnonymousUser(sessionValue); user != nil {
				return user
			}
		}
	}
	return nil
}

func GenerateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// output sample: 4fa92b6712cd8f3491acde53617799ba
	return hex.EncodeToString(b), nil
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
