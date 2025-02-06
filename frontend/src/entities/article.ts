import sanitizeHtml from 'sanitize-html';

export class Article {
    id: number;
    category: string;
    tags: string[];
    title: string;
    link: string;
    description: string;
    feed: string;
    pub_date: Date;
    image_url: string;
    created_at: Date;
    updated_at: Date;

    constructor(data: Partial<Article>) {
        this.id = data.id ?? 0;
        this.category = data.category ?? '';
        this.tags = data.tags ?? [];
        this.title = data.title ?? '';
        this.link = data.link ?? '';
        this.description = data.description ?? '';
        this.feed = data.feed ?? '';
        this.pub_date = data.pub_date ? new Date(data.pub_date) : new Date();
        this.image_url = data.image_url ?? '';
        this.created_at = data.created_at ? new Date(data.created_at) : new Date();
        this.updated_at = data.updated_at ? new Date(data.updated_at) : new Date();
    }

    static fromApiResponse(response: any): Article {
        // descriptionのHTMLタグを除去
        const description = response.description || '';
        const sanitizedDescription = sanitizeHtml(description, {
            allowedTags: [], // すべてのタグを許可しない
        });

        return new Article({
            id: response.id,
            category: response.category,
            tags: response.tags,
            title: response.title,
            link: response.link,
            description: sanitizedDescription,
            feed: response.feed,
            pub_date: response.pub_date,
            image_url: response.image_url,
            created_at: response.created_at,
            updated_at: response.updated_at,
        });
    }
}

