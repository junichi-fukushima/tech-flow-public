package main

import (
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/initializer"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/junichi-fukushima/tech-flow/backend/domain/article"
	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/domain/feed"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/http"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/repository"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/restRepository"
	"github.com/junichi-fukushima/tech-flow/backend/usecase"
)

var (
	feedRepository        feed.FeedRepository
	articleRepository     article.ArticleRepository
	restArticleRepository article.RestArticleRepository
	categoryRepository    category.CategoryRepository
	tagRepository         tag.TagsRepository
	metaRankRepository    metaRank.MetaRankRepository
)

func init() {
	var err error

	// feedRepositoryのインターフェースを取得
	feedRepository, err = repository.NewFeedRepository()
	if err != nil {
		initializer.Logger.Error("failed to initialize FeedRepository", "error", err)
	}

	// articleRepositoryのインターフェースを取得
	articleRepository, err = repository.NewArticleRepository()
	if err != nil {
		initializer.Logger.Error("failed to initialize ArticleRepository", "error", err)
	}

	// NewRestArticleRepositoryのインターフェースを取得
	restArticleRepository = restRepository.NewRestArticleRepository()

	// categoryRepositoryのインターフェースを取得
	categoryRepository, err = repository.NewCategoryRepository()
	if err != nil {
		initializer.Logger.Error("failed to initialize CategoryRepository", "error", err)
	}

	// TagRepositoryのインターフェースを取得
	tagRepository, err = repository.NewTagsRepository()
	if err != nil {
		initializer.Logger.Error("failed to initialize TagsRepository", "error", err)
	}
	metaRankRepository, err = repository.NewMetaRankRepository()
	if err != nil {
		initializer.Logger.Error("failed to initialize MetaRankRepository", "error", err)
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// DIしてユースケースを呼び出す
	rssUsecase := usecase.NewRssUsecase(
		feedRepository,
		articleRepository,
		restArticleRepository,
		categoryRepository,
		tagRepository,
		metaRankRepository,
	)

	// Upsertの処理実行
	err := rssUsecase.UpsertRss()
	if err != nil {
		initializer.Logger.Error("failed to upsert rss", "error", err)
		return http.CreateErrorResponse(err), nil
	}

	return http.CreateSuccessResponse("feedの更新処理に成功しました", nil), nil
}

func main() {
	lambda.Start(handler)
}
