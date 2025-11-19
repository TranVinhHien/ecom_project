"use client"

import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { useRouter } from "@/i18n/routing";
import { Button } from '@/components/ui/button';
import { ShopApiResponse } from '@/types/shop.types';
import { Loading } from '@/components/ui/loading';
import { Store, CheckCircle, MapPin, Star } from 'lucide-react';
import apiClient from '@/lib/apiClient';

interface ShopInfoProps {
  shopId: string;
}

export default function ShopInfo({ shopId }: ShopInfoProps) {
  const router = useRouter();

  const { data, isLoading, error } = useQuery<ShopApiResponse>({
    queryKey: ['shop-info', shopId],
    queryFn: async () => {
      const response = await apiClient.get(`/Shops/${shopId}`, {
        customBaseURL: 'http://localhost:8000/api'
      });
      return response.data;
    },
    enabled: !!shopId,
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-4">
        <Loading size="sm" variant="primary" />
      </div>
    );
  }

  if (error || !data?.result) {
    return null;
  }

  const shop = data.result;

  const handleViewShop = () => {
    router.push(`/shop/${shopId}`);
  };

  return (
    <div className="bg-white rounded-lg p-4 border border-gray-200 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-center gap-3">
        {/* Shop Logo */}
        <div 
          className="relative w-14 h-14 rounded-full overflow-hidden border-2 border-gray-100 flex-shrink-0 cursor-pointer hover:scale-105 transition-transform"
          onClick={handleViewShop}
        >
          <img 
            src={shop.shopLogo} 
            alt={shop.shopName}
            className="w-full h-full object-cover"
            onError={(e) => {
              (e.target as HTMLImageElement).src = '/placeholder-shop.png';
            }}
          />
          {shop.taxInfo.taxActiveStatus && (
            <div className="absolute bottom-0 right-0 bg-blue-600 rounded-full p-0.5">
              <CheckCircle className="w-2.5 h-2.5 text-white" />
            </div>
          )}
        </div>

        {/* Shop Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h3 
              className="text-sm font-bold text-gray-900 hover:text-blue-600 cursor-pointer transition-colors truncate"
              onClick={handleViewShop}
            >
              {shop.shopName}
            </h3>
            {shop.taxInfo.taxActiveStatus && (
              <span className="bg-blue-600 text-white text-[10px] px-1.5 py-0.5 rounded flex items-center gap-0.5 flex-shrink-0">
                <Store className="w-2.5 h-2.5" />
                <span className="font-semibold">OFFICIAL</span>
              </span>
            )}
          </div>
          
          <div className="flex items-center gap-3 text-xs text-gray-600">
            <div className="flex items-center gap-1">
              {/* <Star className="w-3 h-3 text- yellow-500 fill-yellow-500" /> */}
              {/* <span className="font-semibold">4.8</span> */}
            </div>
            <div className="flex items-center gap-1">
              <span className="font-semibold">{shop.followerCount.toLocaleString()}</span>
              <span>Theo dõi</span>
            </div>
          </div>
        </div>

        {/* View Shop Button */}
        <Button
          variant="outline"
          size="sm"
          onClick={handleViewShop}
          className="border-blue-600 text-blue-600 hover:bg-blue-50 text-xs px-3 h-8 flex-shrink-0"
        >
          <Store className="w-3 h-3 mr-1" />
          Xem Shop
        </Button>
      </div>
    </div>
  );
}
