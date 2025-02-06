'use client';

import React from 'react';
import {Alert, Divider} from 'antd';
import {categories} from "@/util/menus";
import {Typography} from "antd";
import MainLayout from "@/components/mainLayout";
import {useArticles} from "@/hooks/useArticles";
import {Articles} from "@/components/articles";
import {PickupArticle} from "@/components/pickupArticle";
import Link from 'next/link';

const {Title} = Typography;

const App: React.FC = () => {
    return (
        <MainLayout>
            {/* ピックアップ記事 */}
            <Pickup/>
            {/* divider */}
            <Divider/>
            {/* カテゴリごとの記事 */}
            <Categories/>
        </MainLayout>
    );
};

function Pickup() {
    let [articlesResponse, loading, error] = useArticles({limit: 1})
    if (error) {
        return <Alert message={`エラーが発生しました。 ${error}`} type="error" style={{marginBottom: 8}}/>
    }


    if (!loading && articlesResponse?.articles?.length === 0) {
        return <></>
    }

    const article = articlesResponse?.articles?.[0]

    return <div>
        <Title level={3}>ピックアップ記事</Title>
        <PickupArticle article={article} loading={loading}/>
    </div>
}

function Categories() {
    return ["あなたへのおすすめ", "新着", ...categories].map((category) => <CategoryArticles category={category}
                                                                                             key={category}/>)
}

function getArticleQuery(category: string) {
    const categoryQuery = (() => {
        switch (category) {
            case "あなたへのおすすめ":
                return "suggest"
            case "新着":
                return undefined
            default:
                return category
        }
    })()
    return {
        category: categoryQuery,
        limit: 8
    }
}

function CategoryArticles({category}: { category: string }) {
    let [articlesResponse, loading, error] = useArticles(getArticleQuery(category))


    return (
        <div style={{marginBottom: 32}}>
            {
                category === "あなたへのおすすめ" ?
                    <Link href={`/recommend`}>
                        <Title level={3} style={{ opacity: 1, transition: 'opacity 0.3s' }}
                            onMouseEnter={(e) => { (e.target as HTMLElement).style.opacity = '0.6' }}
                            onMouseLeave={(e) => { (e.target as HTMLElement).style.opacity = '1' }}>
                            {category}
                        </Title>
                    </Link> :
                    category === "新着" ?
                        <Link href={`/new`}>
                            <Title level={3} style={{ opacity: 1, transition: 'opacity 0.3s' }}
                                onMouseEnter={(e) => { (e.target as HTMLElement).style.opacity = '0.6' }}
                                onMouseLeave={(e) => { (e.target as HTMLElement).style.opacity = '1' }}>
                                {category}
                            </Title>
                        </Link> :
                        <Link href={`/categories/${category}`}>
                            <Title level={3} style={{ opacity: 1, transition: 'opacity 0.3s' }}
                                onMouseEnter={(e) => { (e.target as HTMLElement).style.opacity = '0.6' }}
                                onMouseLeave={(e) => { (e.target as HTMLElement).style.opacity = '1' }}>
                                {category}
                            </Title>
                        </Link>
            }
            {
                error ? <Alert message={`エラーが発生しました。 ${error}`} type="error" style={{marginBottom: 8}}/> :
                    <div style={{padding: '4px 0 0 0'}}>
                        <Articles articles={articlesResponse?.articles || []}
                                  ranking_event_id={articlesResponse?.meta?.ranking_event_id} loading={loading}/>
                    </div>
            }
        </div>
    )
}


export default App;
