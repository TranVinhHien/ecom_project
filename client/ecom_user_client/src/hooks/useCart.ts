/**
 * Custom hook for cart operations
 * 
 * This hook automatically handles cart operations for both:
 * - Logged in users (via API)
 * - Non-logged in users (via localStorage)
 */

import { useState, useEffect } from "react";
import { useCartStore } from "@/store/cartStore";
import { useAddToCart, useGetCartCount } from "@/services/apiService";
import { INFO_USER } from "@/assets/configs/request";
import { CartItem } from "@/types/cart.types";
import { useToast } from "./use-toast";

export const useCart = () => {
  const { toast } = useToast();
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isHydrated, setIsHydrated] = useState(false);

  // LocalStorage cart operations
  const addToLocalCart = useCartStore((state) => state.addToCart);
  const localCartCount = useCartStore((state) => state.getTotalItems());

  // API cart operations
  const addToApiCartMutation = useAddToCart();
  const { data: apiCartCount, refetch: refetchCartCount } = useGetCartCount();

  // Check login status
  useEffect(() => {
    useCartStore.persist.rehydrate();
    const userInfo = localStorage.getItem(INFO_USER);
    setIsLoggedIn(!!userInfo);
    setIsHydrated(true);
  }, []);

  /**
   * Add item to cart (automatically chooses API or localStorage)
   */
  const addToCart = async (item: CartItem): Promise<boolean> => {
    try {
      if (isLoggedIn) {
        // User is logged in - use API
        await addToApiCartMutation.mutateAsync({
          SkuId: item.sku_id,
          Quantity: item.quantity,
        });

        toast({
          title: "Thành công",
          description: "Đã thêm sản phẩm vào giỏ hàng",
        });
        
        // Immediately invalidate cart count query để cập nhật header ngay lập tức
        refetchCartCount();
        
        return true;
      } else {
        // User is not logged in - use localStorage
        addToLocalCart(item);
        
        toast({
          title: "Thành công",
          description: "Đã thêm sản phẩm vào giỏ hàng tạm thời",
        });
        
        return true;
      }
    } catch (error) {
      console.error("Error adding to cart:", error);
      
      toast({
        title: "Lỗi",
        description: "Không thể thêm sản phẩm vào giỏ hàng",
        variant: "destructive",
      });
      
      return false;
    }
  };

  /**
   * Get current cart count
   */
  const getCartCount = (): number => {
    if (!isHydrated) return 0;
    return isLoggedIn ? (apiCartCount || 0) : localCartCount;
  };

  return {
    addToCart,
    getCartCount,
    isLoggedIn,
    isHydrated,
    isAddingToCart: addToApiCartMutation.isPending,
  };
};
