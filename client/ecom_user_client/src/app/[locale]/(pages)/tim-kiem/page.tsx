"use client";

import { useSearchParams } from "next/navigation";
import { useRouter } from "@/i18n/routing"

import { useState, useEffect } from "react";
import { useGetProducts } from "@/services/apiService";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { ShoppingCart, Heart, SlidersHorizontal } from "lucide-react";
import { Link } from '@/i18n/routing';
import Image from "next/image";
import { useCartStore } from "@/store/cartStore";
import { useToast } from "@/hooks/use-toast";
import ROUTER from "@/assets/configs/routers";

export default function SearchPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const { toast } = useToast();
  const addToCart = useCartStore((state) => state.addToCart);

  // Helper để xử lý image URL
  const getImageUrl = (imageUrl: string | null | undefined) => {
    if (!imageUrl) return "/placeholder.png";
    // Nếu là URL đầy đủ (http/https) thì giữ nguyên
    if (imageUrl.startsWith("http://") || imageUrl.startsWith("https://")) {
      return imageUrl;
    }
    // Nếu không có protocol, thêm http:// vào đầu
    return `http://${imageUrl}`;
  };

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

  // Reset page về 1 khi thay đổi filters
  useEffect(() => {
    setPage(1);
  }, [sortBy, priceMin, priceMax, brand, keywords, cate_path]);

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

  // Format giá
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
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
              <SelectItem value="default">Mặc định</SelectItem>
              <SelectItem value="price_asc">Giá tăng dần</SelectItem>
              <SelectItem value="price_desc">Giá giảm dần</SelectItem>
              <SelectItem value="newest">Mới nhất</SelectItem>
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
              <Card
                key={product.id}
                className="overflow-hidden hover:shadow-lg transition-shadow duration-300 group flex flex-col h-full"
              >
                <Link href={`/product/${product.key}`}>
                  <div className="relative aspect-square overflow-hidden bg-gray-100">
                    <Image
                      src={getImageUrl(product.image)}
                      alt={product.name}
                      fill
                      className="object-cover group-hover:scale-110 transition-transform duration-300"
                      unoptimized
                    />
                    {/* Wishlist button */}
                    <button
                      className="absolute top-2 right-2 p-2 bg-white rounded-full shadow-md hover:bg-gray-100 transition-colors"
                      onClick={(e) => {
                        e.preventDefault();
                        toast({
                          title: "Đã thêm vào yêu thích",
                          description: product.name,
                        });
                      }}
                    >
                      <Heart className="w-4 h-4 text-gray-600" />
                    </button>
                  </div>
                </Link>

                <CardContent className="p-4 flex-1 flex flex-col">
                  <Link href={`${ROUTER.product}/${product.key}`}>
                    <h3 className="font-medium text-sm mb-2 line-clamp-2 hover:text-[hsl(var(--primary))] min-h-[40px]">
                      {product.name}
                    </h3>
                  </Link>

                  <div className="flex items-baseline gap-2 mt-auto">
                    {product.min_price === product.max_price ? (
                      <span className="text-[hsl(var(--primary))] font-bold text-lg">
                        {formatPrice(product.min_price)}
                      </span>
                    ) : (
                      <>
                        <span className="text-[hsl(var(--primary))] font-bold text-lg">
                          {formatPrice(product.min_price)}
                        </span>
                        <span className="text-gray-400 text-sm">
                          - {formatPrice(product.max_price)}
                        </span>
                      </>
                    )}
                  </div>
                </CardContent>

                <CardFooter className="p-4 pt-0">
                  <Link href={`${ROUTER.product}/${product.key}`} className="w-full">
                    <Button className="w-full bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)]">
                      {/* <ShoppingCart className="w-4 h-4 mr-2" /> */}
                      Xem chi tiết
                    </Button>
                  </Link>
                </CardFooter>
              </Card>
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
