import type {Metadata} from "next";
import {Inter} from "next/font/google";
import "./globals.css";
import {AntdRegistry} from '@ant-design/nextjs-registry';
import {NeedAuth} from "@/components/auth";
import "@/api"
import { GoogleAnalytics } from '@next/third-parties/google'

const inter = Inter({subsets: ["latin"]});

export const metadata: Metadata = {
    title: "tech-flow | 世界中の技術情報を整理する",
    description: "tech-flowは、世界中の技術情報を整理するためのプラットフォームです。",
};

export default function RootLayout({
                                       children,
                                   }: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="ja">
        <body className={inter.className}>
        <AntdRegistry>
            <NeedAuth>{children}</NeedAuth>
        </AntdRegistry>
        </body>
        <GoogleAnalytics gaId="G-RRX9QFFJCB" />
        </html>
    );
}

