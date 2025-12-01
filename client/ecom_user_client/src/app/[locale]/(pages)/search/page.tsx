"use client"

import API from "@/assets/configs/api";
import { handleProductImg } from "@/assets/configs/handle_img";
import * as request from "@/assets/helpers/request_without_token";
import { MetaType, ParamType } from "@/assets/types/request";
import C_ProductSimple from "@/resources/components_thuongdung/product";
import productSimple from "@/resources/components_thuongdung/product";
import { ProductSummary, ProductListParams } from "@/types/product.types";
import { useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { Search, StarIcon, SlidersHorizontal, X } from "lucide-react";
import { useTranslations } from "next-intl";
import { useSearchParams } from "next/navigation";
import React, { useRef, useState, useEffect } from "react";
import { useGetProducts, useGetCategories } from "@/services/apiService";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from "@/components/ui/select";
import { Loading } from "@/components/ui/loading";
import { INFO_USER } from "@/assets/configs/request";
import { UserProfile } from "@/types/user.types";

export default function SearchPage() {
  const searchParams = useSearchParams();
  const categoryId = searchParams.get('id');
  const categoryPath = searchParams.get('path'); // Lấy category path từ URL
  const initialKeywords = searchParams.get('keywords') || ''; // Lấy keywords từ URL
  
  const t = useTranslations("System");

  // States for filtering and search
  const [searchKeyword, setSearchKeyword] = useState(initialKeywords);
  const [sortBy, setSortBy] = useState<ProductListParams['sort']>('best_sell');
  const [priceMin, setPriceMin] = useState<number | undefined>();
  const [priceMax, setPriceMax] = useState<number | undefined>();
  const [showFilters, setShowFilters] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [profile, setProfile] = useState<UserProfile | null>(null);

  // Reset page when filters change
  useEffect(() => {
    setCurrentPage(1);
  }, [searchKeyword, sortBy, priceMin, priceMax, categoryPath]);
  // Reset page when filters change
  useEffect(() => {
    setCurrentPage(1);
  }, [searchKeyword, sortBy, priceMin, priceMax, categoryPath]);

  // Fetch category info (giữ nguyên logic cũ)
  const CategoryQuery = useQuery<Category, AxiosError<ResponseType>>({
    refetchOnWindowFocus: false,
    queryKey: ['search-category', categoryId],
    queryFn: async () => {
      const cate = categoryId != null ? `${categoryId}` : "";
      const response: any = await request.get<Category>(`${API.category.getAll}`, {
        params: {
          cate_id: cate
        }
      });
      let responseData = response.data?.result.categories[0] ?? [];
      return responseData || [];
    },
    enabled: !!categoryId, // Chỉ gọi khi có categoryId
  });

  // Fetch products using useGetProducts hook
  const { data: productsData, isLoading, error } = useGetProducts({
    page: currentPage,
    limit: 20,
    keywords: searchKeyword || undefined,
    cate_path: categoryPath || undefined,
    sort: sortBy,
    price_min: priceMin,
    price_max: priceMax,
  });

  const handleSearch = () => {
    setCurrentPage(1);
  };

  const handleApplyFilters = () => {
    setCurrentPage(1);
    setShowFilters(false);
  };

  useEffect(() => {
    const userInfo = localStorage.getItem(INFO_USER);
    if (userInfo) {
      try {
        const userData = JSON.parse(userInfo);
        setProfile(userData);
      }
      catch (error) {
        console.error("Error parsing user data:", error);
      }
    }
  }, []);

  const handleClearFilters = () => {
    setPriceMin(undefined);
    setPriceMax(undefined);
    setSortBy('best_sell');
    setSearchKeyword('');
    setCurrentPage(1);
  };
  
  return (
    <div className="min-h-screen p-8 pb-20 font-[family-name:var(--font-geist-sans)]">
      <div className="max-w-7xl mx-auto">
        {/* Category Title */}
        {categoryId && (
          <div className="mb-6">
            <h1 className="text-2xl font-bold text-gray-800">
              {t("danh_muc_san_pham")}: {CategoryQuery.isFetched ? CategoryQuery.data?.name : t("khong_tim_thay_danh_muc")}
            </h1>
          </div>
        )}

        {/* Search and Filter Bar */}
        <div className="mb-6 bg-white rounded-lg shadow-md p-4 space-y-4">
          {/* Search and Sort Row */}
          <div className="flex gap-3">
            {/* Search Input */}
            <div className="flex-1 relative">
              <Input
                type="text"
                placeholder={t("tim_kiem_san_pham")}
                value={searchKeyword}
                onChange={(e) => setSearchKeyword(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                className="pr-10"
              />
              <Button
                size="sm"
                onClick={handleSearch}
                className="absolute right-1 top-1/2 -translate-y-1/2 h-8"
              >
                <Search className="w-4 h-4" />
              </Button>
            </div>

            {/* Sort Select */}
            <Select value={sortBy} onValueChange={(value: any) => setSortBy(value)}>
              <SelectTrigger className="w-[200px]">
                <SelectValue placeholder="Sắp xếp" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="best_sell">Bán chạy</SelectItem>
                <SelectItem value="price_asc">Giá thấp đến cao</SelectItem>
                <SelectItem value="price_desc">Giá cao đến thấp</SelectItem>
                <SelectItem value="name_asc">Tên A-Z</SelectItem>
                <SelectItem value="name_desc">Tên Z-A</SelectItem>
              </SelectContent>
            </Select>

            {/* Filter Button */}
            <Button
              variant="outline"
              onClick={() => setShowFilters(!showFilters)}
              className="gap-2"
            >
              <SlidersHorizontal className="w-4 h-4" />
              Bộ lọc
            </Button>
          </div>

          {/* Filter Panel */}
          {showFilters && (
            <div className="p-4 border rounded-lg bg-gray-50 space-y-4">
              <div className="flex items-center justify-between mb-3">
                <h3 className="font-semibold">Bộ lọc giá</h3>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setShowFilters(false)}
                >
                  <X className="w-4 h-4" />
                </Button>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-sm text-gray-600 mb-1 block">Giá tối thiểu</label>
                  <Input
                    type="number"
                    placeholder="0"
                    value={priceMin || ''}
                    onChange={(e) => setPriceMin(e.target.value ? Number(e.target.value) : undefined)}
                  />
                </div>
                <div>
                  <label className="text-sm text-gray-600 mb-1 block">Giá tối đa</label>
                  <Input
                    type="number"
                    placeholder="1000000"
                    value={priceMax || ''}
                    onChange={(e) => setPriceMax(e.target.value ? Number(e.target.value) : undefined)}
                  />
                </div>
              </div>

              <div className="flex gap-2 justify-end">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleClearFilters}
                >
                  Xóa bộ lọc
                </Button>
                <Button
                  size="sm"
                  onClick={handleApplyFilters}
                >
                  Áp dụng
                </Button>
              </div>
            </div>
          )}

          {/* Active Filters Display */}
          {(searchKeyword || priceMin || priceMax || categoryPath) && (
            <div className="flex flex-wrap gap-2 items-center">
              <span className="text-sm text-gray-600">Bộ lọc đang áp dụng:</span>
              {searchKeyword && (
                <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full flex items-center gap-1">
                  Từ khóa: {searchKeyword}
                  <button onClick={() => setSearchKeyword('')}>
                    <X className="w-3 h-3" />
                  </button>
                </span>
              )}
              {priceMin && (
                <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full flex items-center gap-1">
                  Giá min: {priceMin.toLocaleString()}₫
                  <button onClick={() => setPriceMin(undefined)}>
                    <X className="w-3 h-3" />
                  </button>
                </span>
              )}
              {priceMax && (
                <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full flex items-center gap-1">
                  Giá max: {priceMax.toLocaleString()}₫
                  <button onClick={() => setPriceMax(undefined)}>
                    <X className="w-3 h-3" />
                  </button>
                </span>
              )}
              {categoryPath && (
                <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">
                  Danh mục: {categoryPath}
                </span>
              )}
              <Button
                variant="ghost"
                size="sm"
                onClick={handleClearFilters}
                className="text-xs h-6"
              >
                Xóa tất cả
              </Button>
            </div>
          )}

          {/* Results Info */}
          {productsData && (
            <div className="text-sm text-gray-600">
              Tìm thấy <span className="font-semibold">{productsData.totalElements}</span> sản phẩm
            </div>
          )}
        </div>

        {/* Search Results */}
        {isLoading ? (
          <div className="flex justify-center items-center py-12">
            <Loading size="lg" variant="primary" />
          </div>
        ) : productsData?.data && productsData.data.length > 0 ? (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6 mb-6">
              {productsData.data.map((product) => (
                <C_ProductSimple key={product.id} product={product} user_id={profile?.id || ""} collection_type="search" />
              ))}
            </div>

            {/* Pagination */}
            {productsData.totalPages > 1 && (
              <div className="flex items-center justify-center gap-2 mt-6">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                  disabled={currentPage === 1}
                >
                  Trang trước
                </Button>
                
                <div className="flex gap-1">
                  {Array.from({ length: Math.min(5, productsData.totalPages) }, (_, i) => {
                    let pageNum;
                    if (productsData.totalPages <= 5) {
                      pageNum = i + 1;
                    } else if (currentPage <= 3) {
                      pageNum = i + 1;
                    } else if (currentPage >= productsData.totalPages - 2) {
                      pageNum = productsData.totalPages - 4 + i;
                    } else {
                      pageNum = currentPage - 2 + i;
                    }
                    
                    return (
                      <Button
                        key={pageNum}
                        variant={currentPage === pageNum ? "default" : "outline"}
                        size="sm"
                        onClick={() => setCurrentPage(pageNum)}
                        className="w-10"
                      >
                        {pageNum}
                      </Button>
                    );
                  })}
                </div>

                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(prev => Math.min(productsData.totalPages, prev + 1))}
                  disabled={currentPage === productsData.totalPages}
                >
                  Trang sau
                </Button>
              </div>
            )}

            {/* Page Info */}
            <div className="text-center text-sm text-gray-500 mt-4">
              Trang {currentPage} / {productsData.totalPages}
            </div>
          </>
        ) : (
          <div className="text-center py-12 bg-white rounded-lg shadow-md">
            <p className="text-gray-500 text-lg mb-4">{t("khong_tim_thay_san_pham_phu_hop")}</p>
            {(searchKeyword || priceMin || priceMax) && (
              <Button
                variant="outline"
                onClick={handleClearFilters}
              >
                Xóa bộ lọc
              </Button>
            )}
          </div>
        )}
      </div>
    </div>
  );
} 