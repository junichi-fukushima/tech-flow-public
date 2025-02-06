## Metarank

### 各ファイルの説明
* `config.yml`
  * metarankの設定ファイル
    * 詳細：[https://docs.metarank.ai/reference/overview](https://docs.metarank.ai/reference/overview) 
* `events.json`
  * 本番DBの各イベント（`item_metadata_events`, `user_metadata_events`, `ranking_events`, `interaction_events`）を格納したファイル
    * 詳細：[https://docs.metarank.ai/reference/event-schema](https://docs.metarank.ai/reference/event-schema)

### metarankの起動
```sh
# portは8000番で起動する（ローカルのAPI Gatewayとぶつかってしまっているため）
docker run -i -t -p 8000:8080 -v $(pwd):/opt/metarank metarank/metarank:latest standalone --config /opt/metarank/config.yml --data /opt/metarank/events.json
```

### 使用例
#### 10件のデータを取得(アイテム指定なし)
* リクエスト
```sh
curl -i http://localhost:8080/recommend/trending -d '{
    "count": 10,
    "items": [],
    "user": "06f6c808-b637-11ef-806f-c6313e953dbd"
}'
```
* レスポンス
```json
{"items":[{"item":"1485","score":7.0},{"item":"6056","score":6.0},{"item":"1486","score":2.0},{"item":"1328","score":1.0},{"item":"1333","score":1.0},{"item":"2155","score":1.0},{"item":"214","score":1.0},{"item":"2923","score":1.0},{"item":"5122","score":1.0}],"took":0}
```
詳細：[https://docs.metarank.ai/reference/api](https://docs.metarank.ai/reference/api)

### その他
* configファイルの自動作成
  * 詳細： [https://docs.metarank.ai/how-to/autofeature](https://docs.metarank.ai/how-to/autofeature) 
```sh
docker run -i -t -p 8080:8080 -v $(pwd):/opt/metarank metarank/metarank:latest autofeature --data /opt/metarank/events.json --out /opt/metarank/auto_config.yaml
```
