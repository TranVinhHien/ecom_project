"use client"

import React, { useState, useEffect, useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useRouter } from "@/i18n/routing";
import { useSearchParams } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { 
  Store, Heart, MessageCircle, Star, Users, Award, 
  MapPin, Phone, Mail, Calendar, CheckCircle, 
  Search, SlidersHorizontal, X, ChevronRight, ChevronLeft, Package,
  ShoppingBag, Tag, Grid, Filter
} from 'lucide-react';
import { ShopApiResponse } from '@/types/shop.types';
import { ProductListParams } from '@/types/product.types';
import { useGetProducts, useGetActiveBanners, useGetCategories } from '@/services/apiService';
import apiClient from '@/lib/apiClient';
import C_ProductSimple from '@/resources/components_thuongdung/product';
import { Loading } from '@/components/ui/loading';
import { getImageUrl } from '@/assets/helpers/convert_tool';

export default function ShopPage({ params }: { params: { id: string } }) {
  const t = useTranslations("System");
  const router = useRouter();
  const searchParams = useSearchParams();
  const [isFollowing, setIsFollowing] = useState(false);
  
  // State for banner slideshow
  const [currentBannerIndex, setCurrentBannerIndex] = useState(0);
  const [isBannerPaused, setIsBannerPaused] = useState(false);
  
  // State for home page products (separate from products tab)
  const [homeProductsPage, setHomeProductsPage] = useState(1);
  const [homeProducts, setHomeProducts] = useState<any[]>([]);
  
  // Get active tab from URL, default to 'shop'
  const activeTab = searchParams.get('tab') || 'shop';
  
  // States for product filtering and search - initialized from URL params
  const [searchKeyword, setSearchKeyword] = useState(searchParams.get('keywords') || '');
  const [sortBy, setSortBy] = useState<ProductListParams['sort']>(
    (searchParams.get('sort') as ProductListParams['sort']) || 'best_sell'
  );
  const [selectedCategory, setSelectedCategory] = useState<string>(searchParams.get('category') || '');
  const [priceMin, setPriceMin] = useState<number | undefined>(
    searchParams.get('price_min') ? Number(searchParams.get('price_min')) : undefined
  );
  const [priceMax, setPriceMax] = useState<number | undefined>(
    searchParams.get('price_max') ? Number(searchParams.get('price_max')) : undefined
  );
  const [showFilters, setShowFilters] = useState(false);
  const [currentPage, setCurrentPage] = useState(
    searchParams.get('page') ? Number(searchParams.get('page')) : 1
  );

  // Function to update URL query parameters
  const updateURLParams = useCallback((updates: Record<string, string | number | undefined | null>) => {
    const current = new URLSearchParams(searchParams.toString());
    
    Object.entries(updates).forEach(([key, value]) => {
      if (value === undefined || value === null || value === '') {
        current.delete(key);
      } else {
        current.set(key, String(value));
      }
    });
    
    const newURL = `?${current.toString()}`;
    router.push(newURL, { scroll: false });
  }, [searchParams, router]);

  // Sync state with URL params when URL changes (browser back/forward)
  useEffect(() => {
    const tab = searchParams.get('tab') || 'shop';
    const keywords = searchParams.get('keywords') || '';
    const sort = (searchParams.get('sort') as ProductListParams['sort']) || 'best_sell';
    const category = searchParams.get('category') || '';
    const priceMinParam = searchParams.get('price_min') ? Number(searchParams.get('price_min')) : undefined;
    const priceMaxParam = searchParams.get('price_max') ? Number(searchParams.get('price_max')) : undefined;
    const page = searchParams.get('page') ? Number(searchParams.get('page')) : 1;

    setSearchKeyword(keywords);
    setSortBy(sort);
    setSelectedCategory(category);
    setPriceMin(priceMinParam);
    setPriceMax(priceMaxParam);
    setCurrentPage(page);
  }, [searchParams]);

  // Fetch shop data
  const { data, isLoading, error } = useQuery<ShopApiResponse>({
    queryKey: ['shop-detail', params.id],
    queryFn: async () => {
      const response = await apiClient.get(`/Shops/${params.id}`, {
        customBaseURL: 'http://localhost:8000/api'
      });
      return response.data;
    },
  });

  // Fetch banners
  const { data: homeBanners } = useGetActiveBanners('HOME');
  const { data: categoryBanners } = useGetActiveBanners('CATEGORY');
  const { data: productBanners } = useGetActiveBanners('PRODUCT');
  const { data: promotionBanners } = useGetActiveBanners('PROMOTION');

  // Fetch categories
  const { data: categoriesData } = useGetCategories();

  // Fetch all products with filters (for products tab)
  const { data: allProducts, isLoading: isLoadingAllProducts } = useGetProducts({
    page: currentPage,
    limit: 20,
    shop_id: params.id,
    keywords: searchKeyword || undefined,
    cate_path: selectedCategory || undefined,
    sort: sortBy,
    price_min: priceMin,
    price_max: priceMax,
  });

  // Fetch products for home page (shop tab)
  const { data: homeProductsData, isLoading: isLoadingHomeProducts } = useGetProducts({
    page: homeProductsPage,
    limit: 60,
    shop_id: params.id,
  });

  const handleFollow = () => {
    setIsFollowing(!isFollowing);
    // TODO: Call API to follow/unfollow shop
  };

  const handleBannerClick = (url: string) => {
    if (url.startsWith('http')) {
      window.open(url, '_blank');
    } else {
      router.push(url);
    }
  };

  // Handler to change tab
  const handleTabChange = (tab: string) => {
    updateURLParams({ tab });
  };

  // Handler for search
  const handleSearch = () => {
    updateURLParams({ 
      keywords: searchKeyword || undefined,
      page: 1 
    });
  };

  // Handler to apply filters
  const handleApplyFilters = () => {
    updateURLParams({
      price_min: priceMin,
      price_max: priceMax,
      page: 1
    });
    setShowFilters(false);
  };

  // Handler to clear all filters
  const handleClearFilters = () => {
    updateURLParams({
      keywords: undefined,
      category: undefined,
      price_min: undefined,
      price_max: undefined,
      sort: undefined,
      page: 1
    });
  };

  // Handler when category changes
  const handleCategoryChange = (category: string) => {
    const categoryValue = category === 'all' ? undefined : category;
    updateURLParams({ 
      category: categoryValue,
      page: 1 
    });
  };

  // Handler when sort changes
  const handleSortChange = (value: string) => {
    const sort = value as ProductListParams['sort'];
    updateURLParams({ 
      sort: sort === 'best_sell' ? undefined : sort,
      page: 1 
    });
  };

  // Handler when page changes
  const handlePageChange = (page: number) => {
    updateURLParams({ page });
  };

  // Filter banners by shopId - Calculate before early returns to ensure hooks are always called
  const shopHomeBanners = homeBanners?.filter(b => b.shopId === params.id) || [];
  const shopCategoryBanners = categoryBanners?.filter(b => b.shopId === params.id) || [];
  const shopProductBanners = productBanners?.filter(b => b.shopId === params.id) || [];
  const shopPromotionBanners = promotionBanners?.filter(b => b.shopId === params.id) || [];

  // Auto slide for banner slideshow - Must be before early returns
  useEffect(() => {
    if (shopHomeBanners.length <= 1 || isBannerPaused) return;
    
    const interval = setInterval(() => {
      setCurrentBannerIndex((prev) => (prev + 1) % shopHomeBanners.length);
    }, 5000); // Change slide every 5 seconds

    return () => clearInterval(interval);
  }, [shopHomeBanners.length, isBannerPaused]);

  // Reset banner index when banners change - Must be before early returns
  useEffect(() => {
    setCurrentBannerIndex(0);
  }, [shopHomeBanners.length]);

  // Update home products when data changes
  useEffect(() => {
    if (homeProductsData?.data && homeProductsData.data.length > 0) {
      if (homeProductsPage === 1) {
        // First page: replace products
        setHomeProducts(homeProductsData.data);
      } else {
        // Subsequent pages: append products (avoid duplicates)
        setHomeProducts(prev => {
          const existingIds = new Set(prev.map(p => p.id));
          const newProducts = homeProductsData.data.filter(p => !existingIds.has(p.id));
          return [...prev, ...newProducts];
        });
      }
    }
  }, [homeProductsData, homeProductsPage]);

  // Reset home products when shop changes
  useEffect(() => {
    setHomeProducts([]);
    setHomeProductsPage(1);
  }, [params.id]);

  // Handler to load more products
  const handleLoadMoreProducts = () => {
    setHomeProductsPage(prev => prev + 1);
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

  // Banner navigation handlers
  const handleBannerNext = () => {
    setCurrentBannerIndex((prev) => (prev + 1) % shopHomeBanners.length);
  };

  const handleBannerPrev = () => {
    setCurrentBannerIndex((prev) => (prev - 1 + shopHomeBanners.length) % shopHomeBanners.length);
  };

  const handleBannerDotClick = (index: number) => {
    setCurrentBannerIndex(index);
  };

  return (
    <div className="bg-gray-50 min-h-screen">
      <div className="max-w-7xl mx-auto py-6 px-4">
        {/* Shop Header - Compact */}
        <div className="bg-white rounded-lg p-6 mb-4 shadow-sm">
          <div className="flex items-center gap-6">
            {/* Shop Logo */}
            <div className="relative w-20 h-20 rounded-full overflow-hidden border-2 border-blue-500 flex-shrink-0 bg-white">
              <img 
                src={getImageUrl(shop.shopLogo)} 
                alt={shop.shopName}
                className="w-full h-full object-cover"
              />
            </div>

            {/* Shop Info */}
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <h1 className="text-2xl font-bold">{shop.shopName}</h1>
                {shop.taxInfo.taxActiveStatus && (
                  <Badge className="bg-blue-600">
                    <CheckCircle className="w-3 h-3 mr-1" />
                    OFFICIAL
                  </Badge>
                )}
              </div>
              
              <div className="flex items-center gap-6 text-sm text-gray-600">
                {/* <div className="flex items-center gap-1">
                  <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                  <span>4.8</span>
                </div> */}
                <div className="flex items-center gap-1">
                  <Users className="w-4 h-4" />
                  <span>{shop.followerCount.toLocaleString()} Người theo dõi</span>
                </div>
                {/* <div className="flex items-center gap-1">
                  <Package className="w-4 h-4" />
                  <span>10000+ Sản phẩm</span>
                </div> */}
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex gap-3">
              <Button
                variant={isFollowing ? "outline" : "default"}
                onClick={handleFollow}
                size="sm"
              >
                <Heart className={`w-4 h-4 mr-2 ${isFollowing ? 'fill-red-500 text-red-500' : ''}`} />
                {isFollowing ? 'Đã theo dõi' : 'Theo dõi'}
              </Button>
              <Button variant="outline" size="sm">
                <MessageCircle className="w-4 h-4 mr-2" />
                Chat
              </Button>
            </div>
          </div>
        </div>

        {/* Main Content Tabs */}
        <Tabs value={activeTab} onValueChange={handleTabChange} className="w-full">
          <TabsList className="w-full bg-white rounded-lg p-1 mb-4 shadow-sm">
            <TabsTrigger value="shop" className="flex-1">
              <Store className="w-4 h-4 mr-2" />
              Cửa Hàng
            </TabsTrigger>
            <TabsTrigger value="products" className="flex-1">
              <ShoppingBag className="w-4 h-4 mr-2" />
              Tất Cả Sản Phẩm
            </TabsTrigger>
            <TabsTrigger value="profile" className="flex-1">
              <Award className="w-4 h-4 mr-2" />
              Hồ Sơ Của Hàng
            </TabsTrigger>
          </TabsList>

          {/* Tab 1: Cửa Hàng - Banners */}
          <TabsContent value="shop" className="space-y-6">
            {/* Home Banners - Slideshow */}
            {shopHomeBanners.length > 0 && (
              <div 
                className="relative rounded-2xl overflow-hidden"
                onMouseEnter={() => setIsBannerPaused(true)}
                onMouseLeave={() => setIsBannerPaused(false)}
              >
                {/* Slideshow Container */}
                <div className="relative w-full h-[400px] md:h-[500px]">
                  {shopHomeBanners.map((banner, index) => (
                    <div
                      key={banner.id}
                      onClick={() => handleBannerClick(banner.bannerUrl)}
                      className={`absolute inset-0 transition-opacity duration-700 ease-in-out cursor-pointer ${
                        index === currentBannerIndex ? 'opacity-100 z-10' : 'opacity-0 z-0'
                      }`}
                    >
                      <img
                        src={getImageUrl(banner.bannerImage)}
                        alt={banner.bannerName}
                        className="w-full h-full object-cover"
                        onError={(e) => {
                          (e.target as HTMLImageElement).src = '/placeholder.png';
                        }}
                      />
                    </div>
                  ))}
                  
                  {/* Navigation Buttons - Only show if more than 1 banner */}
                  {shopHomeBanners.length > 1 && (
                    <>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="absolute left-4 top-1/2 -translate-y-1/2 z-20 bg-white/80 hover:bg-white shadow-lg"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleBannerPrev();
                        }}
                      >
                        <ChevronLeft className="w-6 h-6" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="absolute right-4 top-1/2 -translate-y-1/2 z-20 bg-white/80 hover:bg-white shadow-lg"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleBannerNext();
                        }}
                      >
                        <ChevronRight className="w-6 h-6" />
                      </Button>
                    </>
                  )}
                  
                  {/* Dots Indicator */}
                  {shopHomeBanners.length > 1 && (
                    <div className="absolute bottom-4 left-1/2 -translate-x-1/2 z-20 flex gap-2">
                      {shopHomeBanners.map((_, index) => (
                        <button
                          key={index}
                          onClick={(e) => {
                            e.stopPropagation();
                            handleBannerDotClick(index);
                          }}
                          className={`w-2 h-2 rounded-full transition-all ${
                            index === currentBannerIndex
                              ? 'bg-white w-8'
                              : 'bg-white/50 hover:bg-white/75'
                          }`}
                          aria-label={`Go to slide ${index + 1}`}
                        />
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Category Banners - Large Grid */}
            {shopCategoryBanners.length > 0 && (
              <div className="space-y-4">
                <h3 className="text-2xl font-bold flex items-center gap-2">
                  <Grid className="w-6 h-6 text-green-600" />
                  Bộ Sưu Tập
                </h3>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  {shopCategoryBanners.map((banner) => (
                    <div
                      key={banner.id}
                      onClick={() => handleBannerClick(banner.bannerUrl)}
                      className="relative rounded-2xl overflow-hidden cursor-pointer hover:shadow-2xl transition-all group bg-white"
                    >
                      <img
                        src={getImageUrl(banner.bannerImage)}
                        alt={banner.bannerName}
                        className="w-full h-[280px] object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                      <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-4">
                        <p className="text-white font-bold text-sm">{banner.bannerName}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Product Banners - Wide Format */}
            {shopProductBanners.length > 0 && (
              <div className="space-y-4">
                <h3 className="text-2xl font-bold flex items-center gap-2">
                  <Package className="w-6 h-6 text-purple-600" />
                  Sản Phẩm Nổi Bật
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {shopProductBanners.map((banner) => (
                    <div
                      key={banner.id}
                      onClick={() => handleBannerClick(banner.bannerUrl)}
                      className="relative rounded-2xl overflow-hidden cursor-pointer hover:shadow-2xl transition-all group bg-gradient-to-br from-purple-50 to-blue-50"
                    >
                      <img
                        src={getImageUrl(banner.bannerImage)}
                        alt={banner.bannerName}
                        className="w-full h-[250px] object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Promotion Banners - Full Width Large Display */}
            {shopPromotionBanners.length > 0 && (
              <div className="space-y-4">
                <h3 className="text-2xl font-bold flex items-center gap-2">
                  <Tag className="w-6 h-6 text-red-600" />
                  Khuyến Mãi Hot
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {shopPromotionBanners.map((banner) => (
                    <div
                      key={banner.id}
                      onClick={() => handleBannerClick(banner.bannerUrl)}
                      className="relative rounded-2xl overflow-hidden cursor-pointer hover:shadow-2xl transition-all group"
                    >
                      <img
                        src={getImageUrl(banner.bannerImage)}
                        alt={banner.bannerName}
                        className="w-full h-[200px] object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                      <div className="absolute top-4 right-4 bg-red-600 text-white px-4 py-2 rounded-full font-bold text-sm shadow-lg">
                        HOT
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Shop Products List */}
            <div className="space-y-4">
              <h3 className="text-2xl font-bold flex items-center gap-2">
                <ShoppingBag className="w-6 h-6 text-blue-600" />
                Sản Phẩm Cửa Hàng
              </h3>
              
              {isLoadingHomeProducts && homeProductsPage === 1 ? (
                <div className="flex items-center justify-center py-12">
                  <Loading size="lg" variant="primary" />
                </div>
              ) : homeProducts.length > 0 ? (
                <>
                  <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                    {homeProducts.map((product) => (
                      <C_ProductSimple key={product.id} product={product} />
                    ))}
                  </div>
                  
                  {/* Load More Button */}
                  {homeProductsData && homeProductsData.totalPages > homeProductsPage && (
                    <div className="flex justify-center pt-6">
                      <Button
                        onClick={handleLoadMoreProducts}
                        disabled={isLoadingHomeProducts}
                        variant="outline"
                        size="lg"
                        className="min-w-[200px]"
                      >
                        {isLoadingHomeProducts ? (
                          <>
                            <Loading size="sm" variant="primary" className="mr-2" />
                            Đang tải...
                          </>
                        ) : (
                          <>
                            Xem thêm
                            <ChevronRight className="w-4 h-4 ml-2" />
                          </>
                        )}
                      </Button>
                    </div>
                  )}
                </>
              ) : !isLoadingHomeProducts ? (
                <div className="bg-white rounded-2xl p-12 text-center shadow-sm">
                  <Package className="w-20 h-20 mx-auto mb-4 text-gray-300" />
                  <p className="text-gray-500 text-lg">Cửa hàng chưa có sản phẩm nào</p>
                </div>
              ) : null}
            </div>

            {/* No banners message */}
            {shopHomeBanners.length === 0 && 
             shopCategoryBanners.length === 0 && 
             shopProductBanners.length === 0 && 
             shopPromotionBanners.length === 0 && (
              <div className="bg-white rounded-2xl p-12 text-center shadow-sm">
                <Store className="w-20 h-20 mx-auto mb-4 text-gray-300" />
                <p className="text-gray-500 text-lg">Cửa hàng chưa có banner nào</p>
              </div>
            )}
          </TabsContent>

          {/* Tab 2: Tất Cả Sản Phẩm */}
          <TabsContent value="products" className="space-y-4">
            <Card>
              <CardContent className="p-6">
                {/* Search and Filter Bar */}
                <div className="mb-6 space-y-4">
                  {/* Search Row */}
                  <div className="flex gap-3">
                    <div className="flex-1 relative">
                      <Input
                        type="text"
                        placeholder="Tìm sản phẩm trong cửa hàng..."
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
                  </div>

                  {/* Filter Row - Tiki Style */}
                  <div className="flex flex-wrap gap-3 items-center">
                    {/* Category Filter */}
                    <Select value={selectedCategory || "all"} onValueChange={handleCategoryChange}>
                      <SelectTrigger className="w-[200px]">
                        <Filter className="w-4 h-4 mr-2" />
                        <SelectValue placeholder="Tất cả danh mục" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="all">Tất cả danh mục</SelectItem>
                        {categoriesData?.map((category) => (
                          <SelectItem key={category.path} value={category.path}>
                            {category.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>

                    {/* Price Filter */}
                    <Button
                      variant="outline"
                      onClick={() => setShowFilters(!showFilters)}
                      className="gap-2"
                    >
                      <SlidersHorizontal className="w-4 h-4" />
                      Giá
                    </Button>

                    {/* Sort */}
                    <Select value={sortBy} onValueChange={handleSortChange}>
                      <SelectTrigger className="w-[180px]">
                        <SelectValue placeholder="Sắp xếp" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="best_sell">Bán chạy nhất</SelectItem>
                        <SelectItem value="price_asc">Giá thấp đến cao</SelectItem>
                        <SelectItem value="price_desc">Giá cao đến thấp</SelectItem>
                        <SelectItem value="name_asc">Tên A-Z</SelectItem>
                        <SelectItem value="name_desc">Tên Z-A</SelectItem>
                      </SelectContent>
                    </Select>

                    {/* Clear Filters */}
                    {(searchKeyword || selectedCategory || priceMin || priceMax) && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleClearFilters}
                        className="text-blue-600"
                      >
                        Xóa bộ lọc
                      </Button>
                    )}
                  </div>

                  {/* Price Filter Panel */}
                  {showFilters && (
                    <div className="p-4 border rounded-lg bg-gray-50 space-y-4">
                      <div className="flex items-center justify-between mb-3">
                        <h3 className="font-semibold">Khoảng giá</h3>
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
                          <label className="text-sm text-gray-600 mb-1 block">Từ</label>
                          <Input
                            type="number"
                            placeholder="0 đ"
                            value={priceMin || ''}
                            onChange={(e) => setPriceMin(e.target.value ? Number(e.target.value) : undefined)}
                          />
                        </div>
                        <div>
                          <label className="text-sm text-gray-600 mb-1 block">Đến</label>
                          <Input
                            type="number"
                            placeholder="1000000 đ"
                            value={priceMax || ''}
                            onChange={(e) => setPriceMax(e.target.value ? Number(e.target.value) : undefined)}
                          />
                        </div>
                      </div>

                      <div className="flex gap-2 justify-end">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => {
                            setPriceMin(undefined);
                            setPriceMax(undefined);
                            updateURLParams({ price_min: undefined, price_max: undefined, page: 1 });
                            setShowFilters(false);
                          }}
                        >
                          Xóa
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
                  {(searchKeyword || selectedCategory || priceMin || priceMax) && (
                    <div className="flex flex-wrap gap-2 items-center pt-2 border-t">
                      <span className="text-sm text-gray-600">Đang lọc:</span>
                      {searchKeyword && (
                        <Badge variant="secondary" className="gap-1">
                          {searchKeyword}
                          <X className="w-3 h-3 cursor-pointer" onClick={() => updateURLParams({ keywords: undefined, page: 1 })} />
                        </Badge>
                      )}
                      {selectedCategory && (
                        <Badge variant="secondary" className="gap-1">
                          {categoriesData?.find(c => c.path === selectedCategory)?.name}
                          <X className="w-3 h-3 cursor-pointer" onClick={() => updateURLParams({ category: undefined, page: 1 })} />
                        </Badge>
                      )}
                      {priceMin && (
                        <Badge variant="secondary" className="gap-1">
                          Từ {priceMin.toLocaleString()}₫
                          <X className="w-3 h-3 cursor-pointer" onClick={() => updateURLParams({ price_min: undefined, page: 1 })} />
                        </Badge>
                      )}
                      {priceMax && (
                        <Badge variant="secondary" className="gap-1">
                          Đến {priceMax.toLocaleString()}₫
                          <X className="w-3 h-3 cursor-pointer" onClick={() => updateURLParams({ price_max: undefined, page: 1 })} />
                        </Badge>
                      )}
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
                    {/* Results count */}
                    <div className="text-sm text-gray-600 mb-4">
                      Tìm thấy <span className="font-semibold">{allProducts.totalElements}</span> sản phẩm
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4 mb-6">
                      {allProducts.data.map((product) => (
                        <C_ProductSimple key={product.id} product={product} />
                      ))}
                    </div>

                    {/* Pagination */}
                    {allProducts.totalPages > 1 && (
                      <div className="flex items-center justify-center gap-2 pt-6 border-t">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handlePageChange(Math.max(1, currentPage - 1))}
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
                                onClick={() => handlePageChange(pageNum)}
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
                          onClick={() => handlePageChange(Math.min(allProducts.totalPages, currentPage + 1))}
                          disabled={currentPage === allProducts.totalPages}
                        >
                          Trang sau
                        </Button>
                      </div>
                    )}
                  </>
                ) : (
                  <div className="text-center text-gray-500 py-12">
                    <Package className="w-16 h-16 mx-auto mb-4 text-gray-300" />
                    <p className="text-lg mb-2">Không tìm thấy sản phẩm</p>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleClearFilters}
                    >
                      Xóa bộ lọc
                    </Button>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          {/* Tab 3: Hồ Sơ Của Hàng */}
          <TabsContent value="profile" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Shop Information */}
              <Card className="lg:col-span-2">
                <CardContent className="p-6">
                  <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                    <Store className="w-5 h-5 text-blue-600" />
                    Thông tin cửa hàng
                  </h2>
                  
                  <div className="space-y-4">
                    {shop.shopDescription && (
                      <div className="p-4 bg-blue-50 rounded-lg">
                        <p className="text-sm text-gray-600 mb-1 font-semibold">Mô tả</p>
                        <p className="text-gray-800">{shop.shopDescription}</p>
                      </div>
                    )}

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
                </CardContent>
              </Card>

              {/* Business Information */}
              <Card>
                <CardContent className="p-6">
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
                      <div className="pt-3 flex items-center gap-2 text-green-600 bg-green-50 p-3 rounded-lg">
                        <CheckCircle className="w-5 h-5" />
                        <span className="font-semibold">Đã xác thực</span>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
