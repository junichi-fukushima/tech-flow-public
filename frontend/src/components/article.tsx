import {Article as ArticleEntity} from "@/entities/article";
import {Card, Skeleton, Space, Tooltip, Typography} from "antd";
import Link from "next/link";
import {ArrowDownOutlined, ArrowUpOutlined} from "@ant-design/icons";
import {postClick} from "@/api/click";
import {EllipsisText} from "@/components/EllipsisText";

const {Text} = Typography;

export function Article({article, ranking_event_id, loading}: {
  article?: ArticleEntity,
  ranking_event_id?: string,
  loading?: boolean
}) {
  const link = article?.link || "https://example.com"
  const post = async () => {
    if (!article) return
    await postClick({
      fields: [],
      ranking_event_id: ranking_event_id,
      article_id: article.id
    })
  }

  const article_image = article?.image_url || "/dummy-2.svg"

  return <Card
    title={
      <div style={{padding: '8px 0'}}>
        {loading ?
          <div style={{padding: '8px 0'}}>
            <Skeleton paragraph={false} loading={loading} active />
            <Skeleton paragraph={false} loading={loading} active style={{padding: '8px 0'}} />
            <Skeleton title={{width: 70}} paragraph={false} loading={loading} active />
          </div>
          :
          <Space direction="vertical">
            <Link
              href={link}
              target="_blank"
              onClick={post}
            >
              <div
                style={{
                  whiteSpace: "initial",
                  display: "-webkit-box",
                  WebkitBoxOrient: "vertical",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                  WebkitLineClamp: 2, // 最大2行まで表示
                  lineHeight: "1.5", // 行の高さ
                  height: "3.4em", // 2行分の高さを固定
                  fontSize: "16px", // 必要に応じてフォントサイズを設定
                  padding: "0.5em", // テキストの周囲の余白（任意）
                  boxSizing: "border-box", // パディングを含めた高さ調整
                }}>
                {article!.title}
              </div>
            </Link>
            <div style={{padding: '0 0.5em'}}>
              <EllipsisText text={article!.feed}/>
              {/* 例：2024.10.10 */}
              <Text type="secondary"
                    style={{fontWeight: "initial"}}>{article!.pub_date.toLocaleDateString()}</Text>
            </div>
          </Space>
        }
      </div>
    }
    size="small"
    styles={{
      body: {
        padding: 0,
      }
    }}
  >
    {
      article ? <>
          <Link href={link} target="_blank" onClick={post}>
            <div style={
              {
                width: '100%',
                aspectRatio: '1.91 / 1', /* 標準的なOGP比率（1.91:1）をベースに調整 */
                overflow: 'hidden',
                position: 'relative',
                background: '#f0f0f0', /* 画像がない場合の背景色 */
              }
            }>
              <img
                src={article_image}
                style={{
                  width: '100%',
                  height: '100%',
                  objectFit: "cover", /* 親要素全体を埋めるように画像を表示 */
                  objectPosition: "center",  /* 中央寄せで調整 */
                  display: 'block', /* 画像間の隙間を防ぐ */
                }}
                alt="記事のカバー画像"
              />
            </div>
          </Link>
          <div style={{padding: '8px 12px'}}>
            <Space direction="vertical">
              <Tooltip title="いいね！は実装中です👷">
                <Text type="secondary">
                  <Space>
                    <ArrowUpOutlined />
                    <ArrowDownOutlined />
                    <span>1</span>
                  </Space>
                </Text>
              </Tooltip>
            </Space>
          </div>
        </>
        :
        <img
          src={article_image}
          width={0}
          height={0}
          sizes="100vw"
          style={{width: '100%', height: 'auto'}}
          alt="記事のカバー画像"
        />
    }
  </Card>
}
