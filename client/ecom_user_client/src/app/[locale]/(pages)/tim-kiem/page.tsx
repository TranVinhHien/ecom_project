"use client";

import { useSearchParams } from "next/navigation";
import { useRouter } from "@/i18n/routing"

import { useState, useEffect, useMemo } from "react";
import { useGetActiveBanners, useGetProducts } from "@/services/apiService";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { SlidersHorizontal } from "lucide-react";
import {getImageUrl} from "@/assets/helpers/convert_tool";
import C_ProductSimple from "@/resources/components_thuongdung/product";
import { UserProfile } from "@/types/user.types";
import { INFO_USER } from "@/assets/configs/request";

export default function SearchPage() {
  const searchParams = useSearchParams();
  const router = useRouter();

  // Lấy params từ URL
  const cate_path = searchParams.get("cate_path") || undefined;
  const keywords = searchParams.get("keywords") || undefined;
  const brandParam = searchParams.get("brand") || undefined;
  const shopParam = searchParams.get("shop_id") || undefined;

  // State cho filters và pagination
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [sortBy, setSortBy] = useState<string>("default");
  const [priceMin, setPriceMin] = useState<number | undefined>(undefined);
  const [priceMax, setPriceMax] = useState<number | undefined>(undefined);
  const [brand, setBrand] = useState<string | undefined>(brandParam);
  const [showFilters, setShowFilters] = useState(false);
  // const [profile, setProfile] = useState<UserProfile   | null>(null);
 

  // Reset page về 1 khi thay đổi filters
  useEffect(() => {
    setPage(1);
  }, [sortBy, priceMin, priceMax, brand, keywords, cate_path]);

  const [profile, setProfile] = useState<UserProfile   | null>(null);
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
  // Gọi API
  const { data, isLoading, error } = useGetProducts({
    page,
    limit,
    keywords,
    cate_path,
    brand,
    shop_id: shopParam,
    price_min: priceMin,
    price_max: priceMax,
    sort: sortBy === "default" ? undefined : (sortBy as any),
  });
  const { data: categoryBanners } = useGetActiveBanners("CATEGORY");
  const normalizedCatePath = cate_path ? decodeURIComponent(cate_path) : undefined;
  const matchedCategoryBanner = useMemo(() => {
    console.log('categoryBanners', categoryBanners);
    if (!categoryBanners || !normalizedCatePath) return null;

    const matched = categoryBanners
      .map((banner) => ({
        banner,
        meta: getBannerTargetMeta(banner.bannerUrl),
      }))
      .filter(({ meta }) => meta && meta.catePath === normalizedCatePath)
      .sort((a, b) => (a.banner.bannerOrder || 0) - (b.banner.bannerOrder || 0));
    console.log(matched);
    return matched[0] || null;
  }, [categoryBanners, normalizedCatePath]);

  const handleBannerNavigation = (meta: BannerTargetMeta) => {
    const href = resolveBannerHref(meta);
    if (!href) return;
    if (ABSOLUTE_URL_PATTERN.test(href)) {
      window.open(href, "_blank", "noopener noreferrer");
    } else {
      router.push(href);
    }
  };


  // Render skeleton khi loading
  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {[...Array(10)].map((_, i) => (
            <Card key={i} className="overflow-hidden">
              <Skeleton className="h-48 w-full" />
              <CardContent className="p-4">
                <Skeleton className="h-4 w-full mb-2" />
                <Skeleton className="h-4 w-2/3" />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  // Render error
  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <p className="text-red-500">Đã có lỗi xảy ra: {error.message}</p>
        </div>
      </div>
    );
  }
  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      
      {matchedCategoryBanner?.banner && matchedCategoryBanner.meta && (
        <div className="mb-6">
          <button
            type="button"
            onClick={() => handleBannerNavigation(matchedCategoryBanner.meta!)}
            className="block w-full"
          >
            <div className="relative w-full h-40 md:h-56 lg:h-64 overflow-hidden rounded-xl shadow">

  <img
    src={getImageUrl(matchedCategoryBanner.banner.bannerImage)}
    alt={matchedCategoryBanner.banner.bannerName}
    
    // SỬA: Thêm "block" và "mx-auto". Bỏ "center"
    className="object-cover w-[70vw] block mx-auto h-[250px] md:h-[400px] lg:h-[500px] xl:h-[650px]"
    
    sizes="70vw"
/>
            </div>
          </button>
        </div>
      )}
      <div className="mb-6">
        <h1 className="text-2xl font-bold mb-2">
          {keywords ? `Kết quả tìm kiếm cho "${keywords}"` : "Danh sách sản phẩm"}
        </h1>
        <p className="text-gray-600">
          Tìm thấy {data?.totalElements || 0} sản phẩm
        </p>
      </div>

      {/* Filters Bar */}
      <div className="mb-6">
        <div className="flex flex-wrap items-center gap-4 mb-4">
          {/* Sort */}
          <Select value={sortBy} onValueChange={setSortBy}>
            <SelectTrigger className="w-[200px]">
              <SelectValue placeholder="Sắp xếp theo" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="default">Mới nhất</SelectItem>
              <SelectItem value="price_asc">Giá tăng dần</SelectItem>
              <SelectItem value="price_desc">Giá giảm dần</SelectItem>
              <SelectItem value="name_asc">Tên A-Z</SelectItem>
              <SelectItem value="name_desc">Tên Z-A</SelectItem>
              <SelectItem value="best_sell">Bán chạy nhất</SelectItem>
            </SelectContent>
          </Select>

          {/* Limit */}
          <Select value={limit.toString()} onValueChange={(val) => setLimit(Number(val))}>
            <SelectTrigger className="w-[150px]">
              <SelectValue placeholder="Hiển thị" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="20">20 sản phẩm</SelectItem>
              <SelectItem value="40">40 sản phẩm</SelectItem>
              <SelectItem value="60">60 sản phẩm</SelectItem>
            </SelectContent>
          </Select>

          {/* Toggle Filters */}
          <Button
            variant="outline"
            onClick={() => setShowFilters(!showFilters)}
            className="gap-2"
          >
            <SlidersHorizontal className="w-4 h-4" />
            Bộ lọc {showFilters ? "▲" : "▼"}
          </Button>
        </div>

        {/* Advanced Filters */}
        {showFilters && (
          <Card className="p-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {/* Price Range */}
              <div className="space-y-2">
                <label className="text-sm font-medium">Khoảng giá</label>
                <div className="flex gap-2 items-center">
                  <Input
                    type="number"
                    placeholder="Từ"
                    value={priceMin || ""}
                    onChange={(e) => setPriceMin(e.target.value ? Number(e.target.value) : undefined)}
                    className="w-full"
                  />
                  <span>-</span>
                  <Input
                    type="number"
                    placeholder="Đến"
                    value={priceMax || ""}
                    onChange={(e) => setPriceMax(e.target.value ? Number(e.target.value) : undefined)}
                    className="w-full"
                  />
                </div>
              </div>

              {/* Brand */}
              <div className="space-y-2">
                <label className="text-sm font-medium">Thương hiệu</label>
                <Input
                  placeholder="Nhập thương hiệu..."
                  value={brand || ""}
                  onChange={(e) => setBrand(e.target.value || undefined)}
                />
              </div>

              {/* Clear Filters */}
              <div className="space-y-2 flex items-end">
                <Button
                  variant="outline"
                  onClick={() => {
                    setPriceMin(undefined);
                    setPriceMax(undefined);
                    setBrand(undefined);
                    setSortBy("default");
                  }}
                  className="w-full"
                >
                  Xóa bộ lọc
                </Button>
              </div>
            </div>
          </Card>
        )}
      </div>

      {/* Products Grid */}
      {data?.data && data.data.length > 0 ? (
        <>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 mb-8">
            {data.data.map((product) => (
              <C_ProductSimple key={product.id} product={product} collection_type="search" user_id={profile?.userId} />
            ))}
          </div>

          {/* Pagination */}
          {data.totalPages > 1 && (
            <div className="flex flex-col items-center gap-4">
              <div className="text-sm text-gray-600">
                Trang {data.currentPage} / {data.totalPages} (Tổng {data.totalElements} sản phẩm)
              </div>
              
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  disabled={page === 1}
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                >
                  ← Trang trước
                </Button>

                <div className="flex items-center gap-2">
                  {/* Show pagination numbers */}
                  {(() => {
                    const maxPagesToShow = 5;
                    const totalPages = data.totalPages;
                    const currentPage = data.currentPage;
                    
                    let startPage = Math.max(1, currentPage - Math.floor(maxPagesToShow / 2));
                    let endPage = Math.min(totalPages, startPage + maxPagesToShow - 1);
                    
                    if (endPage - startPage + 1 < maxPagesToShow) {
                      startPage = Math.max(1, endPage - maxPagesToShow + 1);
                    }
                    
                    const pages = [];
                    for (let i = startPage; i <= endPage; i++) {
                      pages.push(i);
                    }
                    
                    return pages.map((pageNum) => (
                      <Button
                        key={pageNum}
                        variant={page === pageNum ? "default" : "outline"}
                        onClick={() => setPage(pageNum)}
                        className={
                          page === pageNum
                            ? "bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)]"
                            : ""
                        }
                      >
                        {pageNum}
                      </Button>
                    ));
                  })()}
                </div>

                <Button
                  variant="outline"
                  disabled={page === data.totalPages}
                  onClick={() => setPage((p) => Math.min(data.totalPages, p + 1))}
                >
                  Trang sau →
                </Button>
              </div>
            </div>
          )}
        </>
      ) : (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg">Không tìm thấy sản phẩm nào</p>
        </div>
      )}
    </div>
  );
}

const TIM_KIEM_PREFIX = "timkiem-";
const ABSOLUTE_URL_PATTERN = /^https?:\/\//i;
const FALLBACK_ORIGIN = "http://localhost:3000";

interface BannerTargetMeta {
  normalizedUrl: string;
  catePath?: string;
}

const getBannerTargetMeta = (rawUrl?: string | null): BannerTargetMeta | null => {
  if (!rawUrl) return null;
  const trimmed = rawUrl.trim();
  if (!trimmed) return null;

  const lower = trimmed.toLowerCase();
  const hasPrefix = lower.startsWith(TIM_KIEM_PREFIX);
  const normalizedUrl = hasPrefix ? trimmed.slice(TIM_KIEM_PREFIX.length) : trimmed;
  const normalizedLower = normalizedUrl.toLowerCase();

  if (!hasPrefix && !normalizedLower.startsWith("/tim-kiem") && !normalizedLower.startsWith("tim-kiem")) {
    return null;
  }

  try {
    const urlObj = buildUrlFromMaybeRelative(normalizedUrl);
    const catePath = urlObj.searchParams.get("cate_path");

    return {
      normalizedUrl,
      catePath: catePath ? decodeURIComponent(catePath) : undefined,
    };
  } catch {
    return null;
  }
};

const buildUrlFromMaybeRelative = (value: string) => {
  if (ABSOLUTE_URL_PATTERN.test(value)) {
    return new URL(value);
  }
  const absolutePath = value.startsWith("/") ? value : `/${value}`;
  return new URL(absolutePath, FALLBACK_ORIGIN);
};

const resolveBannerHref = (meta?: BannerTargetMeta | null) => {
  if (!meta || !meta.normalizedUrl) return undefined;
  if (ABSOLUTE_URL_PATTERN.test(meta.normalizedUrl)) {
    return meta.normalizedUrl;
  }
  return meta.normalizedUrl.startsWith("/") ? meta.normalizedUrl : `/${meta.normalizedUrl}`;
};
