"use client";

import { useState, useEffect, useRef } from "react";
import { useGetOrdersInfinite } from "@/services/apiService";
import { OrderStatus } from "@/types/order.types";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Loader2, Package, Truck, CheckCircle, XCircle, RefreshCcw, Clock, PackageCheck, Codesandbox } from "lucide-react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useRouter } from "@/i18n/routing";
import ROUTER from "@/assets/configs/routers";

const ORDER_TABS: { value: OrderStatus | 'ALL'; label: string; icon: any }[] = [
  { value: 'ALL', label: 'Tất cả đơn hàng', icon: Codesandbox  },
  { value: 'AWAITING_PAYMENT', label: 'Chờ thanh toán', icon: Clock },
  { value: 'PROCESSING', label: 'Đang xử lý', icon: PackageCheck },
  { value: 'SHIPPED', label: 'Đang giao hàng', icon: Truck },
  { value: 'COMPLETED', label: 'Hoàn thành', icon: CheckCircle },
  { value: 'CANCELED', label: 'Đã hủy', icon: XCircle },
  { value: 'REFUNDED', label: 'Trả hàng/Hoàn tiền', icon: RefreshCcw },
];

export default function OrdersPage() {
  const t = useTranslations("System");
  const router = useRouter();
  const [activeTab, setActiveTab] = useState<OrderStatus | 'ALL'>('ALL');
  const limit = 10;

  // Intersection Observer ref để phát hiện khi scroll đến cuối
  const loadMoreRef = useRef<HTMLDivElement>(null);

  // Fetch orders với infinite scroll
  const {
    data,
    isLoading,
    error,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useGetOrdersInfinite({
    limit,
    status: activeTab === 'ALL' ? undefined : activeTab,
  });

  // Setup Intersection Observer để tự động load thêm khi scroll xuống
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        const firstEntry = entries[0];
        // Khi phần tử loadMoreRef xuất hiện trong viewport và còn trang tiếp theo
        if (firstEntry.isIntersecting && hasNextPage && !isFetchingNextPage) {
          fetchNextPage();
        }
      },
      {
        threshold: 0.1, // Trigger khi 10% phần tử hiển thị
        rootMargin: '100px', // Load trước 100px
      }
    );

    if (loadMoreRef.current) {
      observer.observe(loadMoreRef.current);
    }

    return () => {
      if (loadMoreRef.current) {
        observer.unobserve(loadMoreRef.current);
      }
    };
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  // Flatten all orders from all pages
  const allOrders = data?.pages.flatMap((page) => page.data) || [];

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

  const formatDate = (dateString: string | null) => {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleDateString('vi-VN', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getStatusColor = (status: OrderStatus) => {
    switch (status) {
      case 'AWAITING_PAYMENT':
        return 'text-yellow-600 bg-yellow-50';
      case 'PROCESSING':
        return 'text-blue-600 bg-blue-50';
      case 'SHIPPED':
        return 'text-purple-600 bg-purple-50';
      case 'COMPLETED':
        return 'text-green-600 bg-green-50';
      case 'CANCELED':
        return 'text-red-600 bg-red-50';
      case 'REFUNDED':
        return 'text-orange-600 bg-orange-50';
      default:
        return 'text-gray-600 bg-gray-50';
    }
  };

  const getStatusText = (status: OrderStatus) => {
    const tab = ORDER_TABS.find(t => t.value === status);
    return tab?.label || status;
  };

  const handleTabChange = (value: string) => {
    setActiveTab(value as OrderStatus | 'ALL');
    // Reset về đầu khi đổi tab
  };

  // ==================== ACTION HANDLERS ====================
  
  /**
   * Xử lý thanh toán lại đơn hàng
   * @param orderId - ID của đơn hàng cần thanh toán lại
   */
  const handlePayAgain = (orderId: string) => {
    // TODO: Implement payment retry logic
    console.log('Thanh toán lại đơn hàng:', orderId);
  };

  /**
   * Xử lý cập nhật địa chỉ giao hàng
   * @param orderId - ID của đơn hàng cần cập nhật địa chỉ
   */
  const handleUpdateAddress = (orderId: string) => {
    // TODO: Implement address update logic
    console.log('Cập nhật địa chỉ cho đơn hàng:', orderId);
  };

  /**
   * Xử lý đánh giá sản phẩm
   * @param orderId - ID của đơn hàng cần đánh giá
   */
  const handleReview = (orderId: string) => {
    // TODO: Implement review logic
    console.log('Đánh giá đơn hàng:', orderId);
  };

  /**
   * Xử lý yêu cầu trả hàng
   * @param orderId - ID của đơn hàng cần trả
   */
  const handleReturn = (orderId: string) => {
    // TODO: Implement return logic
    console.log('Trả hàng cho đơn hàng:', orderId);
  };

  /**
   * Xử lý mua lại đơn hàng
   * @param orderId - ID của đơn hàng cần mua lại
   */
  const handleBuyAgain = (orderId: string) => {
    // TODO: Implement buy again logic
    console.log('Mua lại đơn hàng:', orderId);
  };

  /**
   * Xử lý liên hệ người bán
   * @param orderId - ID của đơn hàng
   * @param shopId - ID của shop
   */
  const handleContactSeller = (orderId: string, shopId: string) => {
    // TODO: Implement contact seller logic
    console.log('Liên hệ người bán cho đơn hàng:', orderId, 'Shop ID:', shopId);
  };

  /**
   * Xử lý xem chi tiết đơn hàng
   * @param orderId - ID của đơn hàng cần xem chi tiết
   */
  const handleViewDetail = (orderId: string) => {
    router.push(`${ROUTER.donhang}/${orderId}`);
  };

  /**
   * Xử lý xem shop
   * @param shopId - ID của shop
   */
  const handleViewShop = (shopId: string) => {
    // TODO: Implement view shop logic
    console.log('Xem shop:', shopId);
  };

  // ==================== RENDER ACTION BUTTONS ====================
  
  /**
   * Render các nút action dựa trên trạng thái đơn hàng
   */
  const renderOrderActions = (orderId: string, shopId: string, status: OrderStatus) => {
    const commonButtons = {
      contactSeller: (
        <Button 
          key="contact" 
          variant="outline" 
          size="sm" 
          className="text-blue-600 border-blue-600 hover:bg-blue-50"
          onClick={() => handleContactSeller(orderId, shopId)}
        >
          Liên Hệ Người Bán
        </Button>
      ),
      viewDetail: (
        <Button 
          key="detail" 
          variant="outline" 
          size="sm"
          onClick={() => handleViewDetail(orderId)}
        >
          Xem Chi Tiết
        </Button>
      ),
    };

    switch (status) {
      case 'AWAITING_PAYMENT':
        return (
          <>
            <Button 
              variant="default" 
              size="sm" 
              className="bg-orange-600 hover:bg-orange-700"
              onClick={() => handlePayAgain(orderId)}
            >
              Thanh Toán Lại
            </Button>
            {commonButtons.contactSeller}
            {commonButtons.viewDetail}
          </>
        );

      case 'PROCESSING':
        return (
          <>
            <Button 
              variant="outline" 
              size="sm" 
              className="text-purple-600 border-purple-600 hover:bg-purple-50"
              onClick={() => handleUpdateAddress(orderId)}
            >
              Cập Nhật Địa Chỉ
            </Button>
            {commonButtons.contactSeller}
            {commonButtons.viewDetail}
          </>
        );

      case 'SHIPPED':
        return (
          <>
            {commonButtons.contactSeller}
            {commonButtons.viewDetail}
          </>
        );

      case 'COMPLETED':
        return (
          <>
            <Button 
              variant="default" 
              size="sm" 
              className="bg-green-600 hover:bg-green-700"
              onClick={() => handleReview(orderId)}
            >
              Đánh Giá
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              className="text-red-600 border-red-600 hover:bg-red-50"
              onClick={() => handleReturn(orderId)}
            >
              Trả Hàng
            </Button>
            {commonButtons.contactSeller}
            {commonButtons.viewDetail}
          </>
        );

      case 'CANCELED':
      case 'REFUNDED':
        return (
          <>
            <Button 
              variant="outline" 
              size="sm" 
              className="text-orange-600 border-orange-600 hover:bg-orange-50"
              onClick={() => handleBuyAgain(orderId)}
            >
              Mua Lại
            </Button>
            {commonButtons.contactSeller}
            {commonButtons.viewDetail}
          </>
        );

      default:
        return (
          <>
            {commonButtons.contactSeller}
            {commonButtons.viewDetail}
          </>
        );
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6 text-[hsl(var(--primary))]">Đơn hàng của tôi</h1>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={handleTabChange} className="mb-6">
        <TabsList className="w-full grid grid-cols-7 h-auto gap-2 bg-transparent">
          {ORDER_TABS.map((tab) => {
            const Icon = tab.icon;
            return (
              <TabsTrigger
                key={tab.value}
                value={tab.value}
                className="flex flex-col items-center gap-2 py-3 data-[state=active]:bg-[hsl(var(--primary))] data-[state=active]:text-white border-2 data-[state=active]:border-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/0.1)]"
              >
                <Icon className="h-6 w-6" />
                <span className="text-xs font-medium">{tab.label}</span>
              </TabsTrigger>
            );
          })}
        </TabsList>
      </Tabs>

      {/* Loading State */}
      {isLoading && (
        <div className="flex justify-center items-center py-16">
          <Loader2 className="h-12 w-12 animate-spin text-[hsl(var(--primary))]" />
        </div>
      )}

      {/* Error State */}
      {error && (
        <div className="text-center py-16">
          <p className="text-red-500 text-lg">Có lỗi xảy ra khi tải đơn hàng</p>
        </div>
      )}

      {/* Orders List */}
      {!isLoading && !error && (
        <>
          {allOrders && allOrders.length > 0 ? (
            <div className="space-y-4">
              {allOrders.map((order) => (
                <Card key={order.shop_order_id} className="border-2 hover:shadow-lg transition-shadow">
                  <CardContent className="p-6">
                    {/* Order Header */}
                    <div className="flex justify-between items-start mb-4 pb-4 border-b">
                      <div className="flex items-center gap-4">
                        <div>
                          <p className="text-sm text-gray-500">Mã đơn hàng</p>
                          <p className="font-bold text-[hsl(var(--primary))]">{order.shop_order_code}</p>
                        </div>
                        <div className="h-8 w-px bg-gray-300" />
                        <div>
                          <p className="text-sm text-gray-500">Ngày đặt</p>
                          <p className="font-medium">{formatDate(order.created_at)}</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-3">
                        <span className={`px-3 py-1 rounded-full text-sm font-semibold ${getStatusColor(order.status)}`}>
                          {getStatusText(order.status)}
                        </span>
                        {order.paid_at && (
                          <span className="px-3 py-1 rounded-full text-sm font-semibold text-green-600 bg-green-50 border border-green-200">
                            ✓ Đã thanh toán
                          </span>
                        )}
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => handleViewShop(order.shop_id)}
                        >
                          Xem Shop
                        </Button>
                      </div>
                    </div>

                    {/* Order Items */}
                    <div className="space-y-3 mb-4">
                      {order.items.map((item: any) => (
                        <div key={item.item_id} className="flex gap-4 p-3 hover:bg-gray-50 rounded-lg transition-colors">
                          <img
                            src={getImageUrl(item.product_image)}
                            alt={item.product_name}
                            className="w-20 h-20 object-cover rounded-lg border"
                          />
                          <div className="flex-1">
                            <h3 className="font-medium text-gray-800 mb-1 line-clamp-2">
                              {item.product_name}
                            </h3>
                            <p className="text-sm text-gray-500 mb-1">{item.sku_attributes}</p>
                            <p className="text-sm text-gray-600">x{item.quantity}</p>
                          </div>
                          <div className="text-right">
                            {item.original_unit_price > item.final_unit_price && (
                              <p className="text-sm text-gray-400 line-through">
                                {formatPrice(item.original_unit_price)}
                              </p>
                            )}
                            <p className="font-semibold text-[hsl(var(--primary))]">
                              {formatPrice(item.final_unit_price)}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>

                    {/* Order Footer */}
                    <div className="flex justify-between items-center pt-4 border-t">
                      <div className="flex gap-2 flex-wrap">
                        {renderOrderActions(order.shop_order_id, order.shop_id, order.status)}
                      </div>
                      <div className="text-right">
                        <p className="text-sm text-gray-600 mb-1">
                          Thành tiền: <span className="text-xl font-bold text-[hsl(var(--primary))]">{formatPrice(order.total_amount)}</span>
                        </p>
                        {order.shipping_fee > 0 && (
                          <p className="text-xs text-gray-500">
                            (Phí vận chuyển: {formatPrice(order.shipping_fee)})
                          </p>
                        )}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}

              {/* Load More Trigger - Invisible element to trigger intersection observer */}
              <div ref={loadMoreRef} className="flex justify-center py-4">
                {isFetchingNextPage && (
                  <div className="flex items-center gap-2 text-[hsl(var(--primary))]">
                    <Loader2 className="h-6 w-6 animate-spin" />
                    <span className="text-sm font-medium">Đang tải thêm đơn hàng...</span>
                  </div>
                )}
                {!hasNextPage && allOrders.length > 0 && (
                  <p className="text-sm text-gray-500">Đã hiển thị tất cả đơn hàng</p>
                )}
              </div>
            </div>
          ) : (
            <div className="text-center py-16">
              <Package className="h-24 w-24 mx-auto text-gray-300 mb-4" />
              <p className="text-gray-500 text-lg">Chưa có đơn hàng nào</p>
              <Link href="/">
                <Button className="mt-4 bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/0.9)]">
                  Tiếp tục mua sắm
                </Button>
              </Link>
            </div>
          )}
        </>
      )}
    </div>
  );
}
