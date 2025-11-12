// Voucher Types
export type VoucherAppliesTo = "ORDER_TOTAL" | "SHIPPING_FEE";
export type VoucherAudienceType = "PUBLIC" | "ASSIGNED";
export type VoucherOwnerType = "PLATFORM" | "SHOP";
export type VoucherDiscountType = "FIXED_AMOUNT" | "PERCENTAGE";

export interface Voucher {
  id: string;
  voucher_code: string;
  name: string;
  discount_type: VoucherDiscountType;
  discount_value: string; // "30000.00"
  applies_to_type: VoucherAppliesTo;
  audience_type: VoucherAudienceType;
  owner_type: VoucherOwnerType;
  owner_id: string; // PLATFORM UUID or SHOP UUID
  min_purchase_amount: string; // "150000.00"
  max_discount_amount: string | null;
  max_usage_per_user: number;
  total_quantity: number;
  used_quantity: number;
  is_active: boolean;
  start_date: Record<string, any>;
  end_date: Record<string, any>;
  created_at: Record<string, any>;
  updated_at: Record<string, any>;
}

export interface VoucherApiResponse {
  code: number;
  message: string;
  status: string;
  result: {
    data: Voucher[];
  };
}

// Voucher by category
export interface CategorizedVouchers {
  platformOrderVouchers: Voucher[]; // Voucher sàn giảm giá đơn hàng (PUBLIC)
  platformShippingVouchers: Voucher[]; // Voucher sàn freeship (PUBLIC)
  assignedVouchers: Voucher[]; // Voucher riêng người dùng (ASSIGNED)
  shopVouchers: Map<string, Voucher[]>; // Voucher theo từng shop (key: shop_id)
}

// Applied voucher
export interface AppliedVoucher {
  shop_id?: string; // Nếu là voucher shop
  voucher_id: string;
  voucher_code: string;
  discount_amount: number;
  applies_to: VoucherAppliesTo;
}

// Shop group (để tính tổng theo shop)
export interface ShopGroup {
  shop_id: string;
  shop_name?: string;
  items: any[]; // CheckoutItem[]
  subtotal: number;
  shipping_fee: number;
  shop_voucher?: AppliedVoucher;
  total: number;
}
