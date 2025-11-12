"use client";

import { useEffect, useState, useMemo } from "react";
import { useRouter } from "@/i18n/routing";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useCheckoutStore, CheckoutItem } from "@/store/checkoutStore";
import { useCartStore } from "@/store/cartStore";
import { useCreateOrder, useGetVouchers } from "@/services/apiService";
import { useToast } from "@/hooks/use-toast";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Textarea } from "@/components/ui/textarea";
import { Separator } from "@/components/ui/separator";
import { Loader2, ShoppingBag, Package, MapPin, Plus, CheckCircle2 } from "lucide-react";
import { Link } from '@/i18n/routing';
import ROUTER from "@/assets/configs/routers";
import VoucherSelector from "@/components/VoucherSelector";
import AddressDialog from "@/components/AddressDialog";
import { 
  Voucher, 
  AppliedVoucher, 
  CategorizedVouchers,
  ShopGroup 
} from "@/types/voucher.types";
import { UserAddress } from "@/types/address.types";
import { getCookieValues } from "@/assets/helpers/cookies";
import { ACCESS_TOKEN, INFO_USER } from "@/assets/configs/request";
import API from "@/assets/configs/api";
import { apiClient } from "@/lib/apiClient";

// Payment Method UUIDs (từ backend)
const PAYMENT_METHODS = {
  COD: "b2c3d4e5-f6a7-8901-2345-67890abcdef1", // Cash on Delivery
  MOMO: "a1b2c3d4-e5f6-7890-1234-567890abcdef", // MoMo
};

// Default shipping fee per shop
const DEFAULT_SHIPPING_FEE = 30000;

interface UserProfile {
  id: string;
  userId: string;
  name: string;
  firstName: string;
  lastName: string;
  dob: string;
  phone_number: string;
  gender: string;
}

export default function CheckoutPage() {
  const router = useRouter();
  const { toast } = useToast();

  // Local state for checkout items (avoid Zustand hydration issues)
  const [checkoutItems, setCheckoutItems] = useState<CheckoutItem[]>([]);
  const [selectedPayment, setSelectedPayment] = useState<string>(PAYMENT_METHODS.COD);
  const [orderNote, setOrderNote] = useState<string>("");
  const [isMounted, setIsMounted] = useState<boolean>(false);
  const [isCheckingAuth, setIsCheckingAuth] = useState<boolean>(true);
  
  // User profile and addresses
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [addresses, setAddresses] = useState<UserAddress[]>([]);
  const [selectedAddress, setSelectedAddress] = useState<UserAddress | null>(null);
  const [isLoadingAddresses, setIsLoadingAddresses] = useState(false);
  const [addressDialogOpen, setAddressDialogOpen] = useState(false);

  // Voucher states
  const [appliedPlatformOrderVoucher, setAppliedPlatformOrderVoucher] = useState<AppliedVoucher | undefined>();
  const [appliedPlatformShippingVoucher, setAppliedPlatformShippingVoucher] = useState<AppliedVoucher | undefined>();
  const [appliedShopVouchers, setAppliedShopVouchers] = useState<Map<string, AppliedVoucher>>(new Map());

  // Zustand store actions only (not state)
  const clearCheckout = useCheckoutStore((state) => state.clearCheckout);
  const clearCart = useCartStore((state) => state.clearCart);

  // Fetch vouchers
  const { data: vouchersData, isLoading: vouchersLoading } = useGetVouchers();

  // React Query Mutation
  const { mutate: createOrder, isPending } = useCreateOrder();

  // Check authentication first
  useEffect(() => {
    const token = getCookieValues<string>(ACCESS_TOKEN);
    
    if (!token) {
      toast({
        title: "Vui lòng đăng nhập",
        description: "Bạn cần đăng nhập để thực hiện thanh toán",
        variant: "destructive",
      });
      
      const timer = setTimeout(() => {
        router.push(ROUTER.auth.login);
      }, 1500);
      
      return () => clearTimeout(timer);
    }
    
    setIsCheckingAuth(false);
  }, [router, toast]);

  // Load user profile and addresses
  useEffect(() => {
    if (isCheckingAuth) return;
    
    const loadUserData = async () => {
      try {
        // Load from localStorage first
        const userInfo = localStorage.getItem(INFO_USER);
        if (userInfo) {
          const userData = JSON.parse(userInfo);
          setProfile(userData);
        }
        
        // Load addresses from API
        await loadAddresses();
      } catch (error) {
        console.error("Error loading user data:", error);
      }
    };
    
    loadUserData();
  }, [isCheckingAuth]);

  // Load addresses function
  const loadAddresses = async () => {
    setIsLoadingAddresses(true);
    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      if (!token) return;

      const response = await apiClient.get(API.user.addresses);
      
      if (response.data.code === 10000) {
        const addressesData = response.data.result || [];
        setAddresses(addressesData);
        
        // Auto select first address if available
        if (addressesData.length > 0 && !selectedAddress) {
          setSelectedAddress(addressesData[0]);
        }
      }
    } catch (error) {
      console.error("Error loading addresses:", error);
      toast({
        title: "Lỗi",
        description: "Không thể tải danh sách địa chỉ",
        variant: "destructive",
      });
    } finally {
      setIsLoadingAddresses(false);
    }
  };

  // Load checkout items from Zustand ONLY on client side
  useEffect(() => {
    if (isCheckingAuth) return; // Wait for auth check
    
    setIsMounted(true);
    
    // Read from Zustand store after mount
    const storeItems = useCheckoutStore.getState().items;
    setCheckoutItems(storeItems);

    // Redirect if empty
    if (storeItems.length === 0) {
      toast({
        title: "Giỏ hàng trống",
        description: "Vui lòng thêm sản phẩm vào giỏ hàng",
        variant: "destructive",
      });
      
      const timer = setTimeout(() => {
        router.push(ROUTER.giohang);
      }, 500);
      
      return () => clearTimeout(timer);
    }
  }, [router, toast, isCheckingAuth]);

  // Helper functions
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

  // Group items by shop_id
  const shopGroups = useMemo((): ShopGroup[] => {
    const grouped = new Map<string, CheckoutItem[]>();
    
    checkoutItems.forEach((item) => {
      const shopId = item.shop_id || 'unknown';
      if (!grouped.has(shopId)) {
        grouped.set(shopId, []);
      }
      grouped.get(shopId)!.push(item);
    });

    return Array.from(grouped.entries()).map(([shop_id, items]) => {
      const subtotal = items.reduce((sum, item) => sum + (item.price || 0) * item.quantity, 0);
      const appliedShopVoucher = appliedShopVouchers.get(shop_id);
      
      return {
        shop_id,
        shop_name: `Shop ${shop_id.slice(0, 8)}`, // Default name
        items,
        subtotal,
        shipping_fee: DEFAULT_SHIPPING_FEE,
        shop_voucher: appliedShopVoucher,
        total: subtotal + DEFAULT_SHIPPING_FEE,
      };
    });
  }, [checkoutItems, appliedShopVouchers]);

  // Categorize vouchers
  const categorizedVouchers = useMemo((): CategorizedVouchers => {
    if (!vouchersData?.data) {
      return {
        platformOrderVouchers: [],
        platformShippingVouchers: [],
        assignedVouchers: [],
        shopVouchers: new Map(),
      };
    }

    const platformOrderVouchers: Voucher[] = [];
    const platformShippingVouchers: Voucher[] = [];
    const assignedVouchers: Voucher[] = [];
    const shopVouchersMap = new Map<string, Voucher[]>();

    vouchersData.data.forEach((voucher) => {
      // ASSIGNED vouchers - riêng cho user
      if (voucher.audience_type === "ASSIGNED") {
        assignedVouchers.push(voucher);
      }
      // PUBLIC PLATFORM vouchers
      else if (voucher.owner_type === "PLATFORM" && voucher.audience_type === "PUBLIC") {
        if (voucher.applies_to_type === "SHIPPING_FEE") {
          platformShippingVouchers.push(voucher);
        } else {
          platformOrderVouchers.push(voucher);
        }
      }
      // SHOP vouchers
      else if (voucher.owner_type === "SHOP") {
        const shopId = voucher.owner_id;
        if (!shopVouchersMap.has(shopId)) {
          shopVouchersMap.set(shopId, []);
        }
        shopVouchersMap.get(shopId)!.push(voucher);
      }
    });

    return {
      platformOrderVouchers,
      platformShippingVouchers,
      assignedVouchers,
      shopVouchers: shopVouchersMap,
    };
  }, [vouchersData]);

  // Calculate discount
  const calculateVoucherDiscount = (voucher: Voucher, amount: number): number => {
    const minAmount = parseFloat(voucher.min_purchase_amount);
    if (amount < minAmount) return 0;

    const discountValue = parseFloat(voucher.discount_value);
    
    if (voucher.discount_type === "FIXED_AMOUNT") {
      return discountValue;
    } else {
      // PERCENTAGE
      let discount = (amount * discountValue) / 100;
      if (voucher.max_discount_amount) {
        const maxDiscount = parseFloat(voucher.max_discount_amount);
        discount = Math.min(discount, maxDiscount);
      }
      return discount;
    }
  };

  // Calculate totals
  const orderSubtotal = useMemo(() => {
    return shopGroups.reduce((sum, group) => sum + group.subtotal, 0);
  }, [shopGroups]);

  const totalShippingFee = useMemo(() => {
    return shopGroups.length * DEFAULT_SHIPPING_FEE;
  }, [shopGroups]);

  // Platform order voucher discount
  const platformOrderDiscount = useMemo(() => {
    if (!appliedPlatformOrderVoucher) return 0;
    const voucher = vouchersData?.data.find(v => v.id === appliedPlatformOrderVoucher.voucher_id);
    if (!voucher) return 0;
    return calculateVoucherDiscount(voucher, orderSubtotal);
  }, [appliedPlatformOrderVoucher, vouchersData, orderSubtotal]);

  // Platform shipping voucher discount
  const platformShippingDiscount = useMemo(() => {
    if (!appliedPlatformShippingVoucher) return 0;
    const voucher = vouchersData?.data.find(v => v.id === appliedPlatformShippingVoucher.voucher_id);
    if (!voucher) return 0;
    return Math.min(calculateVoucherDiscount(voucher, orderSubtotal), totalShippingFee);
  }, [appliedPlatformShippingVoucher, vouchersData, orderSubtotal, totalShippingFee]);

  // Shop vouchers discount
  const shopVouchersDiscount = useMemo(() => {
    let total = 0;
    appliedShopVouchers.forEach((appliedVoucher, shop_id) => {
      const voucher = vouchersData?.data.find(v => v.id === appliedVoucher.voucher_id);
      const shopGroup = shopGroups.find(g => g.shop_id === shop_id);
      if (voucher && shopGroup) {
        total += calculateVoucherDiscount(voucher, shopGroup.subtotal);
      }
    });
    return total;
  }, [appliedShopVouchers, vouchersData, shopGroups]);

  const totalDiscount = platformOrderDiscount + platformShippingDiscount + shopVouchersDiscount;
  const grandTotal = orderSubtotal + totalShippingFee - totalDiscount;

  // Voucher handlers
  const handleApplyPlatformOrderVoucher = (voucher: Voucher | null) => {
    if (!voucher) {
      setAppliedPlatformOrderVoucher(undefined);
      return;
    }

    const discount = calculateVoucherDiscount(voucher, orderSubtotal);
    setAppliedPlatformOrderVoucher({
      voucher_id: voucher.id,
      voucher_code: voucher.voucher_code,
      discount_amount: discount,
      applies_to: voucher.applies_to_type,
    });
  };

  const handleApplyPlatformShippingVoucher = (voucher: Voucher | null) => {
    if (!voucher) {
      setAppliedPlatformShippingVoucher(undefined);
      return;
    }

    const discount = Math.min(calculateVoucherDiscount(voucher, orderSubtotal), totalShippingFee);
    setAppliedPlatformShippingVoucher({
      voucher_id: voucher.id,
      voucher_code: voucher.voucher_code,
      discount_amount: discount,
      applies_to: voucher.applies_to_type,
    });
  };

  const handleApplyShopVoucher = (shop_id: string, voucher: Voucher | null) => {
    const newMap = new Map(appliedShopVouchers);
    
    if (!voucher) {
      newMap.delete(shop_id);
    } else {
      const shopGroup = shopGroups.find(g => g.shop_id === shop_id);
      if (shopGroup) {
        const discount = calculateVoucherDiscount(voucher, shopGroup.subtotal);
        newMap.set(shop_id, {
          shop_id,
          voucher_id: voucher.id,
          voucher_code: voucher.voucher_code,
          discount_amount: discount,
          applies_to: voucher.applies_to_type,
        });
      }
    }
    
    setAppliedShopVouchers(newMap);
  };

  const calculateTotal = () => {
    return grandTotal;
  };

  const handleAddAddress = () => {
    setAddressDialogOpen(true);
  };

  const handleSelectAddress = (address: UserAddress) => {
    setSelectedAddress(address);
  };

  // Handle form submit
  const handlePlaceOrder = () => {
    // Validate selected address
    if (!selectedAddress) {
      toast({
        title: "Chưa chọn địa chỉ",
        description: "Vui lòng chọn địa chỉ giao hàng",
        variant: "destructive",
      });
      return;
    }
    // Prepare voucher_shop array (shop vouchers)
    const voucher_shop: Array<{ voucher_id: string; shop_id: string }> = [];
    appliedShopVouchers.forEach((voucher, shop_id) => {
      voucher_shop.push({
        voucher_id: voucher.voucher_id,
        shop_id: shop_id,
      });
    });

    // Prepare shipping address from selected address
    const shippingAddressData = {
      fullName: selectedAddress.name,
      phone: selectedAddress.phoneNumber,
      address: selectedAddress.address.other,
      ward: selectedAddress.address.ward.fullName,
      district: selectedAddress.address.ward.district.fullName,
      city: selectedAddress.address.ward.district.province.fullName,
      wardId: selectedAddress.address.ward.id,
      districtId: selectedAddress.address.ward.district.id,
      provinceId: selectedAddress.address.ward.district.province.id,
    };

    // Prepare payload for API
    const orderPayload: any = {
      shippingAddress: shippingAddressData,
      paymentMethod: selectedPayment,
      items: checkoutItems.map(item => ({
        sku_id: item.sku_id,
        shop_id: item.shop_id,
        quantity: item.quantity,
      })),
      note: orderNote || "",
    };

    // Add vouchers if applied
    if (voucher_shop.length > 0) {
      orderPayload.voucher_shop = voucher_shop;
    }

    if (appliedPlatformOrderVoucher) {
      orderPayload.voucher_site_id = appliedPlatformOrderVoucher.voucher_id;
    }

    if (appliedPlatformShippingVoucher) {
      orderPayload.voucher_shipping_id = appliedPlatformShippingVoucher.voucher_id;
    }
    // orderPayload.email = "vinhhien12z@gmail.com"
    console.log("Order Payload:", orderPayload);

    // Submit order
    createOrder(orderPayload, {
      onSuccess: (response) => {
        toast({
          title: "Đặt hàng thành công!",
          description: `Mã đơn hàng: ${response.result?.orderCode  || 'N/A'} - Tổng tiền: ${formatPrice(response.result?.grandTotal || 0)}`,
        });

        // Clear cart and checkout
        clearCart();
        clearCheckout();
        
        // Redirect based on payment method
        if (selectedPayment === PAYMENT_METHODS.MOMO && response.result?.paymentUrl) {
          window.location.href = response.result.paymentUrl;
        } else {
          router.push(`${ROUTER.dat_hang_thanh_cong}?order_id=${response.result?.orderId || ''}`);
        }
      },
      onError: (error: any) => {
        toast({
          title: "Đặt hàng thất bại",
          description: error.response?.data?.message || "Có lỗi xảy ra, vui lòng thử lại",
          variant: "destructive",
        });
      },
    });
  };

  // Show loading while checking auth or mounting
  if (isCheckingAuth || !isMounted || checkoutItems.length === 0) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="flex flex-col items-center justify-center gap-4">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">
            {isCheckingAuth ? "Đang kiểm tra đăng nhập..." : "Đang tải..."}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Thanh toán</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left: Shipping & Payment */}
        <div className="lg:col-span-2 space-y-6">
          {/* Shipping Address Selection */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <MapPin className="h-5 w-5" />
                Địa chỉ giao hàng
              </CardTitle>
              <Button 
                variant="outline" 
                size="sm" 
                onClick={handleAddAddress}
                type="button"
              >
                <Plus className="h-4 w-4 mr-2" />
                Thêm địa chỉ mới
              </Button>
            </CardHeader>
            <CardContent>
              {isLoadingAddresses ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-primary" />
                </div>
              ) : addresses.length > 0 ? (
                <div className="space-y-3">
                  {addresses.map((address) => (
                    <div
                      key={address.id}
                      onClick={() => handleSelectAddress(address)}
                      className={`
                        p-4 border-2 rounded-lg cursor-pointer transition-all
                        ${selectedAddress?.id === address.id 
                          ? 'border-primary bg-primary/5' 
                          : 'border-gray-200 hover:border-primary/50'
                        }
                      `}
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            <h4 className="font-semibold">{address.name}</h4>
                            {selectedAddress?.id === address.id && (
                              <CheckCircle2 className="h-5 w-5 text-primary" />
                            )}
                          </div>
                          <p className="text-sm text-muted-foreground mb-1">
                            📱 {address.phoneNumber}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            📍 {address.address.other}, {address.address.ward.fullName}, {address.address.ward.district.fullName}, {address.address.ward.district.province.fullName}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <MapPin className="h-12 w-12 mx-auto mb-3 opacity-20" />
                  <p className="mb-4">Chưa có địa chỉ giao hàng</p>
                  <Button variant="outline" onClick={handleAddAddress} type="button">
                    <Plus className="h-4 w-4 mr-2" />
                    Thêm địa chỉ đầu tiên
                  </Button>
                </div>
              )}
              
              {!selectedAddress && addresses.length > 0 && (
                <p className="text-sm text-red-500 mt-3">
                  ⚠️ Vui lòng chọn địa chỉ giao hàng
                </p>
              )}
            </CardContent>
          </Card>

          {/* Order Note */}
          <Card>
            <CardHeader>
              <CardTitle>Ghi chú đơn hàng</CardTitle>
            </CardHeader>
            <CardContent>
              <Textarea
                value={orderNote}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setOrderNote(e.target.value)}
                placeholder="Giao hàng sau 5h chiều..."
                rows={3}
              />
            </CardContent>
          </Card>

            {/* Payment Method */}
            <Card>
              <CardHeader>
                <CardTitle>Phương thức thanh toán</CardTitle>
              </CardHeader>
              <CardContent>
                <RadioGroup value={selectedPayment} onValueChange={setSelectedPayment}>
                  <div className="flex items-center space-x-2 p-4 border rounded-lg hover:bg-gray-50">
                    <RadioGroupItem value={PAYMENT_METHODS.COD} id="cod" />
                    <Label htmlFor="cod" className="flex-1 cursor-pointer">
                      <div className="font-medium">Thanh toán khi nhận hàng (COD)</div>
                      <div className="text-sm text-gray-500">Thanh toán bằng tiền mặt khi nhận hàng</div>
                    </Label>
                  </div>

                  <div className="flex items-center space-x-2 p-4 border rounded-lg hover:bg-gray-50">
                    <RadioGroupItem value={PAYMENT_METHODS.MOMO} id="momo" />
                    <Label htmlFor="momo" className="flex-1 cursor-pointer">
                      <div className="font-medium">Ví MoMo</div>
                      <div className="text-sm text-gray-500">Thanh toán qua ví điện tử MoMo</div>
                    </Label>
                  </div>
                </RadioGroup>
              </CardContent>
            </Card>
        </div>

        {/* Right: Order Summary */}
        <div className="lg:col-span-1">
            <Card className="sticky top-4">
              <CardHeader>
                <CardTitle>Đơn hàng ({checkoutItems.length} sản phẩm)</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Shop Groups */}
                <div className="space-y-4 max-h-96 overflow-y-auto">
                  {shopGroups.map((shopGroup) => (
                    <div key={shopGroup.shop_id} className="border rounded-lg p-3 space-y-2">
                      <div className="flex items-center gap-2 mb-2">
                        <Package className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm font-semibold">
                          {shopGroup.shop_name || `Shop ${shopGroup.shop_id.slice(0, 8)}`}
                        </span>
                      </div>
                      
                      {/* Items in this shop */}
                      {shopGroup.items.map((item) => (
                        <div key={item.sku_id} className="flex gap-2">
                          <div className="relative w-12 h-12 flex-shrink-0">
                            <img
                              src={getImageUrl(item.image)}
                              alt={item.name || 'Product'}
                              className="w-full h-full object-cover rounded"
                            />
                          </div>
                          <div className="flex-1 min-w-0">
                            <p className="text-xs line-clamp-1">{item.name}</p>
                            <p className="text-xs text-gray-500">x{item.quantity}</p>
                            <p className="text-xs font-semibold text-primary">
                              {formatPrice((item.price || 0) * item.quantity)}
                            </p>
                          </div>
                        </div>
                      ))}
                      
                      <Separator className="my-2" />
                      
                      {/* Shop totals */}
                      <div className="space-y-1 text-xs">
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Tạm tính:</span>
                          <span>{formatPrice(shopGroup.subtotal)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Phí ship:</span>
                          <span>{formatPrice(shopGroup.shipping_fee)}</span>
                        </div>
                        {shopGroup.shop_voucher && (
                          <div className="flex justify-between text-green-600">
                            <span>Voucher shop:</span>
                            <span>-{formatPrice(shopGroup.shop_voucher.discount_amount)}</span>
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>

                <Separator />

                {/* Voucher Selector */}
                <VoucherSelector
                  platformOrderVouchers={categorizedVouchers.platformOrderVouchers}
                  platformShippingVouchers={categorizedVouchers.platformShippingVouchers}
                  assignedVouchers={categorizedVouchers.assignedVouchers}
                  shopVouchers={categorizedVouchers.shopVouchers}
                  shopGroups={shopGroups.map(g => ({
                    shop_id: g.shop_id,
                    shop_name: g.shop_name,
                    subtotal: g.subtotal,
                  }))}
                  totalOrderAmount={orderSubtotal}
                  appliedPlatformOrderVoucher={appliedPlatformOrderVoucher}
                  appliedPlatformShippingVoucher={appliedPlatformShippingVoucher}
                  appliedShopVouchers={appliedShopVouchers}
                  onApplyPlatformOrderVoucher={handleApplyPlatformOrderVoucher}
                  onApplyPlatformShippingVoucher={handleApplyPlatformShippingVoucher}
                  onApplyShopVoucher={handleApplyShopVoucher}
                />

                <Separator />

                {/* Price Summary */}
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>Tạm tính ({shopGroups.length} shop)</span>
                    <span>{formatPrice(orderSubtotal)}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span>Phí vận chuyển</span>
                    <span>{formatPrice(totalShippingFee)}</span>
                  </div>
                  
                  {/* Discounts */}
                  {platformOrderDiscount > 0 && (
                    <div className="flex justify-between text-sm text-green-600">
                      <span>Voucher sàn</span>
                      <span>-{formatPrice(platformOrderDiscount)}</span>
                    </div>
                  )}
                  {platformShippingDiscount > 0 && (
                    <div className="flex justify-between text-sm text-green-600">
                      <span>Voucher ship</span>
                      <span>-{formatPrice(platformShippingDiscount)}</span>
                    </div>
                  )}
                  {shopVouchersDiscount > 0 && (
                    <div className="flex justify-between text-sm text-green-600">
                      <span>Voucher shop</span>
                      <span>-{formatPrice(shopVouchersDiscount)}</span>
                    </div>
                  )}
                  
                  {totalDiscount > 0 && <Separator />}
                  
                  <div className="flex justify-between text-lg font-bold">
                    <span>Tổng cộng</span>
                    <span className="text-primary">{formatPrice(grandTotal)}</span>
                  </div>
                  
                  {totalDiscount > 0 && (
                    <p className="text-xs text-green-600 text-right">
                      Bạn đã tiết kiệm {formatPrice(totalDiscount)}
                    </p>
                  )}
                </div>

                {/* Submit Button */}
                <Button
                  onClick={handlePlaceOrder}
                  className="w-full"
                  size="lg"
                  disabled={isPending || vouchersLoading || !selectedAddress}
                >
                  {isPending ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Đang xử lý...
                    </>
                  ) : (
                    "Đặt hàng"
                  )}
                </Button>

                <Link href={ROUTER.giohang}>
                  <Button variant="outline" className="w-full" size="lg">
                    Quay lại giỏ hàng
                  </Button>
                </Link>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Address Dialog */}
        <AddressDialog
          open={addressDialogOpen}
          onOpenChange={setAddressDialogOpen}
          onSuccess={loadAddresses}
        />
      </div>
    );
  }
