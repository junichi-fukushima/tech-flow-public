### マイグレーションファイルの作成

```sh
migrate create -ext sql -dir . -seq migaration_file_name
```

### マイグレーションの実行

```sh
migrate -database="mysql://ユーザー名:パスワード@ホスト名:ポート番号/データベース名" -path=./ up
```