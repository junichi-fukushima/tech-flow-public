import {Article as ArticleEntity} from "@/entities/article";
import {Article} from "@/components/article";
import {Col, Row} from "antd";
import React from "react";

export function Articles({
                             articles,
                             ranking_event_id,
                             loading
                         }: {
    articles: ArticleEntity[],
    ranking_event_id?: string,
    loading?: boolean
}) {
    if (!loading && articles.length === 0) {
        return <div>記事がありません。</div>
    }
    return <Row gutter={[16, 16]}>
        {
            loading ? [1, 2, 3, 4].map(n =>
                    <Col
                        className="gutter-row"
                        xs={24} sm={8} md={6}
                        key={n}
                    >
                        <Article loading={loading}/>
                    </Col>)
                :
                articles.map((a) =>
                    <Col
                        xs={24} sm={8} md={6}
                        className="gutter-row"
                        key={a.id}
                    >
                        <Article ranking_event_id={ranking_event_id} article={a}/>
                    </Col>)
        }
    </Row>
}
