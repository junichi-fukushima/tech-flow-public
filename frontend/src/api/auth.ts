import axios from './';
import {GetArticlesResponse} from "@/api/article";

// 匿名認証を行う
export const getAuth = async (): Promise<boolean> => {
    const response = await axios.get('/auth', {withCredentials: true})
    return response.data.has_favorite_categories
}
