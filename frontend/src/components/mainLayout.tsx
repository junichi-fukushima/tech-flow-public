'use client';

import React, {useEffect, useState} from 'react';
import {MenuOutlined, UserOutlined,} from '@ant-design/icons';
import {Button, Tooltip} from 'antd';
import {Layout, Menu, theme} from 'antd';
import Image from "next/image";
import {items} from "@/util/menus";
import {Input} from 'antd';
import {useRouter} from 'next/navigation'
import Link from "next/link";
import {useMediaQuery} from "react-responsive";
import GetStarted from "@/components/getStarted";

const {Search} = Input;
const {Header, Content, Sider} = Layout;

const MainLayout: React.FC<{ children: React.ReactNode }> = ({children}) => {
  const router = useRouter();
  const {
    token: {colorBgContainer},
  } = theme.useToken();
  const isMobile = useMediaQuery({
    query: '(max-width: 575px)'
  })
  const [collapsed, setCollapsed] = useState(isMobile);
  useEffect(() => {
    setCollapsed(isMobile);
  }, [isMobile]);


  const handleSearch = (keyword: string) => {
    if (keyword.trim()) {
      router.push(`/search/${encodeURIComponent(keyword)}`);
    }
  };

  return (
    <Layout>
      <Header style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: isMobile ? 'center' : 'space-between',
        background: 'white',
        borderBottom: '1px solid #f0f0f0'
      }}>
        <Link href="/">
          <div style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}>
            <Image
              src="/logo.png"
              width={156}
              height={26.4}
              alt="logo"
            />
          </div>
        </Link>
        {
          !isMobile && <div style={{display: 'flex'}}>
            <Search
              placeholder="探したい記事を入力（例：Go, Java...）"
              style={{width: 400}}
              onSearch={handleSearch}
            />
          </div>
        }

        {
          !isMobile && <div style={{display: 'flex'}}>
            <Tooltip title="会員登録は実装中です👷">
              <Button shape="circle" icon={<UserOutlined />} disabled />
            </Tooltip>
          </div>
        }
      </Header>
      <Layout>
        <Sider
          breakpoint="lg"
          collapsedWidth="0"
          width={224}
          style={{background: colorBgContainer, overflow: 'auto'}}
          collapsed={collapsed}
          collapsible
        >
          <Menu
            mode="inline"
            style={{height: '100%', borderRight: 0}}
            items={items}
            onSelect={(item) => {
              router.push(`/tags/${item.key}`);
            }}
          />
        </Sider>
        <Layout>
          <Content
            style={{
              padding: "24px 36px",
              margin: 0,
              minHeight: "calc(100vh - 64px)",
              borderLeft: '1px solid #f0f0f0',
              background: colorBgContainer,
            }}
          >
            {children}
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};


export default MainLayout;
