"use client"

import API from "@/assets/configs/api";
import { handleProductImg } from "@/assets/configs/handle_img";
import * as request from "@/assets/helpers/request_without_token";
import { MetaType, ParamType } from "@/assets/types/request";
import C_ProductSimple from "@/resources/components_thuongdung/product";
import productSimple from "@/resources/components_thuongdung/product";
import { Search, StarIcon, Camera } from "lucide-react";
import { useTranslations } from "next-intl";
import { useSearchParams } from "next/navigation";
import React, { useRef, useState, useEffect } from "react";
import apiClient from "@/lib/apiClient";
import { UserProfile } from "@/types/user.types";
import { INFO_USER } from "@/assets/configs/request";

function base64ToFile(base64: string, filename: string): File {
  const arr = base64.split(',');
  const mime = arr[0].match(/:(.*?);/)?.[1] || 'image/png';
  const bstr = atob(arr[1]);
  let n = bstr.length;
  const u8arr = new Uint8Array(n);

  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }

  return new File([u8arr], filename, { type: mime });
}

export default function SearchPage() {
  const [searchImage, setSearchImage] = useState<string | null>(null);
  const [products, setProducts] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [prediction, setPrediction] = useState<{label: string, confidence: number} | null>(null);
  
  useEffect(() => {
    // Load image from localStorage when component mounts
    const storedImage = localStorage.getItem('searchImage');
    if (storedImage) {
      setSearchImage(storedImage);
    }
  }, []);

  useEffect(() => {
    const fetchProducts = async () => {
      if (!searchImage) return;

      try {
        setLoading(true);
        setError(null);
        
        const formData = new FormData();
        formData.append('file', base64ToFile(searchImage, 'search_image.png'));
        
        const response = await apiClient.post('/search/image', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
          customBaseURL: API.AI_CHUONG_BASE,
        });

        // const data = response.data;

        const data_search_image:{
          results:{
            matched_image: string,
            price: number,
            product_id: string,
            similarity: number
          }[]
        } = response.data
        console.log("Search image results:", data_search_image)
      
        const productIds = data_search_image.results
            .filter(item => item.similarity > 0.7)
            .map(item => item.product_id);

        if (productIds.length === 0) {
          setProducts([]);
          return;
        }

        // 1. Tạo query string thủ công: "product_ids=id1&product_ids=id2"
        const queryString = productIds
            .map(id => `product_ids=${id}`)
            .join('&');

        // 2. Gọi API (Nối trực tiếp vào URL)
        const response_product = await apiClient.get(`/product/get_products_detail_for_search?${queryString}`, {
            headers: {
                'Content-Type': 'application/json',
            },
            customBaseURL: API.base_product,
        });

        console.log("response_product:", response_product.data);

        // Transform data from API response to match ProductSummary interface
        const productsData = response_product.data?.result?.data || [];
        
        const transformedProducts = productsData.map((item: any) => {
          // Tìm SKU có giá thấp nhất
          const minPriceSku = item.sku?.reduce((min: any, sku: any) => 
            sku.price < min.price ? sku : min, item.sku[0]
          );
          
          // Tìm SKU có giá cao nhất
          const maxPriceSku = item.sku?.reduce((max: any, sku: any) => 
            sku.price > max.price ? sku : max, item.sku[0]
          );

          // Map đúng theo ProductSummary interface
          return {
            id: item.product.id,
            name: item.product.name,
            key: item.product.key,
            image: item.product.image,
            shop_id: '', // Không có trong response, để empty
            brand_id: '', // Không có trong response, chỉ có brand name
            category_id: '', // Không có trong response, chỉ có category name
            min_price: minPriceSku?.price || 0,
            max_price: maxPriceSku?.price || 0,
            min_price_sku_id: minPriceSku?.id || '',
            max_price_sku_id: maxPriceSku?.id || '',
            description: item.product.short_description || '',
            total_sold: 0, // Không có trong response
            short_description: item.product.short_description || '',
            media: null, // Không có trong response
            product_is_permission_check: item.product.product_is_permission_check,
            product_is_permission_return: item.product.product_is_permission_return,
            delete_status: '', // Không có trong response
            create_date: '', // Không có trong response
            update_date: '', // Không có trong response
            rating: {
              product_id: item.product.id,
              total_reviews: 0, // Không có trong response
              average_rating: 0, // Không có trong response
            }
          };
        });

        console.log("Transformed products:", transformedProducts);
        setProducts(transformedProducts);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Lỗi khi tìm kiếm sản phẩm');
        console.error('Error fetching products:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, [searchImage]);

  const t = useTranslations("System");

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

  return (
    <div className="min-h-screen p-8 pb-20 font-[family-name:var(--font-geist-sans)]">
      <div className="max-w-7xl mx-auto">
        {/* Image Search Display */}
        {searchImage && (
          <div className="mb-8 p-6 bg-white rounded-lg shadow-md">
            <h2 className="text-2xl font-semibold mb-4 text-gray-800">{t("hinh_anh_tim_kiem")}</h2>
            <div className="flex items-center gap-6">
              <div className="relative w-72 h-72 rounded-lg overflow-hidden border-2 ">
                <img 
                  src={searchImage} 
                  alt="Search image" 
                  className="w-full h-full object-contain"
                />
              </div>
              <div className="flex-1">
                {loading && (
                  <div className="flex items-center gap-2">
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-[hsl(var(--primary))]"></div>
                    <span className="text-gray-600">{t("dang_tim_kiem")}</span>
                  </div>
                )}
                {prediction && !loading && (
                  <div className="mt-4">
                    <div className="text-lg font-semibold text-gray-800">
                      {t("ket_qua_du_doan")}: {prediction.label}
                    </div>
                    <div className="text-md text-gray-600">
                      {t("do_tin_cay")}: {(prediction.confidence * 100).toFixed(2)}%
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Search Results */}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
          {products.map((product, index) => (
            <C_ProductSimple key={index} product={product} user_id={profile?.id} collection_type="search" />
          ))}
        </div>

        {/* Loading State */}
        {loading && !searchImage && (
          <div className="flex justify-center items-center py-8">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#ee4d2d]"></div>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="text-center py-8">
            <p className="text-red-500 text-lg">{error}</p>
          </div>
        )}

        {/* No Results */}
        {!loading && products.length === 0 && !error && (
          <div className="text-center py-8">
            <p className="text-gray-500 text-lg">{t("khong_tim_thay_san_pham_phu_hop")}</p>
          </div>
        )}
      </div>
    </div>
  );
} 