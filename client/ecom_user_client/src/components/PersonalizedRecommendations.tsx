"use client"

import React, { useState, useEffect } from 'react';
import C_ProductSimple from "@/resources/components_thuongdung/product";
import apiClient from "@/lib/apiClient";
import API from "@/assets/configs/api";
import { ProductSummary } from "@/types/product.types";

interface PersonalizedRecommendationsProps {
  userId: string;
  title?: string;
  limit?: number;
}

interface RecommendationItem {
  price: number;
  product_id: string;
  rating: number;
  reason: string;
  score: number;
}

interface RecommendationResponse {
  count: number;
  recommendations: RecommendationItem[];
  success: boolean;
}

export default function PersonalizedRecommendations({ 
  userId, 
  title = "Gợi ý dành riêng cho bạn",
  limit = 25 
}: PersonalizedRecommendationsProps) {
  const [products, setProducts] = useState<ProductSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRecommendations = async () => {
      if (!userId) return;

      try {
        setLoading(true);
        setError(null);

        // 1. Gọi API lấy danh sách recommendations
        const recommendationsResponse = await apiClient.post<RecommendationResponse>(
          '/recommendations/personalized',
          {
            user_id: userId,
            n: limit
          },
          {
            customBaseURL: 'http://localhost:5000/api'
          }
        );

        const recommendations = recommendationsResponse.data?.recommendations || [];
        
        if (recommendations.length === 0) {
          setProducts([]);
          return;
        }

        // 2. Lấy danh sách product_ids
        const productIds = recommendations.map(item => item.product_id);

        // 3. Tạo query string
        const queryString = productIds
          .map(id => `product_ids=${id}`)
          .join('&');

        // 4. Gọi API lấy chi tiết sản phẩm
        const productsResponse = await apiClient.get(
          `/product/get_products_detail_for_search?${queryString}`,
          {
            headers: {
              'Content-Type': 'application/json',
            },
            customBaseURL: API.base_product,
          }
        );

        console.log("Products response:", productsResponse.data);

        // 5. Transform data theo ProductSummary interface
        const productsData = productsResponse.data?.result?.data || [];
        
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
            shop_id: '',
            brand_id: '',
            category_id: '',
            min_price: minPriceSku?.price || 0,
            max_price: maxPriceSku?.price || 0,
            min_price_sku_id: minPriceSku?.id || '',
            max_price_sku_id: maxPriceSku?.id || '',
            description: item.product.short_description || '',
            total_sold: 0,
            short_description: item.product.short_description || '',
            media: null,
            product_is_permission_check: item.product.product_is_permission_check,
            product_is_permission_return: item.product.product_is_permission_return,
            delete_status: '',
            create_date: '',
            update_date: '',
            rating: {
              product_id: item.product.id,
              total_reviews: 0,
              average_rating: 0,
            }
          };
        });

        console.log("Transformed products:", transformedProducts);
        setProducts(transformedProducts);

      } catch (err) {
        setError(err instanceof Error ? err.message : 'Lỗi khi tải gợi ý sản phẩm');
        console.error('Error fetching recommendations:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchRecommendations();
  }, [userId, limit]);

  // Loading state
  if (loading) {
    return (
      <div className="w-full px-4 md:px-0">
        <section className="bg-[#f5f5f5] py-4 px-4 md:px-6 rounded-lg mb-8 w-full">
          <div className="flex items-center border-b-2 border-[#ee4d2d] pb-2 mb-4">
            <span className="text-base md:text-lg font-bold text-[#ee4d2d] uppercase tracking-wider">{title}</span>
          </div>
          <div className="flex justify-center items-center py-8">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#ee4d2d]"></div>
          </div>
        </section>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="w-full px-4 md:px-0">
        <section className="bg-[#f5f5f5] py-4 px-4 md:px-6 rounded-lg mb-8 w-full">
          <div className="flex items-center border-b-2 border-[#ee4d2d] pb-2 mb-4">
            <span className="text-base md:text-lg font-bold text-[#ee4d2d] uppercase tracking-wider">{title}</span>
          </div>
          <div className="text-center py-8">
            <p className="text-red-500 text-lg">{error}</p>
          </div>
        </section>
      </div>
    );
  }

  // No products
  if (products.length === 0) {
    return null;
  }

  // Render products
  return (
    <div className="w-full px-4 md:px-0">
      <section className="bg-[#f5f5f5] py-4 px-4 md:px-6 rounded-lg mb-8 w-full">
        <div className="flex items-center border-b-2 border-[#ee4d2d] pb-2 mb-4">
          <span className="text-base md:text-lg font-bold text-[#ee4d2d] uppercase tracking-wider">{title}</span>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 mb-8">
          {products.map(product => (
            <C_ProductSimple 
              key={product.id} 
              product={product} 
              collection_type="click" 
              user_id={userId} 
            />
          ))}
        </div>
      </section>
    </div>
  );
}
