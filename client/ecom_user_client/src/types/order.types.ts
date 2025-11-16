// Dùng cho Checkout
export interface OrderItemPayload {
  sku_id: string;
  shop_id: string;
  quantity: number;
}

export interface ShippingAddress {
  fullName: string;
  phone: string;
  address: string;
  district: string;
  city: string;
  postalCode: string;
}

export interface CreateOrderPayload {
  shippingAddress: ShippingAddress;
  paymentMethod: string; // UUID của payment method (COD, MoMo, etc.)
  items: OrderItemPayload[];
  vouchers: string[]; // Array of voucher codes
  note: string; // Ghi chú đơn hàng
}

export interface CreateOrderSuccessResponse {
  code: number;
  message: string;
  status: string;
  result?: {
    orderCode: string;
    orderId: string;
    paymentUrl?: string; // For MoMo, ZaloPay, etc.
    grandTotal: number;
    shopOrders: string[];
  };
}

// Dùng cho Order List
export type OrderStatus = 
  | 'AWAITING_PAYMENT' 
  | 'PROCESSING' 
  | 'SHIPPED' 
  | 'COMPLETED' 
  | 'CANCELLED' 
  | 'REFUNDED';

export interface OrderItem {
  item_id: string;
  product_id: string;
  sku_id: string;
  quantity: number;
  original_unit_price: number;
  final_unit_price: number;
  total_price: number;
  product_name: string;
  product_image: string;
  sku_attributes: string;
  promotions_snapshot: Record<string, any>;
}

export interface ShopOrder {
  shop_order_id: string;
  shop_order_code: string;
  shop_id: string;
  status: OrderStatus;
  subtotal: number;
  shipping_fee: number;
  total_discount: number;
  total_amount: number;
  shop_voucher_code: string;
  shop_voucher_discount: number;
  site_order_discount: number; // New field
  site_shipping_discount: number; // New field
  shipping_method: string;
  tracking_code: string;
  items: OrderItem[];
  created_at: string;
  updated_at: string;
  paid_at: string | null;
  processing_at: string | null;
  shipped_at: string | null;
  completed_at: string | null;
  cancelled_at: string | null;
}

export interface OrderListParams {
  page?: number;
  limit?: number;
  status?: OrderStatus;
}

// New API structure - Nested order with shop order
export interface OrderWithShop {
  order: OrderDetail;
  order_shop: ShopOrder;
}

export interface OrderListResponse {
  code: number;
  message: string;
  status: string;
  result: {
    currentPage: number;
    data: OrderWithShop[];
    pageSize: number;
    totalElements: number;
    totalPages: number;
  };
}

// Dùng cho Order Detail
export interface PaymentMethod {
  id: string;
  name: string;
  code: string;
  type: string;
  is_active: boolean;
}

export interface OrderDetail {
  order_id: string;
  order_code: string;
  user_id: string;
  status: OrderStatus;
  grand_total: number;
  subtotal: number;
  total_shipping_fee: number;
  total_discount: number;
  shipping_address: string | ShippingAddress; // Có thể là string hoặc object
  payment_method: string | PaymentMethod; // Có thể là string hoặc object
  note: string;
  created_at: string;
  updated_at: string;
  site_order_voucher_code: string;
  site_shipping_voucher_code: number;
  paid_at: string | null;
  processing_at: string | null;
  shipped_at: string | null;
  completed_at: string | null;
  cancelled_at: string | null;
}

export interface OrderDetailResponse {
  code: number;
  message: string;
  status: string;
  result: {
    order: OrderDetail;
    order_shop: ShopOrder;
  };
}
