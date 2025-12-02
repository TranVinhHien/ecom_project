"use client"

import C_ProductSimple from "@/resources/components_thuongdung/product";
import { useTranslations } from "next-intl";
import React, { useEffect, useState } from "react";
import { useGetActiveBanners, useGetProducts } from "@/services/apiService";
import { ProductSummary } from "@/types/product.types";
import { UserProfile } from "@/types/user.types";
import { INFO_USER } from "@/assets/configs/request";
import { event_type } from "@/types/collection.types";
import PersonalizedRecommendations from "@/components/PersonalizedRecommendations";
import HomeBannerSlider from "@/components/HomeBannerSlider";

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
      <HomeBannerSlider banners={homeBanners || []}   height="h-[400px] md:h-[500px]" />
      <div className="container mx-auto px-4 md:px-8 lg:px-16 py-8">
        {profile?.id && (
          <PersonalizedRecommendations 
            userId={profile.id} 
            title="Gợi ý dành riêng cho bạn" 
            limit={25}
          />
        )}
        <HomeSuggestion products={products} user_id={profile?.id || ""} title="Bán chạy nhất điện tử" t={t} collection_type="click" />
        <HomeSuggestion products={data_cham_soc_nha_cua?.data || []} user_id={profile?.id || ""} title="Chăm sóc nhà cửa" t={t} collection_type="click" />
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
    <section className="bg-[#f5f5f5] py-4 px-4 md:px-6 rounded-lg mb-8 w-full">
      <div className="flex items-center border-b-2 border-[#ee4d2d] pb-2 mb-4">
        <span className="text-base md:text-lg font-bold text-[#ee4d2d] uppercase tracking-wider">{title}</span>
      </div>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
        {products.length > 0 && products?.map(product => (
          <C_ProductSimple key={product.id} product={product} collection_type={collection_type} user_id={user_id} />
        ))}
      </div>
    </section>
  );
}
