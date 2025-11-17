"use client";

import { useEffect, useState } from "react";
import { useCartStore } from "@/store/cartStore";
import { useCheckoutStore } from "@/store/checkoutStore";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Minus, Plus, Trash2, ShoppingBag, Loader2 } from "lucide-react";
import { Link } from '@/i18n/routing';
import { useRouter } from "@/i18n/routing"
import ROUTER from "@/assets/configs/routers";

export default function CartPage() {
  const router = useRouter();
  const [isHydrated, setIsHydrated] = useState(false);
  
  const items = useCartStore((state) => state.items);
  const removeFromCart = useCartStore((state) => state.removeFromCart);
  const updateQuantity = useCartStore((state) => state.updateQuantity);
  const getTotalPrice = useCartStore((state) => state.getTotalPrice);
  const { setCheckoutItems } = useCheckoutStore();

  // Manually hydrate from localStorage
  useEffect(() => {
    useCartStore.persist.rehydrate();
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

  const handleCheckout = () => {
    // Convert cart items to checkout items
    const checkoutItems = items.map(item => ({
      sku_id: item.sku_id,
      shop_id: item.shop_id,
      quantity: item.quantity,
      // Additional info for display
      name: item.name,
      price: item.price,
      image: item.image,
      sku_name:item.sku_name
    }));

    setCheckoutItems(checkoutItems);
    router.push(ROUTER.thanhtoan);
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

  if (items.length === 0) {
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
              <CardTitle>Sản phẩm ({items.length})</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {items.map((item) => (
                <div key={item.sku_id}>
                  <div className="flex gap-4">
                    {/* Image */}
                    <div className="relative w-24 h-24 flex-shrink-0">
                      <img
                        src={getImageUrl(item.image)|| "/placeholder.png"}
                        alt={item.name}
                        // fill
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
                        onClick={() => removeFromCart(item.sku_id)}
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
                            updateQuantity(item.sku_id, Math.max(1, item.quantity - 1))
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
                            updateQuantity(item.sku_id, item.quantity + 1)
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
              ))}
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
                  <span>Tạm tính(chưa tính tiền ship)</span>
                  <span>{formatPrice(getTotalPrice())}</span>
                </div>
                {/* <div className="flex justify-between text-sm">
                  <span>Phí vận chuyển</span>
                  <span className="text-black-600">30.000đ</span>
                </div> */}
                <Separator />
                <div className="flex justify-between text-lg font-bold">
                  <span>Tổng cộng</span>
                  <span className="text-[hsl(var(--primary))]">
                    {formatPrice(getTotalPrice() )}
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



// Bạn thực hiện khi người dùng chưa đăng nhập thì vãn cho người dùng thêm sản phẩm vào giỏ hàng và lưu vào localstoge. Khi người dùng đăng nhập vào thì cập nhật tất cả vào API giỏ hàng. Tôi hiện có tất cả các API sau: http://localhost:8000/api/Cart  lấy danh sách tất cả sản phẩm trong card trả về : {
//     "result": {
//         "id": "9acc8610-314d-49cf-973d-31d29f371188",
//         "items": [
//             {
//                 "skuId": "019cbb88-0bfa-4389-a01a-98115af5613f",
//                 "productName": "Default",
//                 "price": 3190000,
//                 "quantity": 1,
//                 "isSelected": true,
//                 "shopId": "019cbb88-0bfa-4389-a01a-98115af5613f",
//                 "addedDate": "2025-11-16T07:16:16.780055Z"
//             }
//         ],
//         "totalItems": 1,
//         "totalPrice": 3190000,
//         "selectedTotalPrice": 3190000
//     },
//     "messages": [],
//     "succeeded": true,
//     "code": 200
// }.  tiếp theo là API đẩy sản phẩm vào trong giỏ hàng : http://localhost:8000/api/Cart/items . với body là : {
//   "SkuId": "019cbb88-0bfa-4389-a01a-98115af5613f",
//   "Quantity": 2
// // }. két quả thành công cũng trả về tương tự như khi lấy giỏ hàng : {
//     "result": {
//         "id": "9acc8610-314d-49cf-973d-31d29f371188",
//         "items": [
//             {
//                 "skuId": "019cbb88-0bfa-4389-a01a-98115af5613f",
//                 "productName": "Default",
//                 "price": 3190000,
//                 "quantity": 3,
//                 "isSelected": true,
//                 "shopId": "019cbb88-0bfa-4389-a01a-98115af5613f",
//                 "addedDate": "2025-11-16T07:18:53.0716205Z"
//             }
//         ],
//         "totalItems": 3,
//         "totalPrice": 9570000,
//         "selectedTotalPrice": 9570000
//     },
//     "messages": [],
//     "succeeded": true,
//     "code": 200
// }
// . bạn dùng thông tin mới này để cập nhật lại giỏ hàng trên giao diện người dùng., tiếp theo là API update số lượng của sku với API : http://localhost:8000/api/Cart/items/019cbb88-0bfa-4389-a01a-98115af5613f . Body là : {
//   "Quantity": 1
// } .  tiếp theo là API Xóa 1 sản phẩm ra khỏi giỏ hàng : http://localhost:8000/api/Cart/items/8a2b639f-3000-4506-ab33-a6268c7830eb . và cuối cùng là API clear toàn bộ giỏ hàng : http://localhost:8000/api/Cart Bạn phát triển lại phần giỏ hàng để tích hợp kèm vào các API này giúp tôi. Dữ liệu giỏ hàng phải luôn là lấy trực tiếp từ API , hạn chế lưu ở local trừ trường hợp người dùng chưa đăng nhập, còn người dùng có giỏ hàng trc khi đăng nhập thì khi đăng nhập phải thêm giỏ hàng vào. có thêm 1 API hiển thị số lượng sản phẩm trong giỏ hàng là : http://localhost:8000/api/Cart/count . với trả về là  ; {
//     "result": 3,
//     "messages": [],
//     "succeeded": true,
//     "code": 200
// }.   