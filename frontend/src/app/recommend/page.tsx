'use client';

import React from 'react';
import {Alert, Spin} from 'antd';
import {Typography} from "antd";
import MainLayout from "@/components/mainLayout";
import {useParams} from "next/navigation";
import {useArticles} from "@/hooks/useArticles";
import {Articles} from "@/components/articles";

const {Title} = Typography;

const App: React.FC = () => {
    return (
        <MainLayout>
            <div style={{marginBottom: 32}}>
                <Title level={3}>{"おすすめ"}</Title>
                <RecommandArticles/>
            </div>
        </MainLayout>
    );
};

function RecommandArticles() {
    let [articlesResponse, loading, error] = useArticles({category: "suggest", limit: 20})
    if (error) {
        return <Alert message={`エラーが発生しました。 ${error}`} type="error" style={{marginBottom: 8}}/>
    }

    if (loading || !articlesResponse) {
        return <div style={{padding: '4px 0 0 0'}}><Spin/></div>
    }

    return <div style={{padding: '4px 0 0 0'}}>
        <Articles
            articles={articlesResponse.articles}
            ranking_event_id={articlesResponse.meta?.ranking_event_id}
            loading={loading}
        />
    </div>
}


export default App;
