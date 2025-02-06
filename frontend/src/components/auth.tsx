"use client";

import {useAuth} from "@/hooks/useAuth";
import {Alert, Spin, theme} from "antd";
import GetStarted from "@/components/getStarted";

export function NeedAuth({children}: Readonly<{ children: React.ReactNode }>) {
    const [loading, error,hasFavoriteCategories ] = useAuth();
    const {
        token: {colorBgContainer},
    } = theme.useToken();
    if (loading) {
        return <div
            style={{
                height: '100vh',
                width: '100vw',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: colorBgContainer
            }}>
            <Spin percent="auto" size="large"/>
        </div>;
    }

    return (
        <div>
            {error && <Alert message={"認証に失敗しました"} type="error"/>}
            {!hasFavoriteCategories && <GetStarted/>}
            {children}
        </div>
    );
}
