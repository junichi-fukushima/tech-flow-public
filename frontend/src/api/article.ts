import axios from './';
import {Article} from "@/entities/article";

interface getArticlesProps {
    category?: string;
    tag?: string;
    limit?: number;
    offset?: number;
    keyword?: string;
}

// 記事取得を行う
export const getArticles = async (props: getArticlesProps): Promise<GetArticlesResponse> => {
    // remove comments for dummy data
    // return dummyGetArticleResponse // ランダムに並び替える

    const response = await axios.get('/articles', {params: props, withCredentials: true});
    return GetArticlesResponse.fromApiResponse(response.data)
}


export class GetArticlesResponse {
    articles: Article[];
    meta: {
        total: number;
        limit: number;
        offset: number;
        ranking_event_id: string;
    }

    constructor(data: Partial<GetArticlesResponse>) {
        this.articles = data.articles ?? [];
        this.meta = data.meta ?? {
            total: 0,
            limit: 0,
            offset: 0,
            ranking_event_id: '',
        }
    }

    static fromApiResponse(response: any): GetArticlesResponse {
        return new GetArticlesResponse({
            articles: response.data.map((item: any) => Article.fromApiResponse(item)),
            meta: response.meta,
        });
    }
}


const dummyArticles = [
    new Article({
        id: 1,
        category: 'インフラ',
        tags: ['AWS'],
        title: 'Amazon RDS for Db2: S3経由のバックアップ/リストアでデータ移行',
        link: 'https://qiita.com/nishikyon/items/a1b613442c61e0477f0b',
        description: 'Amazon RDS for Db2のデータ移行をAmazon\n' +
            '                        S3経由で行う方法が説明されています。S3にバックアップを保存し、それをRDSにリストアするプロセスを詳しく解説しています。AWS\n' +
            '                        CLIを使ったコマンドや設定の手順も含まれており、スケーラブルで柔軟なデータ移行方法として推奨されています。',
        pub_date: new Date(),
        image_url: '/dummy-2.svg',
        created_at: new Date(),
        updated_at: new Date(),
    }),
    new Article({
        id: 2,
        category: 'プログラミング言語',
        tags: ['Python'],
        title: 'Pythonでのデバッグ、print()からic()に置き換えよう！',
        link: 'https://qiita.com/ryosuke_ohori/items/11b2ad43f1ae50f25cf5',
        description: 'Amazon RDS for Db2のデータ移行をAmazon\n' +
            '                        S3経由で行う方法が説明されています。S3にバックアップを保存し、それをRDSにリストアするプロセスを詳しく解説しています。AWS\n' +
            '                        CLIを使ったコマンドや設定の手順も含まれており、スケーラブルで柔軟なデータ移行方法として推奨されています。',
        pub_date: new Date(),
        image_url: '/dummy-2.svg',
        created_at: new Date(),
        updated_at: new Date(),
    }),
    new Article({
        id: 3,
        category: 'プログラミング言語',
        tags: ['Go'],
        title: '開発用適当ツールはGoで作るのがオススメ',
        link: 'https://qiita.com/ssc-ksaitou/items/6c66669f1672806ac9bb',
        description: 'Amazon RDS for Db2のデータ移行をAmazon\n' +
            '                        S3経由で行う方法が説明されています。S3にバックアップを保存し、それをRDSにリストアするプロセスを詳しく解説しています。AWS\n' +
            '                        CLIを使ったコマンドや設定の手順も含まれており、スケーラブルで柔軟なデータ移行方法として推奨されています。',
        pub_date: new Date(),
        image_url: '/dummy-2.svg',
        created_at: new Date(),
        updated_at: new Date(),
    }),
    new Article({
        id: 4,
        category: 'インフラ',
        tags: ['Docker'],
        title: 'Docker Desktopの代替として注目されているOrbStackについてまとめてみた',
        link: 'https://qiita.com/shota0616/items/5b5b74d72272627e0f5a',
        description: 'Amazon RDS for Db2のデータ移行をAmazon\n' +
            '                        S3経由で行う方法が説明されています。S3にバックアップを保存し、それをRDSにリストアするプロセスを詳しく解説しています。AWS\n' +
            '                        CLIを使ったコマンドや設定の手順も含まれており、スケーラブルで柔軟なデータ移行方法として推奨されています。',
        pub_date: new Date(),
        image_url: '/dummy-2.svg',
        created_at: new Date(),
        updated_at: new Date(),
    })
]

const dummyGetArticleResponse: GetArticlesResponse = new GetArticlesResponse({
    articles: dummyArticles.sort(() => Math.random() - 0.5),
    meta: {
        total: 4,
        limit: 10,
        offset: 0,
        ranking_event_id: "",
    }
})
