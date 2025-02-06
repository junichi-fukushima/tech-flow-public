package restRepository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/junichi-fukushima/tech-flow/backend/domain/article"
	"github.com/junichi-fukushima/tech-flow/backend/domain/category"
	"github.com/junichi-fukushima/tech-flow/backend/domain/tag"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/initializer"
)

type RestArticleRepositoryImpl struct{}

type RequestPayload struct {
	QueryStringParameters QueryStringParameters `json:"queryStringParameters"`
}

type QueryStringParameters struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	CategoryList string `json:"categoryList"`
	TagList      string `json:"tagList"`
}

type ResponsePayload struct {
	Response struct {
		TagList      []string `json:"tagList"`
		CategoryList []string `json:"categoryList"`
	} `json:"response"`
}

func NewRestArticleRepository() article.RestArticleRepository {
	return &RestArticleRepositoryImpl{}
}

func (r *RestArticleRepositoryImpl) GetTagAndCategoryByClaudeAPI(
	title string,
	description string,
	categories []category.Category,
	tags []*tag.Tag,
) (category.Category, []*tag.Tag, error) {
	// Define API endpoint
	claude_endpoint_api := os.Getenv("CLAUDE_ENDPOINT_API")
	if claude_endpoint_api == "" {
		// 本番で動くことを優先する
		claude_endpoint_api = "https://pdxo460tjj.execute-api.ap-northeast-1.amazonaws.com/Prod/claude"
	}

	// Convert category and tag lists to JSON
	categoryList := make([]string, len(categories))
	for i, category := range categories {
		categoryList[i] = category.Name
	}

	tagList := make([]string, len(tags))
	for i, tag := range tags {
		tagList[i] = tag.Name
	}

	query := url.Values{}
	query.Set("title", title)
	query.Set("description", description)
	query.Set("categoryList", fmt.Sprintf("%v", categoryList))
	query.Set("tagList", fmt.Sprintf("%v", tagList))

	// 取得に失敗した場合、カテゴリー・タグはNONEとする
	// TODO: 事実上不要な処理となったので、消してしまって良い
	categoryNotFound := category.Category{
		ID:   8,
		Name: "NONE",
	}
	tagNotFound := []*tag.Tag{
		{
			ID:         108,
			Name:       "NONE",
			CategoryID: 8,
		},
	}

	urlWithParams := fmt.Sprintf("%s?%s", claude_endpoint_api, query.Encode())

	// Create HTTP request
	req, err := http.NewRequest("GET", urlWithParams, nil)
	if err != nil {
		initializer.Logger.Error("Error creating request:", "error", err)
		return category.Category{}, nil, err
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		initializer.Logger.Error("Error sending request:", "error", err)
		return category.Category{}, nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		initializer.Logger.Error("Non-200 response from Claude API:",
			"status", resp.StatusCode,
			"description", description,
		)

		return categoryNotFound, tagNotFound, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// Read and inspect response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		initializer.Logger.Error("Error reading response body: ", "error", err)
		return categoryNotFound, tagNotFound, err
	}

	// Unmarshal response body
	var responsePayload ResponsePayload
	if err := json.Unmarshal(body, &responsePayload); err != nil {
		initializer.Logger.Error("Error unmarshalling response JSON:", "error", err)
		return categoryNotFound, tagNotFound, err
	}

	// responseからカテゴリ取得
	var mappedCategory category.Category
	if len(responsePayload.Response.CategoryList) > 0 {
		for _, originalCategory := range categories {
			if responsePayload.Response.CategoryList[0] == originalCategory.Name {
				mappedCategory = originalCategory
				break
			}
		}
	}
	if mappedCategory == (category.Category{}) {
		// mappedCategoryが見つからない場合、Nameが"NONE"のものを設定
		mappedCategory = categoryNotFound
	}

	// responseからタグ取得
	mappedTags := []*tag.Tag{}
	for _, responseTag := range responsePayload.Response.TagList {
		for _, originalTag := range tags {
			if responseTag == originalTag.Name {
				mappedTags = append(mappedTags, originalTag)
				break
			}
		}
	}
	// mappedTagが見つからない場合、Nameが"NONE"のものを設定
	if mappedTags == nil {
		mappedTags = tagNotFound
	}

	return mappedCategory, mappedTags, nil
}
