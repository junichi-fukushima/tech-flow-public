import {Skeleton, Space, Tooltip, Typography} from "antd";
import Link from "next/link";
import {ArrowDownOutlined, ArrowUpOutlined} from "@ant-design/icons";
import React from "react";
import {Article as ArticleEntity} from "@/entities/article";
import {postClick} from "@/api/click";
import {useMediaQuery} from 'react-responsive'
import {Article} from "@/components/article";
import {EllipsisText} from "@/components/EllipsisText";

const {Text, Title, Paragraph} = Typography;

export function PickupArticle({article, ranking_event_id, loading}: {
  article?: ArticleEntity,
  ranking_event_id?: string,
  loading?: boolean
}) {
  const isMobile = useMediaQuery({
    query: '(max-width: 575px)'
  })
  if (isMobile) {
    return <Article article={article} ranking_event_id={ranking_event_id} loading={loading} />
  }

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
  return (
    <div style={{display: 'flex', padding: '4px 0 0 0'}}>
      {
        loading ? <Skeleton.Image active style={{width: 315, height: 165}} />
          : <Link
            href={link}
            target="_blank"
            onClick={post}
          >
            <div style={
              {
                width: '100%',
                maxWidth: 315,
                maxHeight: 165,
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
                  minWidth: '315px', /* 最小幅を設定 */
                }}
                alt="ピックアップ記事のカバー画像"
              />
            </div>
          </Link>
      }
      <div style={{
        padding: "4px 24px",
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        flex: 1
      }}>
        <div>
          {
            loading ?
              <Skeleton loading={loading} active paragraph={{rows: 2}} />
              : <><Title level={4}>
                <Link
                  href={link}
                  target="_blank"
                  onClick={post}
                >
                  {article?.title}
                </Link>
              </Title>
                <Paragraph>
                  {article?.description}
                </Paragraph>
              </>
          }
        </div>
        <div style={{width: '100%'}}>
          <Space direction="vertical" style={{width: '100%'}}>
            {
              article ? <div style={{width: '100%'}}>
                  <EllipsisText text={article?.feed} />
                  <Text type="secondary">{article?.pub_date.toLocaleDateString()}</Text>
                </div> :
                <Skeleton title={{width: 70}} paragraph={false} loading={loading} active />
            }
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
      </div>
    </div>
  )
}
