package dto

import (
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"strconv"
	"time"
)

// ItemMetadataEventRequest は、ItemMetadataEventのリクエストを表します
// 例：
//
//	{
//	 "event": "item",
//	 "id": "acaec3aa-b441-11ef-a31e-d63b9a437fb6",
//	 "item": "2",
//	 "timestamp": "1733570363000",
//	 "fields": [
//	   {
//	     "name": "title",
//	     "value": "国内初！店舗ショールームへの MagicLeap導入を通じての学び"
//	   },
//	   {
//	     "name": "description",
//	     "value": "こんにちわ！ x gardenでSpatial ComputingのUnity開発、UXデザインを担当しているFalと申します。 今日は、先日開催されたMagic Leap Meetupで発表した内容をご紹介します。 クライアントは株式会社ニトリ様、システムキッチンをMRでシミュレーションするMagic Leap 1アプリケーションを開発したのですが、その時の学びをまとめてみました。 国内初！店舗ショールームへの MagicLeap導入を通じての学び - Google スライド 今後Magic Leapで開発するクリエイターの方へ、何かしら学びを共有できれば幸いです！ぜひご笑覧ください。"
//	   },
//	   {
//	     "name": "category",
//	     "value": "プログラミング言語"
//	   },
//	   {
//	     "name": "tags",
//	     "value": [
//	       "Go",
//	       "R"
//	     ]
//	   }
//	 ]
//	}
type ItemMetadataEventRequest struct {
	Event     string    `json:"event"`
	ID        string    `json:"id"`
	Item      string    `json:"item"`
	Timestamp time.Time `json:"timestamp"`
	Fields    []Field   `json:"fields"`
}

type Field struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

func (u *ItemMetadataEventRequest) FromItemMetadataEvent(itemMetadataEvent *metaRank.ItemMetadataEvent) *ItemMetadataEventRequest {
	fields := make([]Field, 0, 4)
	fields = append(fields, Field{
		Name:  "title",
		Value: itemMetadataEvent.Article.Title,
	})
	var description string
	if itemMetadataEvent.Article.Description != nil {
		description = *itemMetadataEvent.Article.Description
	}

	fields = append(fields, Field{
		Name:  "description",
		Value: description,
	})

	fields = append(fields, Field{
		Name:  "category",
		Value: itemMetadataEvent.Article.Category.Name,
	})

	tags := make([]string, 0, len(itemMetadataEvent.Article.Tags))
	for _, tag := range itemMetadataEvent.Article.Tags {
		tags = append(tags, tag.Name)
	}

	fields = append(fields, Field{
		Name:  "tags",
		Value: tags,
	})

	return &ItemMetadataEventRequest{
		Event:     "item",
		ID:        itemMetadataEvent.ID,
		Timestamp: itemMetadataEvent.Timestamp,
		Item:      strconv.Itoa(int(itemMetadataEvent.ArticleID)),
		Fields:    fields,
	}
}
