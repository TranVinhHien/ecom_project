"use client";

import { useEffect, useState, useCallback, useRef } from "react";
import { useCartStore } from "@/store/cartStore";
import { useCheckoutStore } from "@/store/checkoutStore";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Minus, Plus, Trash2, ShoppingBag, Loader2 } from "lucide-react";
import { Link } from '@/i18n/routing';
import { useRouter } from "@/i18n/routing";
import ROUTER from "@/assets/configs/routers";
import { INFO_USER } from "@/assets/configs/request";
import { 
  useGetCart, 
  useUpdateCartItem, 
  useDeleteCartItem 
} from "@/services/apiService";
import { ApiCartItem } from "@/types/cart.types";
import { useToast } from "@/hooks/use-toast";
import { Checkbox } from "@/components/ui/checkbox";

export default function CartPage() {
  const router = useRouter();
  const { toast } = useToast();
  const [isHydrated, setIsHydrated] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  
  // Debounce timers for API updates
  const debounceTimers = useRef<Record<string, NodeJS.Timeout>>({});
  
  // Pending updates (lưu tạm các thay đổi chưa sync lên API)
  const [pendingUpdates, setPendingUpdates] = useState<Record<string, number>>({});
  
  // Track items being deleted (để ẩn ngay lập tức khỏi UI)
  const [deletingItems, setDeletingItems] = useState<Set<string>>(new Set());
  
  // Track selection state for API cart items (local state only)
  const [apiCartSelection, setApiCartSelection] = useState<Record<string, boolean>>({});
  
  // LocalStorage cart (for non-logged in users)
  const localItems = useCartStore((state) => state.items);
  const removeFromLocalCart = useCartStore((state) => state.removeFromCart);
  const updateLocalQuantity = useCartStore((state) => state.updateQuantity);
  const toggleLocalSelection = useCartStore((state) => state.toggleSelection);
  const getLocalSelectedTotalPrice = useCartStore((state) => state.getSelectedTotalPrice);
  const getSelectedLocalItems = useCartStore((state) => state.getSelectedItems);
  const { setCheckoutItems } = useCheckoutStore();

  // API cart (for logged in users)
  const { data: apiCart, isLoading: apiLoading, error: apiError } = useGetCart();
  const updateCartMutation = useUpdateCartItem();
  const deleteCartMutation = useDeleteCartItem();
  console.log("🛒 API Cart Data:", apiCart);
  // Check if user is logged in
  useEffect(() => {
    useCartStore.persist.rehydrate();
    const userInfo = localStorage.getItem(INFO_USER);
    setIsLoggedIn(!!userInfo);
    setIsHydrated(true);
  }, []);

  const getImageUrl = (imagePath: string | null | undefined) => {
    if (!imagePath) return '/placeholder.png';
    if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
      return imagePath;
    }
    return `http://${imagePath}`;
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  // Handle quantity update for API cart with debounce (10 seconds)
  const handleApiUpdateQuantity = useCallback((skuId: string, newQuantity: number) => {
    if (newQuantity < 1) return;
    
    // Cập nhật UI ngay lập tức (optimistic update)
    setPendingUpdates(prev => ({ ...prev, [skuId]: newQuantity }));
    
    // Clear existing timer
    if (debounceTimers.current[skuId]) {
      clearTimeout(debounceTimers.current[skuId]);
    }
    
    // Set new timer - sau 10 giây mới gửi API
    debounceTimers.current[skuId] = setTimeout(async () => {
      try {
        await updateCartMutation.mutateAsync({
          skuId,
          payload: { Quantity: newQuantity }
        });
        
        // Xóa pending update sau khi thành công
        // setPendingUpdates(prev => {
        //   const newPending = { ...prev };
        //   delete newPending[skuId];
        //   return newPending;
        // });
        
        toast({
          title: "Đã cập nhật",
          description: "Số lượng sản phẩm đã được cập nhật",
        });
      } catch (error) {
        // Nếu lỗi, xóa pending update để hiển thị lại giá trị cũ
        setPendingUpdates(prev => {
          const newPending = { ...prev };
          delete newPending[skuId];
          return newPending;
        });
        
        toast({
          title: "Lỗi",
          description: "Không thể cập nhật số lượng sản phẩm",
          variant: "destructive",
        });
      }
    }, 2000); // 10 seconds debounce
  }, [updateCartMutation, toast]);

  // Handle remove item from API cart
  const handleApiRemoveItem = async (skuId: string) => {
    // Ẩn item ngay lập tức (optimistic update)
    setDeletingItems(prev => new Set(prev).add(skuId));
    
    try {
      await deleteCartMutation.mutateAsync(skuId);
      toast({
        title: "Đã xóa sản phẩm",
        description: "Sản phẩm đã được xóa khỏi giỏ hàng",
      });
    } catch (error) {
      // Nếu lỗi, hiện lại item
      setDeletingItems(prev => {
        const newSet = new Set(prev);
        newSet.delete(skuId);
        return newSet;
      });
      toast({
        title: "Lỗi",
        description: "Không thể xóa sản phẩm",
        variant: "destructive",
      });
    }
  };

  // Handle checkout
  const handleCheckout = () => {
    if (isLoggedIn && apiCart) {
      // Checkout with API cart data - CHỈ LẤY CÁC SẢN PHẨM ĐÃ CHỌN
      const selectedItems = getApiSelectedItems();
      
      if (selectedItems.length === 0) {
        toast({
          title: "Thông báo",
          description: "Vui lòng chọn ít nhất một sản phẩm để thanh toán",
          variant: "destructive",
        });
        return;
      }
      
      const checkoutItems = selectedItems.map(item => ({
        sku_id: item.skuId,
        shop_id: item.shopId,
        quantity: item.quantity,
        name: item.productName,
        price: item.price,
        image: item.productImage || '',
        sku_name: item.productName
      }));
      
      // ✅ ĐÁNH DẤU: Items này TỪ GIỎ HÀNG (sẽ xóa sau khi thanh toán)
      setCheckoutItems(checkoutItems, true);
    } else {
      // Checkout with localStorage cart data - CHỈ LẤY CÁC SẢN PHẨM ĐÃ CHỌN
      const selectedItems = getSelectedLocalItems();
      
      if (selectedItems.length === 0) {
        toast({
          title: "Thông báo",
          description: "Vui lòng chọn ít nhất một sản phẩm để thanh toán",
          variant: "destructive",
        });
        return;
      }
      
      const checkoutItems = selectedItems.map(item => ({
        sku_id: item.sku_id,
        shop_id: item.shop_id,
        quantity: item.quantity,
        name: item.name,
        price: item.price,
        image: item.image,
        sku_name: item.sku_name
      }));
      
      // ✅ ĐÁNH DẤU: Items này TỪ GIỎ HÀNG (sẽ xóa sau khi thanh toán)
      setCheckoutItems(checkoutItems, true);
    }
    router.push(ROUTER.thanhtoan);
  };

  // Cleanup debounce timers on unmount
  useEffect(() => {
    return () => {
      Object.values(debounceTimers.current).forEach(timer => clearTimeout(timer));
    };
  }, []);

  // Handle selection for API cart
  const handleApiToggleSelection = async (skuId: string) => {
    setApiCartSelection(prev => ({
      ...prev,
      [skuId]: !prev[skuId]
    }));
  };

  // Handle selection for localStorage cart
  const handleLocalToggleSelection = (skuId: string) => {
    toggleLocalSelection(skuId);
  };

  // Select all items
  const handleSelectAll = () => {
    if (isLoggedIn && apiCart) {
      // Select all API cart items
      const allSelected = apiCart.items.every(item => apiCartSelection[item.skuId]);
      const newSelection: Record<string, boolean> = {};
      apiCart.items.forEach(item => {
        newSelection[item.skuId] = !allSelected;
      });
      setApiCartSelection(newSelection);
    } else {
      // For localStorage, select all items
      localItems.forEach(item => {
        if (!item.isSelected) {
          toggleLocalSelection(item.sku_id);
        }
      });
    }
  };

  // Get selected items for API cart
  const getApiSelectedItems = () => {
    if (!apiCart) return [];
    return apiCart.items.filter(item => apiCartSelection[item.skuId]);
  };

  // Calculate selected total price for API cart
  const getApiSelectedTotalPrice = () => {
    return getApiSelectedItems()
      .filter(item => !deletingItems.has(item.skuId)) // Loại bỏ items đang xóa
      .reduce((sum, item) => {
        // Dùng pending quantity nếu có, không thì dùng quantity từ API
        const quantity = pendingUpdates[item.skuId] ?? item.quantity;
        return sum + (item.price * quantity);
      }, 0);
  };

  // Show loading during hydration
  if (!isHydrated) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="flex items-center justify-center">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      </div>
    );
  }

  // Show loading when fetching API cart
  if (isLoggedIn && apiLoading) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="flex items-center justify-center">
          <Loader2 className="h-8 w-8 animate-spin" />
          <span className="ml-2">Đang tải giỏ hàng...</span>
        </div>
      </div>
    );
  }

  // Determine which cart to display
  const displayItems = isLoggedIn ? (apiCart?.items || []) : localItems;
  
  // Calculate total price for SELECTED items only
  const totalPrice = isLoggedIn 
    ? getApiSelectedTotalPrice()
    : getLocalSelectedTotalPrice();
    
  const isEmpty = displayItems.length === 0;

  if (isEmpty) {
    return (
      <div className="container mx-auto px-4 py-16">
        <Card className="max-w-2xl mx-auto text-center">
          <CardContent className="pt-12 pb-8">
            <ShoppingBag className="w-24 h-24 mx-auto text-gray-300 mb-6" />
            <h2 className="text-2xl font-bold mb-4">Giỏ hàng trống</h2>
            <p className="text-gray-600 mb-8">
              Bạn chưa có sản phẩm nào trong giỏ hàng
            </p>
            <Link href="/">
              <Button className="bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)]">
                Tiếp tục mua sắm
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Giỏ hàng của bạn</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left: Cart Items */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Sản phẩm ({displayItems.length})</CardTitle>
                {displayItems.length > 0 && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleSelectAll}
                  >
                    Chọn tất cả
                  </Button>
                )}
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              {isLoggedIn ? (
                // Display API cart items
                apiCart?.items
                  .filter(item => !deletingItems.has(item.skuId)) // Ẩn items đang bị xóa
                  .map((item: ApiCartItem) => {
                    // Lấy quantity từ pendingUpdates hoặc từ API data
                    const displayQuantity = pendingUpdates[item.skuId] ?? item.quantity;
                    
                    return (
                  <div key={item.skuId}>
                    <div className="flex gap-4">
                      {/* Checkbox - Using local selection state */}
                      <div className="flex items-start pt-6">
                        <Checkbox
                          checked={apiCartSelection[item.skuId] || false}
                          onCheckedChange={() => handleApiToggleSelection(item.skuId)}
                          className="mt-1"
                        />
                      </div>
                      {/* Image placeholder - API doesn't return image */}
                      <div className="relative w-24 h-24 flex-shrink-0 bg-gray-100 rounded flex items-center justify-center">
                         <div className="relative w-24 h-24 flex-shrink-0">
                        <img
                          src={getImageUrl(item.productImage) || "/placeholder.png"}
                          alt={item.productImage}
                          className="object-cover rounded"
                        />
                      </div>
                      </div>

                      {/* Info */}
                      <div className="flex-1 min-w-0">
                        <h4 className="font-medium mb-1 line-clamp-2">
                          {item.productName}
                        </h4>
                        {/* <p className="text-sm text-gray-500 mb-2">
                          SKU: {item.skuId}
                        </p> */}
                        <p className="text-lg font-bold text-[hsl(var(--primary))]">
                          {formatPrice(item.price)}
                        </p>
                      </div>

                      {/* Quantity Controls */}
                      <div className="flex flex-col items-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleApiRemoveItem(item.skuId)}
                          disabled={deleteCartMutation.isPending}
                          className="text-red-500 hover:text-red-700"
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>

                        <div className="flex items-center gap-2">
                          <Button
                            variant="outline"
                            size="icon"
                            className="h-8 w-8"
                            onClick={() => handleApiUpdateQuantity(item.skuId, displayQuantity - 1)}
                            disabled={displayQuantity <= 1}
                          >
                            <Minus className="w-3 h-3" />
                          </Button>

                          <span className="w-12 text-center font-medium">
                            {displayQuantity}
                          </span>

                          <Button
                            variant="outline"
                            size="icon"
                            className="h-8 w-8"
                            onClick={() => handleApiUpdateQuantity(item.skuId, displayQuantity + 1)}
                          >
                            <Plus className="w-3 h-3" />
                          </Button>
                        </div>

                        <p className="text-sm text-gray-500">
                          Tổng: {formatPrice(item.price * displayQuantity)}
                        </p>
                      </div>
                    </div>
                    <Separator className="mt-4" />
                  </div>
                    );
                  })
              ) : (
                // Display localStorage cart items
                localItems.map((item) => (
                  <div key={item.sku_id}>
                    <div className="flex gap-4">
                      {/* Checkbox */}
                      <div className="flex items-start pt-6">
                        <Checkbox
                          checked={item.isSelected || false}
                          onCheckedChange={() => handleLocalToggleSelection(item.sku_id)}
                          className="mt-1"
                        />
                      </div>
                      {/* Image */}
                      <div className="relative w-24 h-24 flex-shrink-0">
                        <img
                          src={getImageUrl(item.image) || "/placeholder.png"}
                          alt={item.name}
                          className="object-cover rounded"
                        />
                      </div>

                      {/* Info */}
                      <div className="flex-1 min-w-0">
                        <h4 className="font-medium mb-1 line-clamp-2">
                          {item.name}
                        </h4>
                        <p className="text-sm text-gray-500 mb-2">
                          SKU: {item.sku_id}
                        </p>
                        <p className="text-lg font-bold text-[hsl(var(--primary))]">
                          {formatPrice(item.price)}
                        </p>
                      </div>

                      {/* Quantity Controls */}
                      <div className="flex flex-col items-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => removeFromLocalCart(item.sku_id)}
                          className="text-red-500 hover:text-red-700"
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>

                        <div className="flex items-center gap-2">
                          <Button
                            variant="outline"
                            size="icon"
                            className="h-8 w-8"
                            onClick={() =>
                              updateLocalQuantity(item.sku_id, Math.max(1, item.quantity - 1))
                            }
                          >
                            <Minus className="w-3 h-3" />
                          </Button>

                          <span className="w-12 text-center font-medium">
                            {item.quantity}
                          </span>

                          <Button
                            variant="outline"
                            size="icon"
                            className="h-8 w-8"
                            onClick={() =>
                              updateLocalQuantity(item.sku_id, item.quantity + 1)
                            }
                          >
                            <Plus className="w-3 h-3" />
                          </Button>
                        </div>

                        <p className="text-sm text-gray-500">
                          Tổng: {formatPrice(item.price * item.quantity)}
                        </p>
                      </div>
                    </div>
                    <Separator className="mt-4" />
                  </div>
                ))
              )}
            </CardContent>
          </Card>
        </div>

        {/* Right: Order Summary */}
        <div className="lg:col-span-1">
          <Card className="sticky top-4">
            <CardHeader>
              <CardTitle>Tổng đơn hàng</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Tạm tính (chưa tính tiền ship)</span>
                  <span>{formatPrice(totalPrice)}</span>
                </div>
                <Separator />
                <div className="flex justify-between text-lg font-bold">
                  <span>Tổng cộng</span>
                  <span className="text-[hsl(var(--primary))]">
                    {formatPrice(totalPrice)}
                  </span>
                </div>
              </div>

              <Button
                className="w-full bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)]"
                size="lg"
                onClick={handleCheckout}
              >
                Tiến hành thanh toán
              </Button>

              <Link href="/">
                <Button variant="outline" className="w-full" size="lg">
                  Tiếp tục mua sắm
                </Button>
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}   