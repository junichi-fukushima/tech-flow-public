package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/junichi-fukushima/tech-flow/backend/domain/auth"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/http"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/repository"
	"github.com/junichi-fukushima/tech-flow/backend/usecase"
)

var metaRankRepo metaRank.MetaRankRepository
var authRepo auth.AuthRepository
var clickUC usecase.ClickUsecase
var authUC usecase.AuthUsecase
var impUC usecase.ImpUsecase
var metarankUC usecase.MetarankUsecase

type ClickRequest struct {
	Fields         *json.RawMessage `json:"fields"`
	RankingEventID *string          `json:"ranking_event_id"`
	ArticleID      int64            `json:"article_id"`
}

func (i *ClickRequest) toInteractionEvent() (*metaRank.InteractionEvent, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	item, err := impUC.GetItemMetadataEventByArticleID(i.ArticleID)
	if err != nil {
		return nil, err
	}

	return &metaRank.InteractionEvent{
		ID:                  id.String(),
		Timestamp:           time.Now(),
		Fields:              i.Fields,
		RankingEventID:      i.RankingEventID,
		ItemMetadataEventID: item.ID,
		EventType:           metaRank.EventClick,
	}, nil
}

func init() {
	var err error
	// init click Usecase
	metaRankRepo, err = repository.NewMetaRankRepository()
	if err != nil {
		panic("failed init metaRankRepo")
	}
	clickUC = usecase.NewClickUsecase(metaRankRepo)

	// init auth Usecase
	authRepo, err = repository.NewAuthRepository()
	if err != nil {
		panic("failed init authRepo")
	}
	authUC = usecase.NewAuthUsecase(authRepo)
	impUC = usecase.NewImpUsecase(metaRankRepo, authRepo)

	// init metarank Usecase
	metarankRepo, err := repository.NewMetaRankRepository()
	if err != nil {
		panic("failed init metarank repo")
	}
	metarankUC = usecase.NewMetarankUsecase(metarankRepo)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var clickReq ClickRequest
	err := json.Unmarshal([]byte(request.Body), &clickReq)

	clickEvent, err := clickReq.toInteractionEvent()
	if err != nil {
		return http.CreateErrorResponse(err), err
	}

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

	// find meta user data
	metaUser, err := authUC.FindUserMetaData(user.ID)
	if err != nil {
		return http.CreateErrorResponse(err), err
	}

	clickEvent.UserMetadataEventID = metaUser.ID

	// save click event
	err = clickUC.CreateInteractionEvent(*clickEvent)
	if err != nil {
		return http.CreateErrorResponse(err), err
	}

	// send to metarank
	if err := metarankUC.SendInteractionFeedback(clickEvent, user.ID, strconv.Itoa(int(clickReq.ArticleID))); err != nil {
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
