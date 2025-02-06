package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"

	"github.com/junichi-fukushima/tech-flow/backend/domain/article"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
	customHTTP "github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/http"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/initializer"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/restRepository"

	"github.com/junichi-fukushima/tech-flow/backend/dto"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/repository"
	"github.com/junichi-fukushima/tech-flow/backend/usecase"
)

var (
	articleUC             usecase.ArticleUsecase
	restArticleRepository article.RestArticleRepository
	impUC                 usecase.ImpUsecase
	validate              = validator.New()
	metarankUC            usecase.MetarankUsecase
	authUC                usecase.AuthUsecase
)

func init() {
	var err error
	articleRepo, err := repository.NewArticleRepository()
	if err != nil {
		panic("failed init article repo")
	}
	restArticleRepository = restRepository.NewRestArticleRepository()
	articleUC = usecase.NewArticleUsecase(
		articleRepo,
		restArticleRepository,
	)

	metaRankRepo, err := repository.NewMetaRankRepository()
	if err != nil {
		panic("failed init imp repo")
	}
	authRepo, err := repository.NewAuthRepository()
	if err != nil {
		panic("failed init auth repo")
	}
	impUC = usecase.NewImpUsecase(metaRankRepo, authRepo)

	// init metarank Usecase
	metarankRepo, err := repository.NewMetaRankRepository()
	if err != nil {
		panic("failed init metarank repo")
	}
	metarankUC = usecase.NewMetarankUsecase(metarankRepo)

	authUC = usecase.NewAuthUsecase(authRepo)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// クエリパラメータを構造体にバインド
	req, err := bindQueryParams(request.QueryStringParameters)
	if err != nil {
		return customHTTP.CreateErrorResponse(err), nil
	}

	// バリデーション実行
	if err := validate.Struct(req); err != nil {
		return customHTTP.CreateErrorResponse(err), nil
	}

	rankingEventID, err := uuid.NewUUID()
	if err != nil {
		return customHTTP.CreateErrorResponse(err), nil
	}

	// Cookie から sessionの 値を取得
	cookieHeader, ok := request.Headers["Cookie"]
	if !ok {
		return customHTTP.CreateErrorResponse(fmt.Errorf("missing cookie header")), fmt.Errorf("missing cookie header")
	}
	sessionValue, err := parseSessionValue(cookieHeader)
	if err != nil {
		return customHTTP.CreateErrorResponse(err), err
	}

	var articles []article.Article

	// user情報を取得
	user, err := authUC.FindAnonymousUser(sessionValue)
	if err != nil {
		return customHTTP.CreateErrorResponse(err), err
	}

	// カテゴリがsuggestの場合はMetaRankに問い合わせる
	if req.Category == "suggest" {

		articlesByFavoriteCategories, err := articleUC.GetArticlesByFavoriteCategories(user.ID)
		if err != nil {
			return customHTTP.CreateErrorResponse(err), err
		}

		// articlesを[]metaRank.ItemDetailに変換
		items := make([]metaRank.ItemDetail, len(articlesByFavoriteCategories))

		for i, atc := range articlesByFavoriteCategories {
			items[i] = metaRank.ItemDetail{
				Item: strconv.Itoa(int(*atc.ID)),
			}
		}

		// metarankにレコメンド記事をAPIで問い合わせ
		recommendedArticles, err := metarankUC.GetRecommendation(user.ID, items)
		if err != nil {
			return customHTTP.CreateErrorResponse(err), err
		}

		// limitの件数もしくはrecommendedArticles.Itemsの件数のどちらか多い方で記事を取得
		limit := req.Limit
		if len(recommendedArticles.Items) < req.Limit {
			limit = len(recommendedArticles.Items)
		}
		topArticles := recommendedArticles.Items[:limit]

		articleIDs := make([]int, len(topArticles))

		for i, item := range topArticles {
			itemID, _ := strconv.Atoi(item.Item)
			articleIDs[i] = itemID
		}
		// metarankの結果を元にarticle情報を取得
		articles, err = articleUC.GetArticlesByArticleIDs(articleIDs)

		if err != nil {
			initializer.Logger.Error("Failed to get articles.", "error", err)
			return customHTTP.CreateErrorResponse(err), nil
		}

		//	articlesをarticleIDsでソートする
		articlesMap := make(map[int]article.Article)
		for _, atc := range articles {
			articlesMap[int(*atc.ID)] = atc
		}
		articles = make([]article.Article, len(articleIDs))
		for i, articleID := range articleIDs {
			articles[i] = articlesMap[articleID]
		}
	} else {
		// 記事取得
		articles, err = articleUC.GetArticles(req.Category, req.Tag, req.Limit, req.Offset, req.Keyword)
		if err != nil {
			initializer.Logger.Error("Failed to get articles.", "error", err)
			return customHTTP.CreateErrorResponse(err), nil
		}
	}

	// レスポンス用に変換
	articlesDTO := make([]dto.ArticleResponse, len(articles))
	for i, atc := range articles {
		articlesDTO[i] = convertArticleToDTOArticle(atc)
	}
	response := dto.Response{
		Data: articlesDTO,
		Meta: dto.Meta{
			Total:          len(articlesDTO),
			Limit:          req.Limit,
			Offset:         req.Offset,
			RankingEventID: rankingEventID.String(),
		},
	}
	if len(articlesDTO) > 0 {
		rankingEvent, err := impUC.CreateRankingEvent(sessionValue, response, rankingEventID.String())
		if err != nil {
			return customHTTP.CreateErrorResponse(err), err
		}

		ume, err := metarankUC.GetUserMetadataEventByUserID(user.ID)
		if err != nil {
			return customHTTP.CreateErrorResponse(err), err
		}

		err = metarankUC.SendRankingFeedback(rankingEvent, ume.ID)
		if err != nil {
			return customHTTP.CreateErrorResponse(err), err
		}
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		initializer.Logger.Error("Failed to marshal response.", "error", err)
		return customHTTP.CreateErrorResponse(err), nil
	}

	return customHTTP.CreateSuccessResponse(string(responseBody), nil), nil
}

func bindQueryParams(params map[string]string) (*dto.ArticleRequest, error) {
	limit, _ := strconv.Atoi(params["limit"])
	offset, _ := strconv.Atoi(params["offset"])

	articleRequest := &dto.ArticleRequest{
		Category: params["category"],
		Tag:      params["tag"],
		Limit:    limit,
		Offset:   offset,
		Keyword:  params["keyword"],
	}
	articleRequest.SetDefaults()
	return articleRequest, nil
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

// article.Articleからdto.Articleに変換
func convertArticleToDTOArticle(article article.Article) dto.ArticleResponse {
	return dto.ArticleResponse{
		ID:          *article.ID,
		Feed:        article.Feed.Title,
		Category:    article.Category.Name,
		Tags:        convertTagsToStrings(article.Tags),
		Title:       article.Title,
		Link:        article.Link,
		Description: article.Description,
		PubDate:     article.PubDate,
		ImageURL:    *article.ImageUrl,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
}

func convertTagsToStrings(tags []*tag.Tag) []string {
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	return tagNames
}

func main() {
	lambda.Start(handler)
}
