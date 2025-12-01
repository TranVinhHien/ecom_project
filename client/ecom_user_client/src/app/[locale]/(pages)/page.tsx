"use client"

import C_ProductSimple from "@/resources/components_thuongdung/product";
import { useTranslations } from "next-intl";
import React, { useEffect, useMemo, useState } from "react";
import { useGetActiveBanners, useGetProducts } from "@/services/apiService";
import { ProductSummary } from "@/types/product.types";
import { Banner } from "@/types/shop.types";
import { getImageUrl } from "@/assets/helpers/convert_tool";
import { useRouter } from "@/i18n/routing";
import { Image } from "lucide-react";
import { UserProfile } from "@/types/user.types";
import { INFO_USER } from "@/assets/configs/request";
import { event_type } from "@/types/collection.types";

export default function Home() {
  const t = useTranslations("System");
  // Sử dụng hook useGetProducts giống các trang khác
  const { data, isLoading, error } = useGetProducts({
    page: 1,
    limit: 20,
    // price_min:60000,
    sort:'best_sell',
    cate_path:'/laptop-may-vi-tinh-linh-kien'
  });

  const { data: data_cham_soc_nha_cua } = useGetProducts({
    page: 1,
    limit: 20,
    // price_min:60000,
    sort:'best_sell',
    cate_path:'/cham-soc-nha-cua'
  });
  const { data: homeBanners } = useGetActiveBanners('HOME');

  const [profile, setProfile] = useState<UserProfile | null>(null);
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

  if (isLoading) return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-[#ee4d2d]"></div>
    </div>
  );
  
  if (error) return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="text-red-500 text-xl">Lỗi: {error.message}</div>
    </div>
  );

  const products = data?.data || [];


  return (
    <div className="pt-36">
      <HomeBannerSlider banners={homeBanners || []} />
      <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
        <div className="mb-64" />
        <TopSearchSection items={products} t={t} />
        <div className="mb-8" />
          <HomeSuggestion products={products} user_id={profile?.id || ""} title="Bán chạy nhất điện tử" t={t} collection_type="click" />
        <HomeSuggestion products={data_cham_soc_nha_cua?.data || []} user_id={profile?.id || ""} title="Chăm sóc nhà cửa" t={t} collection_type="click" />
      </div>
    </div>
  );
}

function HomeBannerSlider({ banners }: { banners: Banner[] }) {
  const router = useRouter();
  const sortedBanners = useMemo(() => {
    return [...banners].sort((a, b) => (a.bannerOrder || 0) - (b.bannerOrder || 0));
  }, [banners]);
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    if (sortedBanners.length <= 1) return;
    const interval = setInterval(() => {
      setCurrentIndex((prev) => (prev + 1) % sortedBanners.length);
    }, 5000);
    return () => clearInterval(interval);
  }, [sortedBanners.length]);

  if (sortedBanners.length === 0) {
    return null;
  }

  const handleNavigate = (rawUrl?: string) => {
    if (!rawUrl) return;
    const url = rawUrl.trim();
    if (url.startsWith("http")) {
      window.open(url, "_blank", "noopener noreferrer");
      return;
    }
    router.push(url);
  };

  return (
    <div className="w-full px-4 md:px-8 lg:px-16 mb-10">
      <div className="relative w-full h-48 md:h-64 lg:h-80 rounded-2xl overflow-hidden shadow-lg">
        {sortedBanners.map((banner, index) => (
          <button
            key={banner.id}
            type="button"
            onClick={() => handleNavigate(banner.bannerUrl)}
            className={`absolute inset-0 transition-opacity duration-700 ${currentIndex === index ? "opacity-100 z-10" : "opacity-0"}`}
            aria-label={banner.bannerName}
          >
             {/* <img
      src={getImageUrl(banner.bannerImage)}
      alt={banner.bannerName}      
      className="object-cover w-full h-[180px] md:h-[300px] lg:h-[400px] xl:h-[500px]"
      sizes="100vw"
      /> */}

  <img
    src={getImageUrl(banner.bannerImage)}
    alt={banner.bannerName}
    
    // SỬA: Thêm "block" và "mx-auto". Bỏ "center"
    className="object-cover w-[70vw] block mx-auto h-[250px] md:h-[400px] lg:h-[500px] xl:h-[650px]"
    
    sizes="70vw"
/>
          </button>
        ))}

        {sortedBanners.length > 1 && (
          <>
            <button
              type="button"
              className="absolute left-4 top-1/2 -translate-y-1/2 bg-black/40 text-white rounded-full p-2 z-20"
              onClick={() => setCurrentIndex((prev) => (prev - 1 + sortedBanners.length) % sortedBanners.length)}
              aria-label="Previous banner"
            >
              ‹
            </button>
            <button
              type="button"
              className="absolute right-4 top-1/2 -translate-y-1/2 bg-black/40 text-white rounded-full p-2 z-20"
              onClick={() => setCurrentIndex((prev) => (prev + 1) % sortedBanners.length)}
              aria-label="Next banner"
            >
              ›
            </button>
            <div className="absolute bottom-4 inset-x-0 flex justify-center gap-2 z-20">
              {sortedBanners.map((_, index) => (
                <span
                  key={`dot-${index}`}
                  className={`h-2 w-2 rounded-full ${currentIndex === index ? "bg-white" : "bg-white/50"}`}
                />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}

function TopSearchSection({ items, t }: { items: ProductSummary[], t: any }) {

  // console.log("TopSearchSection items:", items);
  // Hiển thị tối đa 4 item, nếu nhiều hơn thì cho scroll ngang và hiện mũi tên
  const showArrow = items?.length > 4;
  return (
    <div className="w-full px-4 md:px-0">
      <section className="bg-[#f5f5f5] py-4 px-4 md:px-6 rounded-lg mb-8 w-full">
        <div className="flex justify-between items-center mb-2">
          <span className="text-base md:text-lg font-semibold text-[#ee4d2d]">{t("tim_kiem_hang_dau")}</span>
          <a href="#" className="text-xs md:text-sm text-[#ee4d2d] hover:underline">{t("xem_tat_ca")} &gt;</a>
        </div>
        <div className="relative">
          <div className={`grid grid-cols-2 md:grid-cols-4 gap-3 md:gap-4 ${showArrow ? 'overflow-x-auto scrollbar-hide' : ''}`}
               style={{ minWidth: showArrow ? 0 : 'unset' }}>
            {items.length > 0 && items?.slice(0, 4).map((item, idx) => (
              <div key={item.id || idx} className="bg-white rounded-lg shadow border flex flex-col items-center p-2 md:p-3 relative">
                <div className="absolute left-1 md:left-2 top-1 md:top-2 bg-[#ee4d2d] text-white text-xs font-bold px-1.5 md:px-2 py-0.5 rounded">TOP</div>
                <img 
                  src={item.image} 
                  alt={item.name} 
                  className="w-12 h-12 md:w-16 md:h-16 object-contain mb-2"
                  onError={(e) => {
                    const target = e.target as HTMLImageElement;
                    target.src = "/placeholder.png";
                  }}
                />
                <div className="text-lg font-bold text-red-500 mb-1">
                  {(item.min_price || 0).toLocaleString()}₫
                </div>
                <div className="text-sm text-gray-700 text-center line-clamp-2">{item.name}</div>
              </div>
            ))}
          </div>
          {showArrow && (
            <div className="absolute right-0 top-1/2 -translate-y-1/2 bg-white rounded-full shadow p-1 cursor-pointer z-10">
              <svg width="24" height="24" fill="none" stroke="#ee4d2d" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M9 18l6-6-6-6"/></svg>
            </div>
          )}
        </div>
      </section>
    </div>
  );
}

function HomeSuggestion({ products, title, t,user_id,collection_type }: { products: ProductSummary[], title: string, t: any, user_id: string, collection_type: event_type }) {
 
  return (
    <div className="w-full px-4 md:px-0">
      <section className="bg-[#f5f5f5] py-4 px-4 md:px-6 rounded-lg mb-8 w-full">
        <div className="flex items-center border-b-2 border-[#ee4d2d] pb-2 mb-4">
          <span className="text-base md:text-lg font-bold text-[#ee4d2d] uppercase tracking-wider">{title}</span>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 mb-8">
          {products.length > 0 && products?.map(product => (
            <C_ProductSimple key={product.id} product={product} collection_type={collection_type} user_id={user_id} />
          ))}
        </div>
       
      </section>
    </div>
  );
}
