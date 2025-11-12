"use client"

import C_ProductSimple from "@/resources/components_thuongdung/product";
import { useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { useTranslations } from "next-intl";
import { useSearchParams } from "next/navigation";
import React from "react";
import apiClient from "@/lib/apiClient";
import { ProductSummary } from "@/types/product.types";

export default function SearchPage() {
  const searchParams = useSearchParams();
  const query = searchParams.get('query');
  const t = useTranslations("System")

  const SearchQuery = useQuery<ProductSummary[], AxiosError<any>>({
    refetchOnWindowFocus: false,
    queryKey: ['semantic-search-products', query],
    enabled: !!query, // tránh gọi khi query null hoặc rỗng
    queryFn: async () => {
      const payload: any = {
        query: query ?? '',
        page: 1,
        top_k: 50,
      };
      const response = await apiClient.post('/products/semantic-search', payload);
      const data: any = response.data;
      const products = data.products ?? [];
      
      // Map response to ProductSummary format
      const mappedProducts: ProductSummary[] = products.map((product: any) => ({
        id: product.id || product.product_id,
        key: product.key || '',
        name: product.name || '',
        shop_id: product.shop_id || '',
        category_id: product.category_id || '',
        brand_id: product.brand_id || null,
        min_price: product.min_price || 0,
        max_price: product.max_price || 0,
        min_price_sku_id: product.min_price_sku_id || null,
        max_price_sku_id: product.max_price_sku_id || null,
        media: product.media || null,
        description: product.description || null,
        short_description: product.short_description || null,
        can_view: product.can_view ?? true,
        can_buy: product.can_buy ?? true,
        can_comment: product.can_comment ?? true,
        delete_status: product.delete_status || 0,
        created_at: product.created_at || '',
        updated_at: product.updated_at || '',
      }));
      
      return mappedProducts;
    },
  });
  
  if (query === "" || !query) {
    return <div className="min-h-screen p-8">
      <div className="max-w-7xl mx-auto">
        <p className="text-gray-500 text-lg text-center">{t("vui_long_nhap_tu_khoa_tim_kiem")}</p>
      </div>
    </div>;
  }
  
  return (
    <div className="min-h-screen p-8 pb-20 font-[family-name:var(--font-geist-sans)]">
      <div className="max-w-7xl mx-auto">
        {/* Category Title */}
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-800">
            {!SearchQuery.isLoading ? t("ket_qua_tim_kiem_cho") + " " + query : t("dang_tim_kiem")}
          </h1>
        </div>

        {/* Search Results */}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
          {SearchQuery.data?.map((product) => (
            <C_ProductSimple key={product.id} product={product} />
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