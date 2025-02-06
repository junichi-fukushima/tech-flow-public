import re
import json
import os
import time
import boto3

def safe_json_loads(json_string, default=None):
    """
    Safely parse a JSON string. Returns default if parsing fails.
    """
    try:
        return json.loads(json_string)
    except (json.JSONDecodeError, TypeError):
        return default if default is not None else []

def get_bedrock_client():
    """
    Get the Bedrock client based on the environment (local or production).
    """
    region_name = "ap-northeast-1"

    # 本番環境では IAM ロール認証
    if os.getenv("Env") == "production":
        boto3_session = boto3.session.Session(region_name=region_name)
    else:
        # ローカル環境ではアクセスキーとシークレットキーを利用
        boto3_session = boto3.session.Session(
            aws_access_key_id=os.getenv("AWS_ACCESS_KEY_ID"),
            aws_secret_access_key=os.getenv("AWS_SECRET_ACCESS_KEY"),
            region_name=region_name
        )

    return boto3_session.client(service_name="bedrock-runtime")


def lambda_handler(event, context):
    # リクエストデータ取得
    query_string = event.get("queryStringParameters", {})
    if not query_string:
        query_string = {}

    title = query_string.get("title", "")
    description = query_string.get("description", "")
    category_list_text = query_string.get("categoryList", "[]")
    tag_list_text = query_string.get("tagList", "[]")

    # プロンプトのテキストを生成
    # 文頭に\n\nHuman:
    # 文末に\n\nAssistant:を置く必要があるがそれはclaudeの仕様です。
    prompt_text = f"""
    \n\nHuman: ### やってほしいこと
    titleとdescription情報をもとに、記事に関連性があるタグとカテゴリーのみを選定してください
    その上で、最後にJSON形式のレスポンス結果のみを出力してください
    レスポンス結果は以下の形式とします。不要な改行やテキストはレスポンスに含めずプレーンなJSONを返却してください

    "tagList":{{
        "name": "選定したタグ名"
    }},
    "categoryList":{{
        "name": "選定したカテゴリー名"
    }},

    ### title情報
    {title}

    ### description情報
    {description}

    ### あなたが出力するレスポンス結果の例(あくまでsample)
    "tagList":{{
        "Python","Scala"
    }},
    "categoryList":{{
        "プログラミング言語"
    }}

    ### あなたの役割
    あなたは、技術記事のタグとカテゴリを見抜くプロです。
    ユーザーに適切なタグとカテゴリを伝える必要があります。

    ### タグ一覧
    {tag_list_text}

    ### カテゴリ一覧
    {category_list_text}
    ※当てはまるカテゴリがない場合は、"NONE"を返してください
    ※当てはまるタグがない場合は、"NONE"を返してください

    ### 制約条件
     タグは最大5個までとする
     カテゴリは最大1個までとする
     タグとカテゴリどちらも正式名称を使うようにしてください(OK: Ruby / NG: ruby)
     上記でセットしているタグの一覧を元に選定してください\n\nAssistant:
    """

    # 10秒待機を追加(Bedrockの負荷軽減のため)
    time.sleep(10)

    # claudeにプロンプトなげる
    bedrock = get_bedrock_client()

    body = json.dumps({
            "prompt": prompt_text,
            "max_tokens_to_sample": 1024,
            "temperature": 0.1,
            "top_p": 0.9,
    })

    # レスポンス
    response = bedrock.invoke_model(
        body=body,
        modelId='anthropic.claude-instant-v1',
        accept='application/json',
        contentType='application/json'
    )
    response_body = json.loads(response.get('body').read())

    # completion(プロンプトに対しての結果)の内容を取得
    completion = response_body.get('completion', '')

    # 不要な改行とかを消す
    cleaned_completion_removed_lf = re.sub(r'\s*\n\s*', '', completion)

    # シングルクォートをダブルクォートに変換
    cleaned_completion_chage_quote = cleaned_completion_removed_lf.replace("'", '"')

    # 正規表現でjson化する
    fixed_completion = re.sub(r'{([^{}]*)}', r'[\1]', cleaned_completion_chage_quote)
    data = json.loads(fixed_completion)

    # jsonからlist化してそこから値を抽出
    tag_list_by_claude = list(data.get("tagList", []))
    category_list_by_claude = list(data.get("categoryList", []))

    return {
        'statusCode': 200,
        'body': json.dumps({
            'response': {
                'tagList': tag_list_by_claude[:5],  # タグリスト
                'categoryList': category_list_by_claude[:1]  # カテゴリリスト
            }
        }, ensure_ascii=False)  # 確実に正しいJSON形式を維持
    }
