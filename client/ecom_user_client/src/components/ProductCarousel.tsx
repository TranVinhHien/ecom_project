"use client";

import { useRouter } from "@/i18n/routing";
// import Image from "next/image";
import { Card } from "@/components/ui/card";
import ROUTER from "@/assets/configs/routers";

interface Product {
  brand: string;
  category: string;
  product: {
    id: string;
    image: string;
    key: string;
    name: string;
    product_is_permission_check: boolean;
    product_is_permission_return: boolean;
    short_description: string;
  };
  sku: Array<{
    id: string;
    price: number;
    quantity: number;
    sku_name: string;
  }>;
  similarity_score: number;
}

interface ProductCarouselProps {
  products: Product[];
}

export default function ProductCarousel({ products }: ProductCarouselProps) {
  const router = useRouter();

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  const getMinPrice = (skus: Product["sku"]) => {
    if (!skus || skus.length === 0) return 0;
    return Math.min(...skus.map((sku) => sku.price));
  };

  const handleProductClick = (key: string) => {
    router.push(`/${ROUTER.product}/${key}`);
  };

  return (
    <div className="mt-2">
      <p className="text-xs font-semibold mb-2 text-gray-700">Sản phẩm gợi ý:</p>
      <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-gray-100">
        {products.map((item, index) => {
          const minPrice = getMinPrice(item.sku);
          
          return (
            <Card
              key={`${item.product.id}-${index}`}
              className="flex-shrink-0 w-[160px] cursor-pointer hover:shadow-lg transition-shadow duration-200"
              onClick={() => handleProductClick(item.product.key)}
            >
              <div className="p-2">
                <div className="relative w-full h-[120px] mb-2 bg-gray-100 rounded-md overflow-hidden">
                  <img
                    src={item.product.image || "/placeholder-product.png"}
                    alt={item.product.name}
                    // fill
                    className="object-cover"
                    sizes="160px"
                  />
                  {/* {item.similarity_score && (
                    <div className="absolute top-1 right-1 bg-green-500 text-white text-[10px] px-1.5 py-0.5 rounded-full font-semibold">
                      {Math.round(item.similarity_score * 100)}%
                    </div>
                  )} */}
                </div>
                <h4 className="text-xs font-medium line-clamp-2 mb-1 min-h-[32px]">
                  {item.product.name}
                </h4>
                <p className="text-sm font-bold text-[hsl(var(--primary))]">
                  {formatPrice(minPrice)}
                </p>
                {item.brand && (
                  <p className="text-[10px] text-gray-500 mt-1">
                    Thương hiệu: {item.brand}
                  </p>
                )}
              </div>
            </Card>
          );
        })}
      </div>
    </div>
  );
}
