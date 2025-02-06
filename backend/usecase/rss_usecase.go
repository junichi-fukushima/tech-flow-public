package usecase

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/junichi-fukushima/tech-flow/backend/dto"

	"github.com/google/uuid"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/initializer"

	"github.com/junichi-fukushima/tech-flow/backend/domain/article"
	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/domain/feed"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
	"github.com/mmcdole/gofeed"
)

type RssUsecase interface {
	UpsertRss() error
}

type rssUsecase struct {
	feedRepository        feed.FeedRepository
	articleRepository     article.ArticleRepository
	restArticleRepository article.RestArticleRepository
	categoryRepository    category.CategoryRepository
	tagRepository         tag.TagsRepository
	metaRankRepository    metaRank.MetaRankRepository
}

func NewRssUsecase(
	feedRepository feed.FeedRepository,
	articleRepository article.ArticleRepository,
	restArticleRepository article.RestArticleRepository,
	categoryRepository category.CategoryRepository,
	tagRepository tag.TagsRepository,
	metaRankRepository metaRank.MetaRankRepository,
) RssUsecase {
	return &rssUsecase{
		feedRepository:        feedRepository,
		articleRepository:     articleRepository,
		restArticleRepository: restArticleRepository,
		categoryRepository:    categoryRepository,
		tagRepository:         tagRepository,
		metaRankRepository:    metaRankRepository,
	}
}

func (u *rssUsecase) UpsertRss() error {
	// feeds一括取得
	// 15分ごとに実行
	// 実行する時間によって分割用のシードを作ってそれに基づきidを指定してfeedを取得する
	now := time.Now()
	timeSeed := (now.Unix() / 240) % 329 // 4分ごとに処理し、329回で658件更新
	rangeSize := 2                       // 一度に取得するフィード数

	// IDの範囲を計算
	startID := int(timeSeed) * rangeSize
	endID := startID + rangeSize

	// ID範囲の調整（オーバーフロー防止）
	if endID > 658 {
		endID = 658
	}

	// 範囲でフィードを取得
	feedsByRepository, err := u.feedRepository.GetFeedsByIDRange(startID, endID)
	if err != nil {
		return err
	}

	// OGP画像取得用のusecaseを生成
	ogpUsecase := NewOgpUsecase()

	for _, feedByRepository := range feedsByRepository {
		// RSS情報を取得する
		feedByRss, err := gofeed.NewParser().ParseURL(feedByRepository.Link)
		if err != nil || feedByRss == nil {
			initializer.Logger.Warn("Failed to parse feed URL or feed is nil:", feedByRepository.Link, err)
			continue
		}

		// UpdatedParsed(=更新日時)がnilの場合は一旦スキップしとく
		if feedByRss.UpdatedParsed == nil {
			initializer.Logger.Warn("Feed UpdatedParsed is nil, skipping feed:", feedByRepository.Link)
			continue
		}

		// NOTE: 一括で全て更新したい場合はここをコメントアウトする
		// feedが更新されてない場合はスキップ
		if feedByRepository.LastBuildDate != nil &&
			(feedByRepository.LastBuildDate.After(*feedByRss.UpdatedParsed) ||
				feedByRepository.LastBuildDate.Equal(*feedByRss.UpdatedParsed)) {
			initializer.Logger.Info("info", "feed情報が最新ではないので、記事更新処理はskipしました")
			continue
		}

		// feedをドメインモデル化
		var imageURL *string
		if feedByRss.Image != nil {
			imageURL = &feedByRss.Image.URL
		}
		domainModelFeed := feed.NewFeed(
			feedByRepository.ID,
			feedByRss.Title,
			feedByRepository.Link,
			&feedByRss.Description,
			stringPointer(strings.Join(feedByRss.Categories, ",")),
			imageURL,
			&feedByRss.Language,
			feedByRss.UpdatedParsed,
			feedByRepository.CreatedAt,
			time.Now(), // UpdatedAt
		)

		// articleドメインモデル化
		var domainArticles = make([]article.Article, len(feedByRss.Items))
		var guids []string // GUIDを格納するスライス(ItemMetadataEventへの登録時に行う)

		// タグ全部取得(claudeAPIでタグを判定させる時に使う)
		tagUseCase := NewTagUsecase(u.tagRepository)
		tags, err := tagUseCase.GetAllTags()
		if err != nil {
			return err
		}

		// カテゴリ全取得(claudeAPIでタグを判定させる時に使う)
		categoryUseCase := NewCategoryUsecase(u.categoryRepository)
		categories, err := categoryUseCase.GetAllCategories()
		if err != nil {
			return err
		}

		articleUseCase := NewArticleUsecase(
			u.articleRepository,
			u.restArticleRepository,
		)

		// claudeAPIを使って判定するかどうか(基本false)
		// TODO:Error: Failed to create changeset for the stack: handlers, An error occurred (ValidationError) when calling the CreateChangeSet operation: Template format error: Parameter name USE_CLAUDE is non alphanumeric.
		isUseClaudeStr := os.Getenv("USE_CLAUDE")
		isUseClaude, err := strconv.ParseBool(isUseClaudeStr)
		if err != nil {
			// 本番で動くのをデフォルトとする
			isUseClaude = true
		}

		for i, articleByRss := range feedByRss.Items {
			// タグとカテゴリーの判定をする
			var matchedTags []*tag.Tag
			var matchedCategory category.Category

			if isUseClaude {
				// カテゴリ & タグ判定(claude API使用)
				matchedTags, matchedCategory, err = articleUseCase.GetCategoryAndTag(
					articleByRss.Title,
					&articleByRss.Description,
					tags,
					categories,
				)
				// NOTE: claudeのAPIが失敗した場合は、暫定重みづけ処理を実行
				if err != nil {
					initializer.Logger.Warn("Claude API判定に失敗したので、重みづけ判定を実施します:",
						"title", articleByRss.Title,
					)
					// タグ判定(暫定重みづけ)
					matchedTags, err = tagUseCase.DecideTags(
						articleByRss.Title,
						&articleByRss.Description,
					)
					if err != nil {
						return err
					}

					// カテゴリ判定(暫定重みづけ)
					matchedCategory, err = categoryUseCase.DecideCategory(tags)
					if err != nil {
						return err
					}

					// ログ判定よう
					tagsJSON, _ := json.Marshal(matchedTags)
					categoryJSON, _ := json.Marshal(matchedCategory)

					initializer.Logger.Info("Claude API判定に失敗したので、重みづけ判定を実施しました",
						"matchedTags", string(tagsJSON),
						"matchedCategory", string(categoryJSON),
					)
				}

			} else {
				// タグ判定(暫定重みづけ)
				matchedTags, err = tagUseCase.DecideTags(
					articleByRss.Title,
					&articleByRss.Description,
				)
				if err != nil {
					return err
				}

				// カテゴリ判定(暫定重みづけ)
				matchedCategory, err = categoryUseCase.DecideCategory(tags)
				if err != nil {
					return err
				}
			}

			if len(matchedTags) == 0 {
				initializer.Logger.Warn("タグが見つからなかったため処理を中断しました")
				break
			}

			if matchedCategory == (category.Category{}) {
				initializer.Logger.Warn("カテゴリが見つからなかったため処理を中断しました")
				break
			}

			// OGP画像を取得する
			imageUrl := ogpUsecase.FetchOGPImage(articleByRss.Link)

			domainModelArticle := article.NewArticle(article.Article{
				FeedID:      feedByRepository.ID,
				Category:    matchedCategory,
				Title:       articleByRss.Title,
				Link:        articleByRss.Link,
				Description: &articleByRss.Description,
				PubDate:     *articleByRss.PublishedParsed,
				GUID:        articleByRss.GUID,
				ImageUrl:    &imageUrl,
				Tags:        matchedTags,
				CreatedAt:   feedByRepository.CreatedAt,
				UpdatedAt:   time.Now(),
			})
			// 最新のRSS情報を元にUpsertする
			domainArticles[i] = *domainModelArticle

			// GUIDを追加
			guids = append(guids, articleByRss.GUID)
		}

		// 一括更新対象の、feedのUpsert処理実行
		err = u.feedRepository.UpsertRss(domainModelFeed)
		if err != nil {
			return err
		}

		// 記事が空の場合スキップ
		if len(domainArticles) == 0 {
			initializer.Logger.Warn("更新対象のfeedですが、更新対象の記事がありませんでした。:", feedByRepository.Link)
			continue
		}

		// 一括更新対象の、articleのUpsert処理実行
		err = u.articleRepository.UpsertArticles(domainArticles)
		if err != nil {
			return err
		}

		// guidsを元に、articleIDsを取得する
		articleIDs, err := u.articleRepository.GetArticleIDsByGUIDs(guids)
		if err != nil {
			return err
		}

		// すでに登録済みのItemMetadataEventを取得
		existingEvents, err := u.metaRankRepository.GetItemMetadataEventsByArticleIDs(articleIDs)
		if err != nil {
			return err
		}

		// 未登録のItemMetadataEventをCreateする
		existingArticleIDSet := make(map[int64]struct{})
		for _, event := range existingEvents {
			existingArticleIDSet[event.ArticleID] = struct{}{}
		}
		var newEvents []metaRank.ItemMetadataEvent
		for _, articleID := range articleIDs {
			if _, exists := existingArticleIDSet[articleID]; !exists {
				id, _ := uuid.NewUUID()
				newEvent := metaRank.ItemMetadataEvent{
					ID:        id.String(),
					Timestamp: time.Now(),
					Fields:    nil, // 必要に応じてデータを設定
					ArticleID: articleID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				newEvents = append(newEvents, newEvent)
			}
		}

		if len(newEvents) > 0 {
			if err := u.metaRankRepository.CreateItemMetadataEvents(newEvents); err != nil {
				return err
			}
			if err := u.sendItemsFeedback(newEvents); err != nil {
				return err
			}
		}
	}

	return nil
}

func stringPointer(s string) *string {
	return &s
}

// metarankにItemMetadataEventを送信する
func (u *rssUsecase) sendItemsFeedback(events []metaRank.ItemMetadataEvent) error {
	// eventのidをすべて取得
	var eventIDs []string
	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	// eventsを子テーブル、孫テーブル付きで取得
	imes, err := u.metaRankRepository.GetItemMetadataEventByIDs(eventIDs)
	if err != nil {
		return err
	}
	for _, ime := range imes {
		if err := u.sendItemFeedback(&ime); err != nil {
			return err
		}
	}
	return nil
}

func (u *rssUsecase) sendItemFeedback(event *metaRank.ItemMetadataEvent) error {
	req := (&dto.ItemMetadataEventRequest{}).FromItemMetadataEvent(event)
	return u.metaRankRepository.SendFeedback(req)
}
