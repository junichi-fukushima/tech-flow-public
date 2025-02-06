import {Article} from "@/entities/article";
import {useEffect, useState} from "react";
import {getArticles, GetArticlesResponse} from "@/api/article";

interface useArticlesProps {
    category?: string;
    tag?: string;
    limit?: number;
    offset?: number;
    keyword?: string;
}

export const useArticles = (props: useArticlesProps): [GetArticlesResponse | undefined, boolean, string | null] => {
    const [articlesResponse, setArticlesResponse] = useState<GetArticlesResponse>();
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchArticles = async () => {
            try {
                setLoading(true);
                setError(null);
                const articles = await getArticles(props);
                setArticlesResponse(articles);
            } catch (err: any) {
                setError(err.message || "Failed to fetch articles");
            } finally {
                setLoading(false);
            }
        };

        fetchArticles();
    }, [props.category, props.tag, props.limit, props.offset, props.keyword]);

    return [articlesResponse, loading, error];
};
