"use client"

import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useRouter } from "@/i18n/routing";
import { Button } from '@/components/ui/button';
import { ShopApiResponse } from '@/types/shop.types';
import { PaginatedProductsResponse, ProductListParams } from '@/types/product.types';
import { Loading } from '@/components/ui/loading';
import { Store, CheckCircle, MapPin, Star, Heart, MessageCircle, Phone, Mail, Calendar, Award, Users, TrendingUp, Search, SlidersHorizontal, X } from 'lucide-react';
import apiClient from '@/lib/apiClient';
import API from '@/assets/configs/api';
import { getImageUrl } from '@/assets/helpers/convert_tool';
import { Card } from '@/components/ui/card';
import { useGetProducts } from '@/services/apiService';
import C_ProductSimple from '@/resources/components_thuongdung/product';
import { Input } from '@/components/ui/input';
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select';

export default function ShopPage({ params }: { params: { id: string } }) {
  const t = useTranslations("System");
  const router = useRouter();
  const [isFollowing, setIsFollowing] = useState(false);
  
  // States for product filtering and search
  const [searchKeyword, setSearchKeyword] = useState('');
  const [sortBy, setSortBy] = useState<ProductListParams['sort']>('best_sell');
  const [priceMin, setPriceMin] = useState<number | undefined>();
  const [priceMax, setPriceMax] = useState<number | undefined>();
  const [showFilters, setShowFilters] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);

  const { data, isLoading, error } = useQuery<ShopApiResponse>({
    queryKey: ['shop-detail', params.id],
    queryFn: async () => {
      const response = await apiClient.get(`/Shops/${params.id}`, {
        customBaseURL: 'http://localhost:8000/api'
      });
      return response.data;
    },
  });

  // Fetch best selling products
  const { data: bestSellingProducts, isLoading: isLoadingBestSelling } = useGetProducts({
    page: 1,
    limit: 10,
    shop_id: params.id,
    sort: 'best_sell',
  });

  // Fetch all products with filters
  const { data: allProducts, isLoading: isLoadingAllProducts } = useGetProducts({
    page: currentPage,
    limit: 20,
    shop_id: params.id,
    keywords: searchKeyword || undefined,
    sort: sortBy,
    price_min: priceMin,
    price_max: priceMax,
  });


  const handleFollow = () => {
    setIsFollowing(!isFollowing);
    // TODO: Call API to follow/unfollow shop
  };

  const handleProductClick = (productKey: string) => {
    router.push(`/product/${productKey}`);
  };

  const handleSearch = () => {
    setCurrentPage(1); // Reset về trang 1 khi search
  };

  const handleApplyFilters = () => {
    setCurrentPage(1); // Reset về trang 1 khi apply filters
    setShowFilters(false);
  };

  const handleClearFilters = () => {
    setPriceMin(undefined);
    setPriceMax(undefined);
    setSortBy('best_sell');
    setSearchKeyword('');
    setCurrentPage(1);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loading size="lg" variant="primary" />
      </div>
    );
  }

  if (error || !data?.result) {
    return (
      <div className="flex items-center justify-center min-h-[400px] text-red-500">
        {t("co_loi_xay_ra_khi_tai_du_lieu")}
      </div>
    );
  }

  const shop = data.result;

  return (
    <div className="max-w-7xl mx-auto py-8 px-4">
      {/* Shop Header Section */}
      <div className="bg-gradient-to-r from-blue-600 to-indigo-600 rounded-lg p-8 mb-6 text-white shadow-lg">
        <div className="flex items-start gap-6">
          {/* Shop Logo */}
          <div className="relative w-32 h-32 rounded-full overflow-hidden border-4 border-white shadow-xl flex-shrink-0 bg-white">
            <img 
              src={shop.shopLogo} 
              alt={shop.shopName}
              className="w-full h-full object-cover"
              onError={(e) => {
                (e.target as HTMLImageElement).src = '/placeholder-shop.png';
              }}
            />
            {shop.taxInfo.taxActiveStatus && (
              <div className="absolute bottom-2 right-2 bg-blue-600 rounded-full p-2 shadow-md">
                <CheckCircle className="w-5 h-5 text-white" />
              </div>
            )}
          </div>

          {/* Shop Info */}
          <div className="flex-1">
            <div className="flex items-start justify-between mb-4">
              <div>
                <div className="flex items-center gap-3 mb-2">
                  <h1 className="text-3xl font-bold">{shop.shopName}</h1>
                  {shop.taxInfo.taxActiveStatus && (
                    <span className="bg-white text-blue-600 text-sm px-3 py-1 rounded-full flex items-center gap-1 font-semibold">
                      <Store className="w-4 h-4" />
                      OFFICIAL
                    </span>
                  )}
                </div>
                
                {shop.shopDescription && (
                  <p className="text-blue-50 mb-3 max-w-2xl">
                    {shop.shopDescription}
                  </p>
                )}
              </div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-4 gap-4 mb-4">
              <div className="bg-white/10 backdrop-blur-sm rounded-lg p-3">
                <div className="flex items-center gap-2 mb-1">
                  <Star className="w-5 h-5 text-yellow-300 fill-yellow-300" />
                  <span className="text-2xl font-bold">4.8</span>
                </div>
                <p className="text-sm text-blue-100">Đánh giá</p>
              </div>

              <div className="bg-white/10 backdrop-blur-sm rounded-lg p-3">
                <div className="flex items-center gap-2 mb-1">
                  <Users className="w-5 h-5" />
                  <span className="text-2xl font-bold">{shop.followerCount.toLocaleString()}</span>
                </div>
                <p className="text-sm text-blue-100">Người theo dõi</p>
              </div>

              <div className="bg-white/10 backdrop-blur-sm rounded-lg p-3">
                <div className="flex items-center gap-2 mb-1">
                  <Award className="w-5 h-5" />
                  <span className="text-2xl font-bold">98%</span>
                </div>
                <p className="text-sm text-blue-100">Phản hồi chat</p>
              </div>

              <div className="bg-white/10 backdrop-blur-sm rounded-lg p-3">
                <div className="flex items-center gap-2 mb-1">
                  <Store className="w-5 h-5" />
                  <span className="text-2xl font-bold">10000+</span>
                </div>
                <p className="text-sm text-blue-100">Sản phẩm</p>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex gap-3">
              <Button
                variant={isFollowing ? "secondary" : "default"}
                onClick={handleFollow}
                className={isFollowing ? "bg-white text-blue-600 hover:bg-gray-100" : "bg-white text-blue-600 hover:bg-gray-100"}
              >
                <Heart className={`w-4 h-4 mr-2 ${isFollowing ? 'fill-red-500 text-red-500' : ''}`} />
                {isFollowing ? 'Đã theo dõi' : 'Theo dõi'}
              </Button>
              <Button
                variant="secondary"
                className="bg-white/10 border-white text-white hover:bg-white/20"
              >
                <MessageCircle className="w-4 h-4 mr-2" />
                Chat ngay
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Shop Details Section */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        {/* Shop Information Card */}
        <div className="lg:col-span-2 bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
            <Store className="w-5 h-5 text-blue-600" />
            Thông tin cửa hàng
          </h2>
          
          <div className="space-y-4">
            {shop.shopAddress && (
              <div className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
                <MapPin className="w-5 h-5 text-blue-600 flex-shrink-0 mt-1" />
                <div>
                  <p className="font-semibold text-sm text-gray-600 mb-1">Địa chỉ</p>
                  <p className="text-gray-800">{shop.shopAddress}</p>
                </div>
              </div>
            )}

            {shop.shopPhone && (
              <div className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
                <Phone className="w-5 h-5 text-blue-600 flex-shrink-0 mt-1" />
                <div>
                  <p className="font-semibold text-sm text-gray-600 mb-1">Số điện thoại</p>
                  <p className="text-gray-800">{shop.shopPhone}</p>
                </div>
              </div>
            )}

            {shop.shopEmail && (
              <div className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
                <Mail className="w-5 h-5 text-blue-600 flex-shrink-0 mt-1" />
                <div>
                  <p className="font-semibold text-sm text-gray-600 mb-1">Email</p>
                  <p className="text-gray-800">{shop.shopEmail}</p>
                </div>
              </div>
            )}

            {shop.createdDate && (
              <div className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
                <Calendar className="w-5 h-5 text-blue-600 flex-shrink-0 mt-1" />
                <div>
                  <p className="font-semibold text-sm text-gray-600 mb-1">Tham gia</p>
                  <p className="text-gray-800">
                    {new Date(shop.createdDate).toLocaleDateString('vi-VN', { 
                      year: 'numeric', 
                      month: 'long', 
                      day: 'numeric' 
                    })}
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Business Information Card */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
            <Award className="w-5 h-5 text-blue-600" />
            Thông tin doanh nghiệp
          </h2>
          
          <div className="space-y-3">
            {shop.taxInfo.taxNationalName && (
              <div className="pb-3 border-b">
                <p className="text-sm text-gray-600 mb-1">Tên doanh nghiệp</p>
                <p className="font-semibold text-gray-800">{shop.taxInfo.taxNationalName}</p>
              </div>
            )}

            {shop.taxInfo.taxCode && (
              <div className="pb-3 border-b">
                <p className="text-sm text-gray-600 mb-1">Mã số thuế</p>
                <p className="font-semibold text-gray-800">{shop.taxInfo.taxCode}</p>
              </div>
            )}

            {shop.taxInfo.taxPresentName && (
              <div className="pb-3 border-b">
                <p className="text-sm text-gray-600 mb-1">Người đại diện</p>
                <p className="font-semibold text-gray-800">{shop.taxInfo.taxPresentName}</p>
              </div>
            )}

            {shop.taxInfo.taxBusinessType && (
              <div className="pb-3 border-b">
                <p className="text-sm text-gray-600 mb-1">Loại hình kinh doanh</p>
                <p className="font-semibold text-gray-800">{shop.taxInfo.taxBusinessType}</p>
              </div>
            )}

            {shop.taxInfo.taxActiveDate && (
              <div>
                <p className="text-sm text-gray-600 mb-1">Ngày hoạt động</p>
                <p className="font-semibold text-gray-800">
                  {new Date(shop.taxInfo.taxActiveDate).toLocaleDateString('vi-VN', {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric'
                  })}
                </p>
              </div>
            )}

            {shop.taxInfo.taxActiveStatus && (
              <div className="pt-3 flex items-center gap-2 text-green-600">
                <CheckCircle className="w-5 h-5" />
                <span className="font-semibold">Đã xác thực</span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Products Section - Placeholder */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <TrendingUp className="w-6 h-6 text-blue-600" />
            Sản phẩm bán chạy
          </h2>
          {bestSellingProducts?.totalElements && (
            <span className="text-sm text-gray-500">
              {bestSellingProducts.totalElements} sản phẩm
            </span>
          )}
        </div>

        {isLoadingBestSelling ? (
          <div className="flex items-center justify-center py-12">
            <Loading size="lg" variant="primary" />
          </div>
        ) : bestSellingProducts?.data && bestSellingProducts.data.length > 0 ? (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
            {bestSellingProducts.data.map((product) => (
            <C_ProductSimple key={product.id} product={product} />
            ))}
          </div>
        ) : (
          <div className="text-center text-gray-500 py-12">
            <Store className="w-16 h-16 mx-auto mb-4 text-gray-300" />
            <p className="text-lg">Chưa có sản phẩm nào</p>
          </div>
        )}
      </div>

      {/* All Products Section with Search and Filters */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
          <Store className="w-6 h-6 text-blue-600" />
          Tất cả sản phẩm
        </h2>

        {/* Search and Filter Bar */}
        <div className="mb-6 space-y-4">
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
          {(searchKeyword || priceMin || priceMax) && (
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
        </div>

        {/* Products Grid */}
        {isLoadingAllProducts ? (
          <div className="flex items-center justify-center py-12">
            <Loading size="lg" variant="primary" />
          </div>
        ) : allProducts?.data && allProducts.data.length > 0 ? (
          <>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4 mb-6">
              {allProducts.data.map((product) => (
                <C_ProductSimple key={product.id} product={product} />
              ))}
            </div>

            {/* Pagination */}
            {allProducts.totalPages > 1 && (
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
                  {Array.from({ length: Math.min(5, allProducts.totalPages) }, (_, i) => {
                    let pageNum;
                    if (allProducts.totalPages <= 5) {
                      pageNum = i + 1;
                    } else if (currentPage <= 3) {
                      pageNum = i + 1;
                    } else if (currentPage >= allProducts.totalPages - 2) {
                      pageNum = allProducts.totalPages - 4 + i;
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
                  onClick={() => setCurrentPage(prev => Math.min(allProducts.totalPages, prev + 1))}
                  disabled={currentPage === allProducts.totalPages}
                >
                  Trang sau
                </Button>
              </div>
            )}

            {/* Results Info */}
            <div className="text-center text-sm text-gray-500 mt-4">
              Hiển thị {allProducts.data.length} / {allProducts.totalElements} sản phẩm
              (Trang {currentPage} / {allProducts.totalPages})
            </div>
          </>
        ) : (
          <div className="text-center text-gray-500 py-12">
            <Store className="w-16 h-16 mx-auto mb-4 text-gray-300" />
            <p className="text-lg">Không tìm thấy sản phẩm phù hợp</p>
            <Button
              variant="outline"
              size="sm"
              onClick={handleClearFilters}
              className="mt-4"
            >
              Xóa bộ lọc
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
