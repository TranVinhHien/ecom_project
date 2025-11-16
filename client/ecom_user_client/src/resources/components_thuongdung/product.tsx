import { Heart, Star, StarIcon } from "lucide-react"
import { Link } from '@/i18n/routing';
import { useTranslations } from "next-intl";
import { ProductSummary } from "@/types/product.types";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import Image from "next/image";
import { toast } from "@/hooks/use-toast";
import ROUTER from "@/assets/configs/routers";
import { Button } from "@/components/ui/button";
import {getImageUrl, formatPrice} from "@/assets/helpers/convert_tool";

const C_ProductSimple = ({ product }: { product: ProductSummary }) => {
  const t = useTranslations("System");
  console.log("C_ProductSimple product:", product);
  return (
    
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
                    {/* Rating */}
                    {product.rating.total_reviews > 0 && (
                      <div className="flex items-center ml-auto  text-xs font-semibold px-2 py-1 rounded">
                        <StarIcon className="w-3 h-3 mr-1 text-yellow-500" />
                        <span> {product.rating.average_rating.toFixed(1)} ({product.rating.total_reviews})</span>
                      </div>
                    )}
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
          
  );
}

export default C_ProductSimple;