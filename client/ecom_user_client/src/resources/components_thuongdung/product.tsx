import { StarIcon } from "lucide-react"
import { Link } from '@/i18n/routing';
import { useTranslations } from "next-intl";
import { ProductSummary } from "@/types/product.types";

const C_ProductSimple = ({ product }: { product: ProductSummary }) => {
  const t = useTranslations("System");
  
  return (
    <div className="bg-white rounded-lg border shadow hover:shadow-lg transition-shadow flex flex-col h-full">
      <Link href={`/product/${product.key}`} className="flex flex-col h-full">
        <img 
          src={product.image} 
          alt={product.name} 
          className="w-full h-48 object-cover rounded-t-lg"
          onError={(e) => {
            const target = e.target as HTMLImageElement;
            target.src = "/placeholder.png";
          }}
        />
        <div className="p-4 flex flex-col flex-1">
          <div className="font-medium text-base line-clamp-2 mb-2 min-h-[48px]">
            {product.name}
          </div>
          <div className="mt-auto space-y-2">
            <div className="flex items-center justify-between">
              <div className="text-lg font-bold text-red-500">
                {(product.min_price || 0).toLocaleString()}₫
                {product.max_price && product.max_price !== product.min_price && (
                  <span className="text-sm"> - {product.max_price.toLocaleString()}₫</span>
                )}
              </div>
            </div>
            <div className="text-sm text-gray-500">
              {product.shop_id}
            </div>
          </div>
        </div>
      </Link>
    </div>
  );
}

export default C_ProductSimple;