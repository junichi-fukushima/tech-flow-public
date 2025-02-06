## 概要
- claude + AWS lambda

## 前提
- aws cliが導入されていること

```
# aws-sam-cliの導入をする
brew install aws-sam-cli
```

## ローカルでの動かし方
- .aws/configに以下の記述を足す

```
[profile bedrock_profile]
region = ap-northeast-1
output = json
```

※パスは、/Users/f-junichi/.aws/configにある

- ローカルで立ち上げる
```
# build-imageの作成
sam build

# localでAPI起動(goのAPIが3000で動いているので9000で起動)
sam local start-api --env-vars env.json --port 9000
```

※localで動かすには、AWS_ACCESS_KEY_IDとAWS_SECRET_ACCESS_KEYが必要！福嶋に聞いてください

## 使用しているユーザ
- [tech-flow-local](https://us-east-1.console.aws.amazon.com/iam/home?region=ap-northeast-1#/users/details/tech-flow-local?section=permissions)
  - 権限はAmazonBedrockFullAccessのみを与えている


## 試しに動かしてみる

claudeにリクエストを投げる@本番環境


```
curl -G "https://pdxo460tjj.execute-api.ap-northeast-1.amazonaws.com/Prod/claude" \
--data-urlencode "categoryList=[AI+IoT+NONE+インフラ+フレームワーク+プログラミング言語+マネジメント+最新技術]" \
--data-urlencode "description=はじめに こんにちは！ディップ株式会社でAndroidエンジニアをやっている氏家拓海(@southcloud_7960)です。 DroidKaigi 2024に参加したので、参加レポートを書こうと思います。 また、私は2022年からDroidKaigiに参加していますが、今年は新しい挑戦としてボランティアスタッフをやってみました！そのことについても触れながらご紹介していきます。 会場について 会場でそんなに写真を撮っていなかったのですが、大きなDroidKaigiのロゴがあり、みんなの写真スポットになっていました。これだけでも！と思い、私もパシャリ。名札はかわいくデコっておきました。 スタッ…" \
--data-urlencode "tagList=[Python+Java+JavaScript+Ruby+C++C#+PHP+Swift+Go+Kotlin+TypeScript+Rust+Perl+R+Dart+Lua+Haskell+Julia+Scala+React+Angular+Vue+Django+Ruby+on+Rails+Spring+Laravel+Flask+Express+ASP.NET+Svelte+Next.js+Nuxt.js+Symfony+Meteor+CodeIgniter+CakePHP+Play+Phoenix+Ember.js+AWS+Azure+GCP+Docker+Kubernetes+Ansible+Terraform+Chef+Puppet+Vagrant+OpenStack+VMware+Jenkins+Nginx+Apache+Cloudflare+ELK+Stack+Grafana+Prometheus+CI/CD+Quantum+Computing+Blockchain+Edge+Computing+5G+AR+VR+Metaverse+Web3+Cybersecurity+Biotech+Autonomous+Vehicles+3D+Printing+Nanotechnology+Robotics+Fintech+Digital+Twins+Machine+Learning+Deep+Learning+Natural+Language+Processing+Computer+Vision+Generative+AI+Reinforcement+Learning+AI+Ethics+Neural+Networks+Speech+Recognition+Chatbots+Recommendation+Systems+AI+Governance+Explainable+AI+AI+in+Healthcare+AI+in+Finance+IoT+Devices+Smart+Homes+Smart+Cities+Edge+Devices+Sensors+Actuators+Wearables+Industrial+IoT+Connected+Cars+Embedded+Systems+Wireless+Communication+M2M+IoT+Security+Zigbee+LoRaWAN+MQTT+マネジメント+NONE]" \
--data-urlencode "title=DroidKaigi2024セッション視聴とスタッフ参加レポート"

```

claudeにリクエストを投げる@ローカル環境
```
curl -G "http://127.0.0.1:9000/claude" \
--data-urlencode "categoryList=[AI+IoT+NONE+インフラ+フレームワーク+プログラミング言語+マネジメント+最新技術]" \
--data-urlencode "description=はじめに こんにちは！ディップ株式会社でAndroidエンジニアをやっている氏家拓海(@southcloud_7960)です。 DroidKaigi 2024に参加したので、参加レポートを書こうと思います。 また、私は2022年からDroidKaigiに参加していますが、今年は新しい挑戦としてボランティアスタッフをやってみました！そのことについても触れながらご紹介していきます。 会場について 会場でそんなに写真を撮っていなかったのですが、大きなDroidKaigiのロゴがあり、みんなの写真スポットになっていました。これだけでも！と思い、私もパシャリ。名札はかわいくデコっておきました。 スタッ…" \
--data-urlencode "tagList=[Python+Java+JavaScript+Ruby+C++C#+PHP+Swift+Go+Kotlin+TypeScript+Rust+Perl+R+Dart+Lua+Haskell+Julia+Scala+React+Angular+Vue+Django+Ruby+on+Rails+Spring+Laravel+Flask+Express+ASP.NET+Svelte+Next.js+Nuxt.js+Symfony+Meteor+CodeIgniter+CakePHP+Play+Phoenix+Ember.js+AWS+Azure+GCP+Docker+Kubernetes+Ansible+Terraform+Chef+Puppet+Vagrant+OpenStack+VMware+Jenkins+Nginx+Apache+Cloudflare+ELK+Stack+Grafana+Prometheus+CI/CD+Quantum+Computing+Blockchain+Edge+Computing+5G+AR+VR+Metaverse+Web3+Cybersecurity+Biotech+Autonomous+Vehicles+3D+Printing+Nanotechnology+Robotics+Fintech+Digital+Twins+Machine+Learning+Deep+Learning+Natural+Language+Processing+Computer+Vision+Generative+AI+Reinforcement+Learning+AI+Ethics+Neural+Networks+Speech+Recognition+Chatbots+Recommendation+Systems+AI+Governance+Explainable+AI+AI+in+Healthcare+AI+in+Finance+IoT+Devices+Smart+Homes+Smart+Cities+Edge+Devices+Sensors+Actuators+Wearables+Industrial+IoT+Connected+Cars+Embedded+Systems+Wireless+Communication+M2M+IoT+Security+Zigbee+LoRaWAN+MQTT+マネジメント+NONE]" \
--data-urlencode "title=DroidKaigi2024セッション視聴とスタッフ参加レポート"

```

レスポンス結果
```
{'response': ' "tagList":{"Cloud","Python"},"categoryList":{"インフラ"}'}

```
