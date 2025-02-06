# feedテーブル生成スクリプト
import csv

# CSVファイルのパス
file_name = "feeds_list.csv"

# 出力するINSERT文のテーブル名
table_name = "feeds"

# CSVファイルを読み込み、INSERT文を生成
def generate_insert_statements(file_name, table_name):
    sql_statements = []
    current_id = 1  # IDの初期値

    with open(file_name, mode='r', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)

        # ヘッダーを確認
        print("ヘッダー名:", reader.fieldnames)

        for row in reader:
            try:
                # 必要なフィールドが存在するか確認
                if 'title' not in row or 'link' not in row:
                    raise KeyError("必要なカラム名がCSVファイルに存在しません: 'title', 'link'")

                title = row['title'].replace("'", "''")  # SQLインジェクション対策
                link = row['link'].replace("'", "''")    # SQLインジェクション対策

                # IDと現在時刻を含むINSERT文を作成
                sql = (
                    f"INSERT INTO {table_name} (id, title, link, created_at, updated_at) "
                    f"VALUES ({current_id}, '{title}', '{link}', NOW(), NOW());"
                )
                sql_statements.append(sql)

                # IDをインクリメント
                current_id += 1
            except KeyError as e:
                print(f"エラー: {e}")

    return "\n".join(sql_statements)

# INSERT文を生成
try:
    insert_statements = generate_insert_statements(file_name, table_name)
    print(insert_statements)
except FileNotFoundError:
    print(f"Error: The file '{file_name}' was not found.")
