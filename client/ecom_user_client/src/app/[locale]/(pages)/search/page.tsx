"use client"

import API from "@/assets/configs/api";
import { handleProductImg } from "@/assets/configs/handle_img";
import * as request from "@/assets/helpers/request_without_token";
import { MetaType, ParamType } from "@/assets/types/request";
import C_ProductSimple from "@/resources/components_thuongdung/product";
import productSimple from "@/resources/components_thuongdung/product";
import { ProductSummary } from "@/types/product.types";
import { useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { Search, StarIcon } from "lucide-react";
import { useTranslations } from "next-intl";
import { useSearchParams } from "next/navigation";
import React, { useRef, useState, useEffect } from "react";

export default function SearchPage() {
  const searchParams = useSearchParams();
  const categoryId = searchParams.get('id') ;
  const [meta, setMeta] = useState<MetaType>(request.defaultMeta);
  
  const paramsRef = useRef<ParamType>({
    page: meta.currentPage,
    limit: 60,
    orderBy: 'id',
    orderDirection: 'ASC',
  });
  const t = useTranslations("System")


  const CategoryQuery = useQuery<Category, AxiosError<ResponseType>>({
    refetchOnWindowFocus: false,
    queryKey: ['search-category', categoryId],
    queryFn: async () => {
        // console.log(typeof(categoryId))
        const cate = categoryId!=null ? `${categoryId}` : ""
        const response: any = await request.get<Category>(`${API.category.getAll}`, {
         params: {
            cate_id:cate
        }
      });
      let responseData = response.data?.result.categories[0] ?? [];
      return responseData || [];
    },
  });


  const SearchQuery = useQuery<ProductSummary[], AxiosError<ResponseType>>({
    refetchOnWindowFocus: false,
    queryKey: ['search-products', categoryId],
    queryFn: async () => {
        // console.log(typeof(categoryId))
        const cate = categoryId!=null ? `category_id=${categoryId}` : ""
        const response: any = await request.get<ProductSummary[]>(`${API.product.getAll}`, {
        params: {
          ...paramsRef.current,
          search: cate
        }
      });
      let responseData = response.data?.result.products ?? [];

      if (response.data?.result.currentPage && response.data?.result.totalPages) {
        setMeta({
          currentPage: response.data?.result.currentPage,
          hasNextPage: response.data?.result.currentPage + 1 === response.data?.result.totalPages ? false : true,
          hasPreviousPage: response.data?.result.currentPage - 1 === 0 ? false : true,
          limit: paramsRef.current.limit,
          totalPages: response.data?.result.totalPages,
        });
      }
      return responseData || [];
    },
  });

  return (
    <div className="min-h-screen p-8 pb-20 font-[family-name:var(--font-geist-sans)]">
      <div className="max-w-7xl mx-auto">
        {/* Category Title */}
        {categoryId && (
          <div className="mb-6">
            <h1 className="text-2xl font-bold text-gray-800">
            {t("danh_muc_san_pham")} : {CategoryQuery.isFetched?CategoryQuery.data?.name:t("khong_tim_thay_danh_muc")}
            </h1>
          </div>
        )}

        {/* Search Results */}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
          {SearchQuery.data?.map((product) => (
          <C_ProductSimple product={product}/>
          ))}
        </div>

        {/* Loading State */}
        {SearchQuery.isLoading && (
          <div className="flex justify-center items-center py-8">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#ee4d2d]"></div>
          </div>
        )}

        {/* No Results */}
        {!SearchQuery.isLoading && SearchQuery.data?.length === 0 && (
          <div className="text-center py-8">
            <p className="text-gray-500 text-lg">{t("khong_tim_thay_san_pham_phu_hop")}</p>
          </div>
        )}
      </div>
    </div>
  );
} 