import axios from './';
import {findCategoryIdByCategory} from "@/util/menus";

interface postFavoritesCategoriesProps {
  liked_categories: string[];
}

// ユーザのお気に入りカテゴリを登録する
export const postFavoritesCategories = async (props: postFavoritesCategoriesProps): Promise<boolean> => {
  const categoryIds = toCategoryIds(props.liked_categories)
  const response = await axios.post('/favorites/categories', {
    liked_categories: categoryIds
  }, {withCredentials: true})
  return response.data
}

const toCategoryIds = (categories: string[]): number[] => {
  return categories.map((category) => findCategoryIdByCategory(category)).filter((categoryId) => categoryId !== undefined)
}
