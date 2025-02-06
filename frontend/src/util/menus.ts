import type {MenuProps} from "antd";
import {
  CloudOutlined,
  FireOutlined,
  LaptopOutlined,
  OpenAIOutlined,
  ToolOutlined,
  TruckOutlined
} from "@ant-design/icons";
import React from "react";

const convertStringToTags = (str: string) => {
  return str.split(', ').map((tag) => tag.trim())
}

// 1,プログラミング言語
// 2,フレームワーク
// 3,インフラ
// 4,最新技術
// 5,AI
// 6,IoT
// 7,マネジメント
// 8,NONE
const menus = [
  {
    categoryId: 1,
    category: 'プログラミング言語',
    icon: React.createElement(LaptopOutlined),
    tags: convertStringToTags('Python, Java, JavaScript, Ruby, C++, C#, PHP, Swift, Go, Kotlin, TypeScript, Rust, Perl, R, Dart, Lua, Haskell, Julia, Scala')
  },
  {
    categoryId: 2,
    category: 'フレームワーク',
    icon: React.createElement(ToolOutlined),
    tags: convertStringToTags('React, Angular, Vue, Django, Ruby on Rails, Spring, Laravel, Flask, Express, ASP.NET, Svelte, Next.js, Nuxt.js, Symfony, Meteor, CodeIgniter, CakePHP, Play, Phoenix, Ember.js')
  },
  {
    categoryId: 3,
    category: 'インフラ',
    icon: React.createElement(CloudOutlined),
    tags: convertStringToTags('AWS, Azure, GCP, Docker, Kubernetes, Ansible, Terraform, Chef, Puppet, Vagrant, OpenStack, VMware, Jenkins, Nginx, Apache, Cloudflare, ELK Stack, Grafana, Prometheus, CI/CD')
  },
  {
    categoryId: 4,
    category: '最新技術',
    icon: React.createElement(FireOutlined),
    tags: convertStringToTags('Quantum Computing, Blockchain, Edge Computing, 5G, AR, VR, Metaverse, Web3, Cybersecurity, Biotech, Autonomous Vehicles, 3D Printing, Nanotechnology, Robotics, Fintech, Digital Twins')
  },
  {
    categoryId: 5,
    category: 'AI',
    icon: React.createElement(OpenAIOutlined),
    tags: convertStringToTags('Machine Learning, Deep Learning, Natural Language Processing, Neural Networks, GPT, BERT, Reinforcement Learning, Image Recognition, Speech Recognition, Computer Vision, Predictive Analytics, AI Ethics, Data Mining, Generative AI, Chatbots, Autonomous Systems')
  },
  {
    categoryId: 6,
    category: 'IoT',
    icon: React.createElement(TruckOutlined),
    tags: convertStringToTags('IoT Devices, Smart Homes, Smart Cities, Edge Devices, Sensors, Actuators, Wearables, Industrial IoT, Connected Cars, Embedded Systems, Wireless Communication, M2M, IoT Security, Zigbee, LoRaWAN, MQTT')
  },
]

export const categories = menus.map((category) => category.category)

export const findCategoryIdByCategory = (category: string) => {
  return menus.find((menu) => menu.category === category)?.categoryId
}

export const items: MenuProps['items'] = menus.map((category) => {
  return {
    key: category.category,
    label: category.category,
    icon: category.icon,
    children: category.tags.map((tag) => {
      return {
        key: tag,
        label: tag
      }
    })
  }
})
