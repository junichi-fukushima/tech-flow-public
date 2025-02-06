"use client"

import {useEffect, useState} from "react";
import {getAuth} from "@/api/auth";

export const useAuth = (): [boolean, string | null, boolean] => {
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [hasFavoriteCategories, setHasFavoriteCategories] = useState<boolean>(false);

  useEffect(() => {
    const authUser = async () => {
      try {
        setLoading(true);
        setError(null);
        const has_favorite_categories = await getAuth();
        setHasFavoriteCategories(has_favorite_categories);
      } catch (err: any) {
        setError(err.message || "Failed to auth user");
      } finally {
        setLoading(false);
      }
    };

    authUser();
  }, []);

  return [loading, error, hasFavoriteCategories];
};
