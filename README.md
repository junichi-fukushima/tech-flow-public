## 世界観
世界中の技術情報を整理する

![スクリーンショット 2024-12-19 23 58 08](https://github.com/user-attachments/assets/794b427b-7bda-42fb-9b2a-999805f648eb)



## リポジトリ構成

```bash
root
│   ├── .devcontainer
│   ├── .github                        <-- GithubActions/issue_template
│   ├── backend                        <-- バックエンドリポジトリ(主にgoのAPIたち)
│   │   │   ├── template.yaml
│   │   │   ├── template-deploy.yaml
│   ├── packages                       <-- goのAPI以外のバックエンドパッケージ
│   │   │   │   ├── claude             <-- カテゴリ/タグ判定するclaudeAPI
│   │   │   │   │   ├── template.yaml
│   │   │   │   │   ├── template-deploy.yaml
│   │   │   │   ├── metarank           <-- metarank(ローカルで動かす用)
│   ├── docs                           <-- 設計ドキュメント/openApi
│   ├── frontend                       <-- フロントエンドリポジトリ(Next.js)
│   ├── operation                      <-- table/sqldata投入など
```

## 開発環境情報&設計指針
- [backend](https://github.com/junichi-fukushima/tech-flow/blob/main/backend/README.md): golang
- [フロントエンド](https://github.com/junichi-fukushima/tech-flow/blob/main/frontend/README.md): Next.js
- その他
  - [claude API](https://github.com/junichi-fukushima/tech-flow/blob/main/packages/claude/README.md): カテゴリ・タグ判定用
  - [metarank](https://github.com/junichi-fukushima/tech-flow/blob/main/packages/metarank/README.md): ランキングアルゴリズムのOSS

## 全体構成
### デプロイ機構
![image](https://github.com/user-attachments/assets/6bb25281-4345-4cb3-bead-cb9d80809184)



### アーキテクチャ
![image](https://github.com/user-attachments/assets/109f5d82-ee9b-4754-ad70-adf1786776f6)


## バックエンドのアーキテクチャスタイル
![20230621183841](https://github.com/user-attachments/assets/a52d2286-2c65-4bb0-aeca-5a1e77d1ee43)


# 関連資料
- [notion: tech-flow](https://www.notion.so/2024-koyamapbl/tech-flow-8071a4e461664ab8a46c5507416f59ef)
- [Gドライブ: tech-flow](https://drive.google.com/drive/u/1/folders/12RHd5dC39RNLxgkipuGZ7PIPPi_sDcS3)
- [figma: tech-flow](https://www.figma.com/board/z2tS99rLbqKe88BidiG04Y/%E6%8A%80%E8%A1%93%E3%82%AD%E3%83%A5%E3%83%AC%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3%E3%82%A2%E3%83%97%E3%83%AA%E3%82%B1%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3%E3%82%A4%E3%83%A1%E3%83%BC%E3%82%B8?node-id=0-1&node-type=canvas&t=3eiTSLMUFDbXRwuO-0)
- [機能一覧](https://docs.google.com/spreadsheets/d/155nc4Wu7NXVLMKXLMqCQahpgcdRng_H_izJ9Dbss0iw/edit?gid=2059564945#gid=2059564945)


# 環境立ち上げ

- goのAPI・claudeのAPI・DBの立ち上げ

```
f-junichi@MacBook-Pro-2 tech-flow %make setupAll
```

- goのAPI・DBの立ち上げ

```
f-junichi@MacBook-Pro-2 tech-flow %make setupGolang
```

- goのAPI・claudeのAPI

```
f-junichi@MacBook-Pro-2 tech-flow %make setupClaude
```

※各パッケージでやりたい場合は各lambda関数があるREADME.mdを参照
※処理途中で終了する場合でうまくいかない場合は、docker stop $(docker ps -aq)でコンテナを停止する

# Tips

以下エラーが出るときは、docker stop $(docker ps -aq)でコンテナ停止した上でAPIを立ち上げると良さそう

```
fatal error: invalid function symbol table
nitialized
stopTheWorld: not stopped (status !
runtime stack:
18 Dec 2024 12:11:11,635 [ERROR] (rapid) Init failed error=Runtime exited with error: signal: segmentation fault InvokeID=
18 Dec 2024 12:11:11,635 [ERROR] (rapid) Invoke failed error=Runtime exited with error: signal: segmentation fault InvokeID=a59acbfe-dc8b-4046-b5ea-a6abc6215018
18 Dec 2024 12:11:11,635 [ERROR] (rapid) Invoke DONE failed: Sandbox.Failure

Invalid lambda response received: Lambda response must be valid json
2024-12-18 21:11:12 127.0.0.1 - - [18/Dec/2024 21:11:12] "GET /rss HTTP/1.1" 502 -
2024-12-18 21:11:12 127.0.0.1 - - [18/Dec/2024 21:11:12] "GET /favicon.ico HTTP/1.1" 403 -
```
