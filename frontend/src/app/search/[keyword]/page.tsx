'use client';

import React from 'react';
import {Alert} from 'antd';
import {Typography} from "antd";
import MainLayout from "@/components/mainLayout";
import {useParams} from "next/navigation";
import {useArticles} from "@/hooks/useArticles";
import {Articles} from "@/components/articles";

const {Title} = Typography;

const App: React.FC = () => {
    const params = useParams()
    const keyword = decodeURIComponent(params.keyword as string)

    return (
        <MainLayout>
            <div>
                <div style={{marginBottom: 32}}>
                    <SearchArticles keyword={keyword}/>
                </div>
            </div>
        </MainLayout>
    );
};


function SearchArticles({ keyword }: { keyword: string }) {
    let [articlesResponse, loading, error] = useArticles({
        keyword: keyword,
        limit: 16
    })
    if (error) {
        return <Alert message={`SearchArticlesエラーが発生しました。 ${error}`} type="error" style={{marginBottom: 8}}/>
    }

    if (loading || !articlesResponse) {
        return <div style={{padding: '4px 0 0 0'}}>Loading...</div>
    }


    return (
        <div style={{marginBottom: 32}}>
            <Title level={3}>{keyword}の検索結果</Title>
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
