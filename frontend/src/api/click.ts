import axios from './';

interface postClickProps {
    fields: any[];
    ranking_event_id?: string;
    article_id: number;
}

// 記事取得を行う
export const postClick = async (props: postClickProps): Promise<string> => {
    if (!props.ranking_event_id || props.ranking_event_id === "") {
        delete props.ranking_event_id
    }
    const response = await axios.post('/clicks', props, {withCredentials: true});
    return response.data
}
