// Dùng cho Zustand Store (localStorage) - khi chưa đăng nhập
export interface CartItem {
  sku_id: string;
  shop_id: string;
  quantity: number;
  name: string;
  price: number;
  image: string;
  sku_name: string;
  isSelected: boolean; // Checkbox selection state
}

// ================ API Response Types ================

// Cart Item từ API
export interface ApiCartItem {
  skuId: string;
  productName: string;
  productImage: string;
  price: number;
  quantity: number;
  isSelected: boolean;
  shopId: string;
  addedDate: string;

}

// Cart Data từ API
export interface ApiCartData {
  id: string;
  items: ApiCartItem[];
  totalItems: number;
  totalPrice: number;
  selectedTotalPrice: number;
}

// Response structure từ API
export interface ApiCartResponse {
  result: ApiCartData;
  messages: string[];
  succeeded: boolean;
  code: number;
}

// Response cho count API
export interface ApiCartCountResponse {
  result: number;
  messages: string[];
  succeeded: boolean;
  code: number;
}

// Payload để add item vào cart
export interface AddToCartPayload {
  SkuId: string;
  Quantity: number;
}

// Payload để update quantity
export interface UpdateCartItemPayload {
  Quantity: number;
}
