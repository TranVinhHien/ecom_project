"use client"

import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useRouter, usePathname } from "@/i18n/routing"
import { ShoppingCart, Heart, Minus, Plus, ChevronDown, ChevronUp } from 'lucide-react';
import apiClient from '@/lib/apiClient';
import { Loading } from "@/components/ui/loading";
import { cn } from '@/lib/utils';
import { ProductDetailApiResponse, ProductSKU, ProductOptionValue } from '@/types/product.types';
import { useCartStore } from '@/store/cartStore';
import { useCheckoutStore } from '@/store/checkoutStore';
import ROUTER from '@/assets/configs/routers';
import API from '@/assets/configs/api';
import ProductComments from '@/components/ProductComments';
import { useChatStore } from '@/store/chatStore';
import { getImageUrl } from '@/assets/helpers/convert_tool';
import ImageGalleryModal from '@/components/ImageGalleryModal';

// Image Slider Component
const ProductImageSlider = ({ images, currentIndex, setCurrentIndex, optionImage, onImageClick }: {
  images: string[];
  currentIndex: number;
  setCurrentIndex: (index: number) => void;
  optionImage: string | null;
  onImageClick: (index: number) => void;
}) => {
  const [isAutoSlide, setIsAutoSlide] = useState(true);
// const [isVideoPlaying, setIsVideoPlaying] = useState(false);

  useEffect(() => {
    if (!isAutoSlide) return;
    
    const interval = setInterval(() => {
      setCurrentIndex((currentIndex + 1) % images.length);
    }, 10000);
    return () => clearInterval(interval);
  }, [currentIndex, images.length, isAutoSlide, setCurrentIndex]);

  // Display option image temporarily when selected
  const displayImage = optionImage || images[currentIndex];

  return (
    <div className="sticky top-8">
      {/* Main Image */}
      <div 
        className="relative w-full aspect-square bg-gray-100 rounded-lg overflow-hidden mb-4 cursor-pointer group"
        onClick={() => onImageClick(currentIndex)}
      >
        <Button
          variant="ghost"
          size="icon"
          className="absolute top-1/2 left-2 z-10 -translate-y-1/2 bg-white/80 hover:bg-white"
          onClick={() => {
            setIsAutoSlide(true);
            setCurrentIndex((currentIndex - 1 + images.length) % images.length);
          }}
        >
          <svg width="24" height="24" fill="none" stroke="currentColor"><path d="M15 18l-6-6 6-6"/></svg>
        </Button>
        
       {
              displayImage.endsWith(".mp4") ? 
                      <video
              src={getImageUrl(displayImage)}
              className="w-full h-full object-cover"
              controls
              autoPlay
              muted
              playsInline
              onPlay={() => setIsAutoSlide(false)}   // ✅ báo là video đang chạy
              onEnded={() => {
                setIsAutoSlide(true)     // ✅ báo video xong
                setCurrentIndex((currentIndex + 1) % images.length); // chuyển slide sau khi video kết thúc
              }}

            />
            :   <img 
              src={getImageUrl(displayImage)} 
              alt={`Thumbnail `}
              className="w-full h-full object-cover"
              onError={(e) => {
                (e.target as HTMLImageElement).src = '/placeholder.png';
              }}
            />
            }
          
        
        <Button
          variant="ghost"
          size="icon"
          className="absolute top-1/2 right-2 z-10 -translate-y-1/2 bg-white/80 hover:bg-white"
          onClick={() => {
            setIsAutoSlide(true);
            setCurrentIndex((currentIndex + 1) % images.length);
          }}
        >
          <svg width="24" height="24" fill="none" stroke="currentColor"><path d="M9 6l6 6-6 6"/></svg>
        </Button>

        {/* Zoom overlay on hover */}
        <div className="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-all flex items-center justify-center">
          <svg 
            className="w-12 h-12 text-white opacity-0 group-hover:opacity-100 transition-opacity"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
          </svg>
        </div>

        {/* Image Counter */}
        <div className="absolute bottom-2 right-2 bg-black/60 text-white text-xs px-2 py-1 rounded">
          {optionImage ? 'Option Preview' : `${currentIndex + 1} / ${images.length}`}
        </div>
      </div>

      {/* Thumbnails */}
      <div className="grid grid-cols-6 gap-2">
        {images.map((img, idx) => (
          <button
            key={idx}
            className={cn(
              "aspect-square border-2 rounded overflow-hidden transition-all",
              currentIndex === idx && !optionImage ? "border-primary" : "border-gray-200 hover:border-gray-400"
            )}
            onClick={() => {
              setIsAutoSlide(true);
              setCurrentIndex(idx);
            }}
          >
            {
              img.endsWith(".mp4") ? 
              <video
                src={getImageUrl(img)}
                className="w-full h-full object-cover"
              /> :   <img 
              src={getImageUrl(img)} 
              alt={`Thumbnail ${idx + 1}`}
              className="w-full h-full object-cover"
              onLoad={() => setIsAutoSlide(true)} // ✅ ảnh thì cho phép auto chuyển sau 8s

              onError={(e) => {
                (e.target as HTMLImageElement).src = '/placeholder.png';
              }}
            />
            }
          
          </button>
        ))}
      </div>
    </div>
  );
};

export default function ProductDetailPage({ params }: { params: { id: string } }) {
  
  const t = useTranslations("System");
  const router = useRouter();
  const { addToCart } = useCartStore();
  const { setCheckoutItems } = useCheckoutStore();
  const { openChatWithMessage } = useChatStore();

  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [selectedOptions, setSelectedOptions] = useState<Record<string, string>>({});
  const [selectedSku, setSelectedSku] = useState<ProductSKU | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [optionPreviewImage, setOptionPreviewImage] = useState<string | null>(null);
  const [previewTimer, setPreviewTimer] = useState<NodeJS.Timeout | null>(null);
  const [isDescriptionExpanded, setIsDescriptionExpanded] = useState(false);
  const [isGalleryOpen, setIsGalleryOpen] = useState(false);
  const [galleryInitialIndex, setGalleryInitialIndex] = useState(0);

  // Fetch product detail
  const { data, isLoading, error } = useQuery<ProductDetailApiResponse>({
    queryKey: ['product-detail', params.id],
    queryFn: async () => {
      const response = await apiClient.get(`/product/getdetail/${params.id}`,{
        customBaseURL:process.env.NEXT_PUBLIC_API_GATEWAY_URL||API.base_product 
      });
      console.log('Product Detail Response:', response.data);
      return response.data;
    },
  });

  const productData = data?.result?.data;
  const product = productData?.product;
  const skus = productData?.sku || [];
  const options = productData?.option || [];

  // Parse media (main image + additional media)
  const images = React.useMemo(() => {
    if (!product) return [];
    const allImages = [product.image];
    if (product.media) {
      // convert media from json to array
      const mediaArray = JSON.parse(product.media);

      // const mediaArray = product.media.split(',').map(m => m.trim());
      // console.log('Parsed mediaArray:', mediaArray);
      allImages.push(...mediaArray);
    }
    return allImages.filter(Boolean);
  }, [product]);

  // Find matching SKU based on selected options
  useEffect(() => {
    if (!options.length || !skus.length) return;

    const selectedOptionIds = Object.values(selectedOptions);
    if (selectedOptionIds.length !== options.length) {
      setSelectedSku(null);
      return;
    }

    const matchedSku = skus.find(sku => {
      return selectedOptionIds.every(optionId => 
        sku.option_value_ids.includes(optionId)
      );
    });

    setSelectedSku(matchedSku || null);
  }, [selectedOptions, skus, options]);

  // Handle option selection
  const handleOptionSelect = (optionName: string, optionValueId: string, image?: string) => {
    setSelectedOptions(prev => ({
      ...prev,
      [optionName]: optionValueId
    }));

    // If option has image, show it temporarily
    if (image) {
      // Clear existing timer if any
      if (previewTimer) {
        clearTimeout(previewTimer);
      }

      // Show option image
      setOptionPreviewImage(image);

      // After 5 seconds, return to normal slider
      const timer = setTimeout(() => {
        setOptionPreviewImage(null);
      }, 5000);

      setPreviewTimer(timer);
      
    }
  };

  // Cleanup timer on unmount
  useEffect(() => {
    return () => {
      if (previewTimer) {
        clearTimeout(previewTimer);
      }
    };
  }, [previewTimer]);

  // Handle image click to open gallery
  const handleImageClick = (index: number) => {
    setGalleryInitialIndex(index);
    setIsGalleryOpen(true);
  };

  // Generate notes for product images
  const imageNotes = React.useMemo(() => {
    if (!product) return [];
    return images.map((_, idx) => {
      if (idx === 0) {
        return `Ảnh chính: ${product.name}`;
      }
      return `Ảnh ${idx}: ${product.name}`;
    });
  }, [images, product]);

  // Handle Add to Cart
  const handleAddToCart = () => {
    if (!selectedSku || !product) {
      alert(t("vui_long_chon_day_du_thuoc_tinh"));
      return;
    }

    if (selectedSku.quantity === 0) {
      alert(t("san_pham_het_hang"));
      return;
    }

    if (quantity > selectedSku.quantity) {
      alert(t("so_luong_vuot_qua_ton_kho"));
      return;
    }

    const selectedOptionsText = Object.entries(selectedOptions).map(([key, value]) => {
      const option = options.find(opt => opt.option_name === key);
      const optionValue = option?.values.find(v => v.option_value_id === value);
      return `${key}: ${optionValue?.value || ''}`;
    }).join(', ');

    addToCart({
      sku_id: selectedSku.id,
      shop_id: product.shop_id,
      name: `${product.name} (${selectedOptionsText})`,
      image: product.image,
      price: selectedSku.price,
      quantity: quantity,
      sku_name:selectedSku.sku_name,
    });

    alert(t("them_vao_gio_hang_thanh_cong"));
  };

  // Handle Ask AI for Product Suggestion
  const handleAskAI = () => {
    if (!product) return;
    
    const message = `Gợi ý sản phẩm ${product.name}`;
    openChatWithMessage(message, product.key);
  };

  // Handle Buy Now
  const handleBuyNow = () => {
    if (!selectedSku || !product) {
      alert(t("vui_long_chon_day_du_thuoc_tinh"));
      return;
    }

    if (selectedSku.quantity === 0) {
      alert(t("san_pham_het_hang"));
      return;
    }

    if (quantity > selectedSku.quantity) {
      alert(t("so_luong_vuot_qua_ton_kho"));
      return;
    }

    const selectedOptionsText = Object.entries(selectedOptions).map(([key, value]) => {
      const option = options.find(opt => opt.option_name === key);
      const optionValue = option?.values.find(v => v.option_value_id === value);
      return `${key}: ${optionValue?.value || ''}`;
    }).join(', ');

    // Set checkout items in Zustand store
    setCheckoutItems([{
      sku_id: selectedSku.id,
      shop_id: product.shop_id,
      quantity: quantity,
      // Additional info for display
      name: `${product.name} (${selectedOptionsText})`,
      price: selectedSku.price,
      image: product.image,
    }]);

    router.push(ROUTER.thanhtoan);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loading size="lg" variant="primary" />
      </div>
    );
  }

  if (error || !productData) {
    return (
      <div className="flex items-center justify-center min-h-[400px] text-red-500">
        {t("co_loi_xay_ra_khi_tai_du_lieu")}
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto py-8 px-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Left: Images */}
        <ProductImageSlider 
          images={images}
          currentIndex={currentImageIndex}
          setCurrentIndex={setCurrentImageIndex}
          optionImage={optionPreviewImage}
          onImageClick={handleImageClick}
        />

        {/* Right: Product Info */}
        <div>
          {/* Brand & Category */}
          <div className="flex gap-2 text-sm text-gray-600 mb-2">
            {productData.brand && (
              <span className="text-primary font-medium">{productData.brand.name}</span>
            )}
            {productData.category && (
              <span>/ {productData.category.name}</span>
            )}
          </div>

          {/* Product Name */}
          <h1 className="text-2xl font-bold mb-4">{product?.name}</h1>

          {/* Short Description */}
          {product?.short_description && (
            <p className="text-gray-600 mb-4">{product.short_description}</p>
          )}

          {/* Price */}
          <div className="bg-gray-50 p-4 rounded-lg mb-6">
            <div className="text-3xl font-bold text-primary">
              {selectedSku ? (
                `${selectedSku.price.toLocaleString()}₫`
              ) : (
                `${product?.min_price.toLocaleString()}₫ - ${product?.max_price.toLocaleString()}₫`
              )}
            </div>
          </div>

          {/* Options Selection */}
          {options.map(option => (
            <div key={option.option_name} className="mb-6">
              <div className="font-semibold mb-3">{option.option_name}:</div>
              <div className="flex flex-wrap gap-2">
                {option.values.map(value => {
                  const isSelected = selectedOptions[option.option_name] === value.option_value_id;
                  const hasImage = value.image;
                  // console.log('Option Value Image:', value);
                  return (
                    <Button
                      key={value.option_value_id}
                      variant={isSelected ? "default" : "outline"}
                      className={cn(
                        "transition-all",
                        isSelected && "ring-2 ring-primary ring-offset-2"
                      )}
                      onClick={() => handleOptionSelect(option.option_name, value.option_value_id, hasImage || undefined)}
                    >
                      {hasImage && (
                        <img 
                          src={getImageUrl(hasImage)} 
                          alt={value.value}
                          className="w-6 h-6 object-cover rounded mr-2"
                        />
                      )}
                      {/* convert to html */}
                      {/* ensure HTML content doesn't expand the button: wrap and break long words */}
                      <div
                        className="max-w-[220px] break-words overflow-hidden text-left text-sm leading-tight"
                        // eslint-disable-next-line react/no-danger
                        dangerouslySetInnerHTML={{ __html: value.value }}
                      />
                    </Button>
                  );
                })}
              </div>
            </div>
          ))}

          {/* Stock Info */}
          {selectedSku && (
            <div className="flex items-center gap-2 mb-4 text-sm">
              <span className="text-gray-600">{t("con_lai")}:</span>
              <span className="font-bold text-lg">
                {selectedSku.quantity} {t("san_pham")}
              </span>
            </div>
          )}

          {/* Quantity Selector */}
          <div className="flex items-center gap-4 mb-6">
            <span className="text-gray-700">{t("so_luong")}:</span>
            <div className="flex items-center border rounded">
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setQuantity(Math.max(1, quantity - 1))}
                disabled={!selectedSku}
              >
                <Minus className="w-4 h-4" />
              </Button>
              <input
                type="number"
                value={quantity}
                onChange={(e) => {
                  const val = parseInt(e.target.value) || 1;
                  setQuantity(Math.max(1, Math.min(selectedSku?.quantity || 1, val)));
                }}
                className="w-16 text-center border-x"
                disabled={!selectedSku}
              />
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setQuantity(Math.min(selectedSku?.quantity || 1, quantity + 1))}
                disabled={!selectedSku}
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 mb-6">
            <Button
              variant="outline"
              className="flex-1"
              onClick={handleAddToCart}
              disabled={!selectedSku || selectedSku.quantity === 0}
            >
              <ShoppingCart className="w-4 h-4 mr-2" />
              {t("them_vao_gio_hang")}
            </Button>
            <Button
              className="flex-1"
              onClick={handleBuyNow}
              disabled={!selectedSku || selectedSku.quantity === 0}
            >
              {t("mua_ngay")}
            </Button>
            <Button 
              variant="outline" 
              size="icon"
              onClick={handleAskAI}
              title={t("goi_y_tu_ai")}
            >
              <Heart className="w-4 h-4" />
            </Button>
          </div>

          {/* Additional Info */}
          <div className="space-y-2 text-sm text-gray-600 mb-6">
            {product?.product_is_permission_check && (
              <div className="flex items-center gap-2">
                <svg className="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd"/>
                </svg>
                <span>{t("san_pham_chinh_hang")}</span>
              </div>
            )}
            {product?.product_is_permission_return && (
              <div className="flex items-center gap-2">
                <svg className="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd"/>
                </svg>
                <span>{t("doi_tra_trong_7_ngay")}</span>
              </div>
            )}
          </div>

          {/* Description */}
          {product?.description && (

            <div className="border-t pt-6">
              <Button
                variant="ghost"
                className="w-full flex items-center justify-between p-4 hover:bg-gray-50"
                onClick={() => setIsDescriptionExpanded(!isDescriptionExpanded)}
              >
                <h2 className="text-xl font-bold">{t("mo_ta_san_pham")}</h2>
                {isDescriptionExpanded ? (
                  <ChevronUp className="w-5 h-5" />
                ) : (
                  <ChevronDown className="w-5 h-5" />
                )}
              </Button>
              
              {isDescriptionExpanded && (
                <>
                  <div 
                    className="prose max-w-none text-gray-700 mt-4 px-4"
                    dangerouslySetInnerHTML={{ __html: product.description }}
                  />
                  {/* Collapse button at bottom */}
                  <div className="flex justify-center pt-4 pb-2">
                    <Button
                      variant="outline"
                      onClick={() => setIsDescriptionExpanded(false)}
                      className="gap-2"
                    >
                      <ChevronUp className="w-4 h-4" />
                      {t("an_bot")}
                    </Button>
                  </div>
                </>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Product Comments Section */}
      <div className="mt-8">
        <ProductComments productId={data.result.data.product.id} />
      </div>

      {/* Image Gallery Modal */}
      <ImageGalleryModal
        images={images}
        initialIndex={galleryInitialIndex}
        isOpen={isGalleryOpen}
        onClose={() => setIsGalleryOpen(false)}
        notes={imageNotes}
      />
    </div>
  );
}
