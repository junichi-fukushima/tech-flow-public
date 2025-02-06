# 使用技術
- golang: 1.23.2
- gorm: 1.25
- mysql 8.0



# バックエンド開発の仕方
## ローカル環境構築
### devcontainerを使う場合
vs codeでdevcontainerを起動する

### ローカルで構築する場合
```
# aws-sam-cliの導入をする
brew install aws-sam-cli

# DBの立ち上げをしておく
cd backend
docker compose up
```

### 共通
```
# backend配下へ移動する
cd backend

# build-imageの作成
sam build

# localでAPI起動
sam local start-api --env-vars env.json
```

### 注意点

- goのライブラリ関係は、関数ごとに作らず以下のフォルダで共通管理します。
```
backend
```


## 本番環境へのデプロイ
```
sam deploy --guided
```

- backend/samconfig.tomlの以下のセクションにDBの接続情報が決まったら書く。(アプリケーション側の調整の可能性あり)

```
[default.build.parameters]
cached = true
parallel = true
DB_USERNAME= {{TODO}}
DB_PASSWORD= {{TODO}}
DB_HOST= {{TODO}}
DB_PORT= {{TODO}}
DB_NAME= {{TODO}}
```

## バックエンドのアーキテクチャスタイル
![20230621183841](https://github.com/user-attachments/assets/a52d2286-2c65-4bb0-aeca-5a1e77d1ee43)

## リポジトリ構成
```bash
backend
│   ├── handlers                 <-- Lambda関数格納フォルダ
│   │   ├── clicks
│   │   │   ├── main.go
│   │   ├── impressions
│   │   │   ├── main.go
│   │   ├── articles
│   │   │   ├── main.go
│   │   ├── rss
│   │   │   ├── main.go
│   │   ├── ...(その他)
│   │   │   ├── main.go
│   ├── dto                     <-- data tramsfer object
│   ├── usecase                 <-- ユースケース
│   ├── domain                  <-- ドメイン層
│   │   │   article
│   │   │   │   article_model.go
│   │   │   │   article_repositoy.go
│   │   │   category
│   │   │   │   category_model.go
│   │   │   │   category_repositoy.go
│   │   │   feed
│   │   │   │   feed_model.go
│   │   │   │   feed_repositoy.go
│   │   │   tag
│   │   │   │   tag_model.go
│   │   │   │   tag_repositoy.go
│   ├── infrastructure          <-- 永続化層/外部API/DBへのアクセス
│   │   │   external            <-- 外部API/DB
│   │   │   │   gorm
│   │   │   repository          <-- 永続化層へのアクセス
│   │   │   │   article_repositoy.go
│   │   │   │   category_repositoy.go
│   │   │   │   feed_repositoy.go
│   │   │   │   tag_repositoy.go
│   ├── env.json                <-- ローカルで読み込む環境変数
│   ├── docker-compose.yml      <-- DBのみ
│   ├── go.mod
│   ├── go.sum
│   ├── template.yaml           <-- Lambda関数ビルド用のファイル
│   ├── README.md               <-- バックエンド用のReadMe
```

## 設計思想
- クリーンアーキテクチャを設計思想として実装ています。
  - dto: バリデーション / ユースケースが使いやすいようなデータ構造に返却します
  - usecase: アプリケーションロジック
  - domain
    - model: データロジックなど
    - repository: 永続化層のインターフェース(基本永続化そうはここからアクセス！)
  - infrastructure: 永続化層や外部APIへのアクセス
