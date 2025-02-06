package usecase

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/initializer"
)

type OgpUsecase interface {
	FetchOGPImage(url string) string
}

type ogpUsecase struct {
	client *http.Client
}

func NewOgpUsecase() OgpUsecase {
	return &ogpUsecase{
		client: &http.Client{
			Timeout: 3 * time.Second, // タイムアウトを3秒に設定
		},
	}
}

func (u *ogpUsecase) FetchOGPImage(url string) string {
	resp, err := u.client.Get(url)
	if err != nil {
		initializer.Logger.Warn(fmt.Sprintf("Failed to fetch URL: %v", err))
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// HTTPステータスコードが正常ではない場合は警告を出力
		initializer.Logger.Warn("HTTPステータスコードが正常ではありません。", url, fmt.Errorf("HTTP error: %s", resp.Status))
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		initializer.Logger.Warn("HTMLの解析時にエラーが発生しました。", url, fmt.Errorf("failed to parse HTML: %v", err))
		return ""
	}

	// OGP画像URLを取得
	return doc.Find("meta[property='og:image']").AttrOr("content", "")
}
