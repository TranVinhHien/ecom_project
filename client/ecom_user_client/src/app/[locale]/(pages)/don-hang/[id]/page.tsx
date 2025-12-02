"use client";

import { useGetOrderDetail } from "@/services/apiService";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Loader2, ArrowLeft, Package, MapPin, CreditCard, FileText, Truck, CheckCircle } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { ShippingAddress } from "@/types/order.types";

export default function OrderDetailPage({ params }: { params: { id: string } }) {
  const t = useTranslations("System");
  const router = useRouter();
  const { data: orderDetail, isLoading, error } = useGetOrderDetail(params.id);

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
  
  // Parse shipping address từ JSON string hoặc object
  const parseShippingAddress = (addressData: string | ShippingAddress | any): ShippingAddress | null => {
    try {
      // Nếu đã là object với đầy đủ properties, trả về luôn
      if (typeof addressData === 'object' && addressData !== null) {
        if (addressData.fullName || addressData.address) {
          return addressData as ShippingAddress;
        }
      }
      // Nếu là JSON string, parse nó
      if (typeof addressData === 'string') {
        const parsed = JSON.parse(addressData);
        return parsed as ShippingAddress;
      }
      return null;
    } catch (e) {
      console.error('Error parsing shipping address:', e);
      return null;
    }
  };
  
  // Parse payment method từ JSON string, object, hoặc ID
  const parsePaymentMethod = (paymentData: string | any): { name: string; code: string; type: string } => {
    try {
      // Nếu đã là object với properties
      if (typeof paymentData === 'object' && paymentData !== null) {
        return {
          name: paymentData.name || 'N/A',
          code: paymentData.code || '',
          type: paymentData.type || ''
        };
      }
      // Nếu là JSON string, parse nó
      if (typeof paymentData === 'string') {
        try {
          const parsed = JSON.parse(paymentData);
          if (typeof parsed === 'object' && parsed !== null) {
            return {
              name: parsed.name || 'N/A',
              code: parsed.code || '',
              type: parsed.type || ''
            };
          }
          // Nếu parse ra string thì coi như tên phương thức
          return { name: String(parsed), code: '', type: '' };
        } catch {
          // Nếu không parse được, coi như là tên phương thức
          return { name: paymentData, code: '', type: '' };
        }
      }
      return { name: 'N/A', code: '', type: '' };
    } catch (e) {
      console.error('Error parsing payment method:', e);
      return { name: 'N/A', code: '', type: '' };
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'AWAITING_PAYMENT':
        return 'text-yellow-600 bg-yellow-50 border-yellow-200';
      case 'PROCESSING':
        return 'text-blue-600 bg-blue-50 border-blue-200';
      case 'SHIPPED':
        return 'text-purple-600 bg-purple-50 border-purple-200';
      case 'COMPLETED':
        return 'text-green-600 bg-green-50 border-green-200';
      case 'CANCELED':
        return 'text-red-600 bg-red-50 border-red-200';
      case 'REFUNDED':
        return 'text-orange-600 bg-orange-50 border-orange-200';
      default:
        return 'text-gray-600 bg-gray-50 border-gray-200';
    }
  };

  const getStatusText = (status: string) => {
    const statusMap: Record<string, string> = {
      'AWAITING_PAYMENT': 'Chờ thanh toán',
      'PROCESSING': 'Đang xử lý',
      'SHIPPED': 'Đang giao hàng',
      'COMPLETED': 'Hoàn thành',
      'CANCELED': 'Đã hủy',
      'REFUNDED': 'Hoàn tiền',
    };
    return statusMap[status] || status;
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <Loader2 className="h-12 w-12 animate-spin text-[hsl(var(--primary))]" />
      </div>
    );
  }

  if (error || !orderDetail) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center py-16">
          <Package className="h-24 w-24 mx-auto text-gray-300 mb-4" />
          <p className="text-red-500 text-lg mb-4">Không tìm thấy đơn hàng</p>
          <Link href="/don-hang">
            <Button className="bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/0.9)]">
              Quay lại danh sách đơn hàng
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  const { order, order_shop } = orderDetail;
  
  // Parse shipping address và payment method
  const shippingAddress = parseShippingAddress(order.shipping_address);
  const paymentMethod = parsePaymentMethod(order.payment_method);

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header with Back Button */}
      <div className="flex items-center gap-4 mb-6">
        <Button
          variant="outline"
          size="icon"
          onClick={() => router.back()}
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div className="flex-1">
          <h1 className="text-3xl font-bold text-[hsl(var(--primary))]">Chi tiết đơn hàng</h1>
          <p className="text-gray-600">Mã đơn hàng: <span className="font-semibold">{order.order_code}</span></p>
        </div>
        <div className="flex items-center gap-2">
          <span className={`px-4 py-2 rounded-full text-sm font-semibold border ${getStatusColor(order_shop.status)}`}>
            {getStatusText(order_shop.status)}
          </span>
          {order_shop.paid_at && (
            <span className="px-4 py-2 rounded-full text-sm font-semibold text-green-600 bg-green-50 border border-green-200">
              ✓ Đã thanh toán
            </span>
          )}
        </div>
      </div>

      {/* Order Timeline */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Truck className="h-5 w-5" />
            Trạng thái đơn hàng
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex justify-between items-center">
            {/* Timeline Steps */}
            <div className="flex items-center gap-4 w-full">
              {/* Đơn hàng đã đặt */}
              <div className="flex flex-col items-center flex-1">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                  order_shop.created_at ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-500'
                }`}>
                  <Package className="h-6 w-6" />
                </div>
                <p className="text-xs mt-2 text-center font-medium">Đơn Hàng Đã Đặt</p>
                <p className="text-xs text-gray-500">{formatDate(order_shop.created_at)}</p>
              </div>

              <div className="flex-1 h-1 bg-gray-200 relative">
                <div className={`h-full ${order_shop.paid_at ? 'bg-green-500' : 'bg-gray-200'}`} />
              </div>

              {/* Đã thanh toán */}
              <div className="flex flex-col items-center flex-1">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                  order_shop.paid_at ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-500'
                }`}>
                  <CreditCard className="h-6 w-6" />
                </div>
                <p className="text-xs mt-2 text-center font-medium">Đã Thanh Toán</p>
                <p className="text-xs text-gray-500">{formatDate(order_shop.paid_at)}</p>
              </div>

              <div className="flex-1 h-1 bg-gray-200 relative">
                <div className={`h-full ${order_shop.processing_at ? 'bg-green-500' : 'bg-gray-200'}`} />
              </div>

              {/* Đang xử lý */}
              <div className="flex flex-col items-center flex-1">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                  order_shop.processing_at ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-500'
                }`}>
                  <Package className="h-6 w-6" />
                </div>
                <p className="text-xs mt-2 text-center font-medium">Đang Xử Lý</p>
                <p className="text-xs text-gray-500">{formatDate(order_shop.processing_at)}</p>
              </div>

              <div className="flex-1 h-1 bg-gray-200 relative">
                <div className={`h-full ${order_shop.shipped_at ? 'bg-green-500' : 'bg-gray-200'}`} />
              </div>

              {/* Đang giao */}
              <div className="flex flex-col items-center flex-1">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                  order_shop.shipped_at ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-500'
                }`}>
                  <Truck className="h-6 w-6" />
                </div>
                <p className="text-xs mt-2 text-center font-medium">Đang giao Hàng</p>
                <p className="text-xs text-gray-500">{formatDate(order_shop.shipped_at)}</p>
              </div>

              <div className="flex-1 h-1 bg-gray-200 relative">
                <div className={`h-full ${order_shop.completed_at ? 'bg-green-500' : 'bg-gray-200'}`} />
              </div>

              {/* Hoàn thành */}
              <div className="flex flex-col items-center flex-1">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                  order_shop.completed_at ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-500'
                }`}>
                  <CheckCircle className="h-6 w-6" />
                </div>
                <p className="text-xs mt-2 text-center font-medium">Hoàn Thành</p>
                <p className="text-xs text-gray-500">{formatDate(order_shop.completed_at)}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column */}
        <div className="lg:col-span-2 space-y-6">
          {/* Shipping Address */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <MapPin className="h-5 w-5" />
                Địa chỉ nhận hàng
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <p className="font-semibold text-lg">{shippingAddress?.fullName || 'N/A'}</p>
                <p className="text-gray-600">Điện thoại: <span className="font-medium">{shippingAddress?.phone || 'N/A'}</span></p>
                <p className="text-gray-600">
                  Địa chỉ: <span className="font-medium">
                    {shippingAddress ? `${shippingAddress.address}, ${shippingAddress.district}, ${shippingAddress.city}, ${shippingAddress.postalCode}` : 'N/A'}
                  </span>
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Order Items */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Package className="h-5 w-5" />
                Sản phẩm ({order_shop.items.length})
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {order_shop.items.map((item:any) => (
                  <div key={item.item_id} className="flex gap-4 p-4 border rounded-lg hover:bg-gray-50 transition-colors">
                    <img
                      src={getImageUrl(item.product_image)}
                      alt={item.product_name}
                      className="w-24 h-24 object-cover rounded-lg border"
                    />
                    <div className="flex-1">
                      <h3 className="font-medium text-gray-800 mb-2 line-clamp-2">
                        {item.product_name}
                      </h3>
                      {item.sku_attributes && (
                        <p className="text-sm text-gray-500 mb-2">{item.sku_attributes}</p>
                      )}
                      <p className="text-sm text-gray-600">Số lượng: x{item.quantity}</p>
                    </div>
                    <div className="text-right">
                      {item.original_unit_price > item.final_unit_price && (
                        <p className="text-sm text-gray-400 line-through">
                          {formatPrice(item.original_unit_price)}
                        </p>
                      )}
                      <p className="font-semibold text-lg text-[hsl(var(--primary))]">
                        {formatPrice(item.final_unit_price)}
                      </p>
                      <p className="text-sm text-gray-500 mt-1">
                        Tổng: {formatPrice(item.total_price)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Order Note */}
          {order.note && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  Ghi chú
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-700">{order.note}</p>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Right Column - Order Summary */}
        <div className="space-y-6">
          {/* Payment Method */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <CreditCard className="h-5 w-5" />
                Phương thức thanh toán
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <p className="font-semibold text-lg">{paymentMethod?.name || 'N/A'}</p>
                <p className="text-sm text-gray-600">Loại: <span className="font-medium">{paymentMethod?.type || 'N/A'}</span></p>
                {paymentMethod?.code && <p className="text-sm text-gray-600">Mã: <span className="font-medium">{paymentMethod.code}</span></p>}
              </div>
            </CardContent>
          </Card>

          {/* Price Summary */}
          <Card>
            <CardHeader>
              <CardTitle>Tổng quan đơn hàng</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {/* Subtotal */}
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Tổng tiền hàng:</span>
                <span className="font-medium">{formatPrice(order_shop.subtotal)}</span>
              </div>
              
              {/* Shipping Fee */}
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Phí vận chuyển:</span>
                <span className="font-medium">{formatPrice(order_shop.shipping_fee)}</span>
              </div>

              {/* Vouchers Section */}
              {(order_shop.shop_voucher_discount > 0 || order_shop.site_order_discount > 0 || order_shop.site_shipping_discount > 0) && (
                <div className="border-t border-b py-3 space-y-2">
                  <p className="text-xs font-semibold text-gray-700 uppercase mb-2">Mã giảm giá</p>
                  
                  {/* Shop Voucher */}
                  {order_shop.shop_voucher_discount > 0 && (
                    <div className="flex justify-between text-sm">
                      <div className="flex flex-col">
                        <span className="text-gray-600">Voucher Shop</span>
                        {order.site_order_voucher_code && (
                          <span className="text-xs text-green-600 font-medium">({order.site_order_voucher_code})</span>
                        )}
                      </div>
                      <span className="font-medium text-green-600">-{formatPrice(order_shop.shop_voucher_discount)}</span>
                    </div>
                  )}

                  {/* Site Order Voucher */}
                  {order_shop.site_order_discount > 0 && (
                    <div className="flex justify-between text-sm">
                      <div className="flex flex-col">
                        <span className="text-gray-600">Voucher Sàn (Đơn hàng)</span>
                        {order.site_order_voucher_code && (
                          <span className="text-xs text-green-600 font-medium">({order.site_order_voucher_code})</span>
                        )}
                      </div>
                      <span className="font-medium text-green-600">-{formatPrice(order_shop.site_order_discount)}</span>
                    </div>
                  )}

                  {/* Site Shipping Voucher */}
                  {order_shop.site_shipping_discount > 0 && (
                    <div className="flex justify-between text-sm">
                      <div className="flex flex-col">
                        <span className="text-gray-600">Voucher Sàn (Vận chuyển)</span>
                        {order.site_shipping_voucher_code && (
                          <span className="text-xs text-green-600 font-medium">({order.site_shipping_voucher_code})</span>
                        )}
                      </div>
                      <span className="font-medium text-green-600">-{formatPrice(order_shop.site_shipping_discount)}</span>
                    </div>
                  )}

                  {/* Total Discount */}
                  {order_shop.total_discount > 0 && (
                    <div className="flex justify-between text-sm font-semibold pt-2 border-t">
                      <span className="text-gray-700">Tổng giảm giá:</span>
                      <span className="text-green-600">-{formatPrice(order_shop.total_discount)}</span>
                    </div>
                  )}
                </div>
              )}

              {/* Grand Total */}
              <div className="border-t-2 pt-3">
                <div className="flex justify-between items-center">
                  <span className="font-semibold text-lg">Tổng cộng:</span>
                  <span className="font-bold text-2xl text-[hsl(var(--primary))]">
                    {formatPrice(order_shop.total_amount-order_shop.total_discount)}
                  </span>
                </div>
              </div>

              {/* Grand Total from Order (if different) */}
              {order.grand_total !== order_shop.total_amount && (
                <div className="bg-blue-50 p-3 rounded-lg border border-blue-200">
                  <div className="flex justify-between items-center text-sm">
                    <span className="text-blue-700 font-medium">Tổng thanh toán (Tổng đơn):</span>
                    <span className="font-bold text-blue-700">{formatPrice(order.grand_total)}</span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Order Info */}
          <Card>
            <CardHeader>
              <CardTitle>Thông tin đơn hàng</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Mã đơn hàng:</span>
                <span className="font-medium">{order.order_code}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Mã đơn shop:</span>
                <span className="font-medium">{order_shop.shop_order_code}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Shop ID:</span>
                <span className="font-medium">{order_shop.shop_id}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Ngày đặt:</span>
                <span className="font-medium">{formatDate(order.created_at)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Cập nhật cuối:</span>
                <span className="font-medium">{formatDate(order.updated_at)}</span>
              </div>
              {order_shop.tracking_code && (
                <div className="flex justify-between">
                  <span className="text-gray-600">Mã vận đơn:</span>
                  <span className="font-medium">{order_shop.tracking_code}</span>
                </div>
              )}
              {order_shop.shipping_method && (
                <div className="flex justify-between">
                  <span className="text-gray-600">Đơn vị vận chuyển:</span>
                  <span className="font-medium">{order_shop.shipping_method}</span>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Action Buttons */}
          <div className="space-y-2">
            <Button className="w-full bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/0.9)]">
              Liên Hệ Người Bán
            </Button>
            <Button variant="outline" className="w-full" onClick={() => router.back()}>
              Quay Lại
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
