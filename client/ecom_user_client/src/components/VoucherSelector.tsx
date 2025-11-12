"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Ticket, Check, X, Info } from "lucide-react";
import { Voucher, AppliedVoucher } from "@/types/voucher.types";
import { cn } from "@/lib/utils";

interface VoucherSelectorProps {
  // Vouchers categorized
  platformOrderVouchers: Voucher[];
  platformShippingVouchers: Voucher[];
  assignedVouchers: Voucher[];
  shopVouchers: Map<string, Voucher[]>;
  
  // Current order info
  shopGroups: Array<{
    shop_id: string;
    shop_name?: string;
    subtotal: number;
  }>;
  totalOrderAmount: number;
  
  // Applied vouchers
  appliedPlatformOrderVoucher?: AppliedVoucher;
  appliedPlatformShippingVoucher?: AppliedVoucher;
  appliedShopVouchers: Map<string, AppliedVoucher>;
  
  // Callbacks
  onApplyPlatformOrderVoucher: (voucher: Voucher | null) => void;
  onApplyPlatformShippingVoucher: (voucher: Voucher | null) => void;
  onApplyShopVoucher: (shop_id: string, voucher: Voucher | null) => void;
}

export default function VoucherSelector({
  platformOrderVouchers,
  platformShippingVouchers,
  assignedVouchers,
  shopVouchers,
  shopGroups,
  totalOrderAmount,
  appliedPlatformOrderVoucher,
  appliedPlatformShippingVoucher,
  appliedShopVouchers,
  onApplyPlatformOrderVoucher,
  onApplyPlatformShippingVoucher,
  onApplyShopVoucher,
}: VoucherSelectorProps) {
  // Merge assigned vouchers into their respective categories
  const assignedShippingVouchers = assignedVouchers.filter(
    v => v.applies_to_type === "SHIPPING_FEE"
  );
  const assignedPlatformOrderVouchers = assignedVouchers.filter(
    v => v.applies_to_type === "ORDER_TOTAL" && v.owner_type === "PLATFORM"
  );
  
  // Merge vouchers for display
  const allPlatformOrderVouchers = [...platformOrderVouchers, ...assignedPlatformOrderVouchers];
  const allPlatformShippingVouchers = [...platformShippingVouchers, ...assignedShippingVouchers];
  
  // Merge assigned shop vouchers into shop vouchers map
  const allShopVouchers = new Map(shopVouchers);
  assignedVouchers.filter(
    v => v.applies_to_type === "ORDER_TOTAL" && v.owner_type === "SHOP" && v.owner_id
  ).forEach(voucher => {
    const shopId = voucher.owner_id!;
    const existing = allShopVouchers.get(shopId) || [];
    allShopVouchers.set(shopId, [...existing, voucher]);
  });
  const [open, setOpen] = useState(false);

  const formatPrice = (price: number | string) => {
    const numPrice = typeof price === "string" ? parseFloat(price) : price;
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(numPrice);
  };

  const isVoucherEligible = (voucher: Voucher, amount: number): boolean => {
    const minAmount = parseFloat(voucher.min_purchase_amount);
    return amount >= minAmount;
  };

  const calculateDiscount = (voucher: Voucher, amount: number): number => {
    if (!isVoucherEligible(voucher, amount)) return 0;

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

  // Sort vouchers by discount value (highest first)
  const sortVouchersByValue = (vouchers: Voucher[], applicableAmount: number): Voucher[] => {
    return [...vouchers].sort((a, b) => {
      const discountA = calculateDiscount(a, applicableAmount);
      const discountB = calculateDiscount(b, applicableAmount);
      return discountB - discountA; // Descending order
    });
  };

  const renderVoucherCard = (
    voucher: Voucher,
    isApplied: boolean,
    canApply: boolean,
    onToggle: () => void,
    applicableAmount?: number
  ) => {
    const discount = applicableAmount ? calculateDiscount(voucher, applicableAmount) : 0;
    const eligible = applicableAmount ? isVoucherEligible(voucher, applicableAmount) : false;

    return (
      <Card
        key={voucher.id}
        className={cn(
          "relative overflow-hidden transition-all",
          isApplied && "border-antique-gold border-2",
          !eligible && "opacity-50"
        )}
      >
        <CardContent className="p-4">
          <div className="flex items-start justify-between gap-3">
            <div className="flex-1">
              {/* Voucher Icon & Code */}
              <div className="flex items-center gap-2 mb-2">
                <Ticket className="h-5 w-5 text-antique-gold" />
                <Badge variant={voucher.audience_type === "ASSIGNED" ? "default" : "secondary"}>
                  {voucher.voucher_code}
                </Badge>
                {voucher.audience_type === "ASSIGNED" && (
                  <Badge variant="destructive" className="text-xs">
                    Riêng
                  </Badge>
                )}
              </div>

              {/* Voucher Name */}
              <h4 className="font-semibold text-base mb-1">{voucher.name}</h4>

              {/* Discount Info */}
              <p className="text-sm text-muted-foreground mb-2">
                {voucher.discount_type === "FIXED_AMOUNT"
                  ? `Giảm ${formatPrice(voucher.discount_value)}`
                  : `Giảm ${voucher.discount_value}%${
                      voucher.max_discount_amount
                        ? ` (Tối đa ${formatPrice(voucher.max_discount_amount)})`
                        : ""
                    }`}
              </p>

              {/* Min Purchase */}
              <p className="text-xs text-muted-foreground">
                Đơn tối thiểu: {formatPrice(voucher.min_purchase_amount)}
              </p>

              {/* Discount Amount */}
              {eligible && discount > 0 && (
                <p className="text-sm font-semibold text-green-600 mt-2">
                  Bạn sẽ giảm: {formatPrice(discount)}
                </p>
              )}

              {/* Not eligible warning */}
              {!eligible && applicableAmount !== undefined && (
                <p className="text-xs text-red-500 mt-2 flex items-center gap-1">
                  <Info className="h-3 w-3" />
                  Chưa đủ điều kiện áp dụng
                </p>
              )}

              {/* Usage info */}
              <p className="text-xs text-muted-foreground mt-1">
                Còn {voucher.total_quantity - voucher.used_quantity}/{voucher.total_quantity} lượt
              </p>
            </div>

            {/* Apply Button */}
            <div className="flex flex-col items-end gap-2">
              {canApply && eligible ? (
                <Button
                  size="sm"
                  variant={isApplied ? "outline" : "default"}
                  onClick={onToggle}
                  className="min-w-[80px]"
                >
                  {isApplied ? (
                    <>
                      <Check className="h-4 w-4 mr-1" />
                      Đã dùng
                    </>
                  ) : (
                    "Áp dụng"
                  )}
                </Button>
              ) : (
                <Button size="sm" variant="ghost" disabled className="min-w-[80px]">
                  Không thể dùng
                </Button>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    );
  };

  const getTotalAppliedVouchers = () => {
    let count = 0;
    if (appliedPlatformOrderVoucher) count++;
    if (appliedPlatformShippingVoucher) count++;
    count += appliedShopVouchers.size;
    return count;
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" className="w-full justify-between">
          <span className="flex items-center gap-2">
            <Ticket className="h-4 w-4" />
            Chọn Voucher
          </span>
          {getTotalAppliedVouchers() > 0 && (
            <Badge variant="default">{getTotalAppliedVouchers()} đã chọn</Badge>
          )}
        </Button>
      </DialogTrigger>

      <DialogContent className="max-w-3xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Chọn Voucher</DialogTitle>
          <DialogDescription>
            Chọn voucher phù hợp để tiết kiệm tối đa cho đơn hàng của bạn
          </DialogDescription>
        </DialogHeader>

        <Tabs defaultValue="platform-order" className="w-full">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="platform-order">
              Sàn ({allPlatformOrderVouchers.length})
            </TabsTrigger>
            <TabsTrigger value="platform-shipping">
              Freeship ({allPlatformShippingVouchers.length})
            </TabsTrigger>
            <TabsTrigger value="shop">
              Shop ({allShopVouchers.size})
            </TabsTrigger>
          </TabsList>

          {/* Platform Order Vouchers */}
          <TabsContent value="platform-order" className="space-y-3 mt-4">
            <p className="text-sm text-muted-foreground">
              Áp dụng cho tổng đơn hàng: {formatPrice(totalOrderAmount)}
            </p>
            <Separator />
            {allPlatformOrderVouchers.length === 0 ? (
              <p className="text-center text-muted-foreground py-8">Không có voucher</p>
            ) : (
              sortVouchersByValue(allPlatformOrderVouchers, totalOrderAmount).map((voucher) =>
                renderVoucherCard(
                  voucher,
                  appliedPlatformOrderVoucher?.voucher_id === voucher.id,
                  true,
                  () => {
                    if (appliedPlatformOrderVoucher?.voucher_id === voucher.id) {
                      onApplyPlatformOrderVoucher(null);
                    } else {
                      onApplyPlatformOrderVoucher(voucher);
                    }
                  },
                  totalOrderAmount
                )
              )
            )}
          </TabsContent>

          {/* Platform Shipping Vouchers */}
          <TabsContent value="platform-shipping" className="space-y-3 mt-4">
            <p className="text-sm text-muted-foreground">
              Áp dụng cho phí vận chuyển tổng: {formatPrice(shopGroups.length * 30000)}
            </p>
            <Separator />
            {allPlatformShippingVouchers.length === 0 ? (
              <p className="text-center text-muted-foreground py-8">Không có voucher</p>
            ) : (
              sortVouchersByValue(allPlatformShippingVouchers, shopGroups.length * 30000).map((voucher) =>
                renderVoucherCard(
                  voucher,
                  appliedPlatformShippingVoucher?.voucher_id === voucher.id,
                  true,
                  () => {
                    if (appliedPlatformShippingVoucher?.voucher_id === voucher.id) {
                      onApplyPlatformShippingVoucher(null);
                    } else {
                      onApplyPlatformShippingVoucher(voucher);
                    }
                  },
                  shopGroups.length * 30000 // Tổng phí ship
                )
              )
            )}
          </TabsContent>

          {/* Shop Vouchers */}
          <TabsContent value="shop" className="space-y-4 mt-4">
            {shopGroups.length === 0 ? (
              <p className="text-center text-muted-foreground py-8">Không có shop nào</p>
            ) : (
              shopGroups.map((shopGroup) => {
                const shopVoucherList = allShopVouchers.get(shopGroup.shop_id) || [];
                return (
                  <div key={shopGroup.shop_id} className="space-y-3">
                    <div className="flex items-center justify-between">
                      <h4 className="font-semibold">
                        {shopGroup.shop_name || `Shop ${shopGroup.shop_id.slice(0, 8)}`}
                      </h4>
                      <p className="text-sm text-muted-foreground">
                        Tổng: {formatPrice(shopGroup.subtotal)}
                      </p>
                    </div>
                    <Separator />
                    {shopVoucherList.length === 0 ? (
                      <p className="text-sm text-muted-foreground">Không có voucher</p>
                    ) : (
                      sortVouchersByValue(shopVoucherList, shopGroup.subtotal).map((voucher) =>
                        renderVoucherCard(
                          voucher,
                          appliedShopVouchers.get(shopGroup.shop_id)?.voucher_id === voucher.id,
                          true,
                          () => {
                            if (
                              appliedShopVouchers.get(shopGroup.shop_id)?.voucher_id === voucher.id
                            ) {
                              onApplyShopVoucher(shopGroup.shop_id, null);
                            } else {
                              onApplyShopVoucher(shopGroup.shop_id, voucher);
                            }
                          },
                          shopGroup.subtotal
                        )
                      )
                    )}
                  </div>
                );
              })
            )}
          </TabsContent>
        </Tabs>

        <div className="flex justify-end gap-3 mt-6">
          <Button variant="outline" onClick={() => setOpen(false)}>
            Đóng
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
