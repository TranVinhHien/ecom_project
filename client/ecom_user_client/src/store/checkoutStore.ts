import { create } from 'zustand';

export interface CheckoutItem {
  sku_id: string;
  shop_id: string;
  quantity: number;
  // Additional info for display (not sent to API)
  name?: string;
  price?: number;
  image?: string;
  sku_name?: string;
}

export interface ShippingAddress {
  fullName: string;
  phone: string;
  address: string;
  district: string;
  city: string;
  postalCode: string;
}

interface CheckoutStore {
  // Checkout items
  items: CheckoutItem[];
  
  // Flag để phân biệt nguồn gốc:
  // - true: Từ giỏ hàng (sau khi thanh toán thành công sẽ xóa khỏi giỏ)
  // - false: Từ trang chi tiết (Mua ngay - không xóa gì)
  isFromCart: boolean;
  
  // Shipping info
  shippingAddress: ShippingAddress | null;
  
  // Payment method (UUID)
  paymentMethod: string;
  
  // Vouchers
  vouchers: string[];
  
  // Note
  note: string;
  
  // Actions
  setCheckoutItems: (items: CheckoutItem[], fromCart?: boolean) => void;
  setIsFromCart: (isFromCart: boolean) => void;
  setShippingAddress: (address: ShippingAddress) => void;
  setPaymentMethod: (method: string) => void;
  setVouchers: (vouchers: string[]) => void;
  setNote: (note: string) => void;
  clearCheckout: () => void;
}

export const useCheckoutStore = create<CheckoutStore>((set) => ({
  items: [],
  isFromCart: false, // Mặc định là không phải từ giỏ hàng
  shippingAddress: null,
  paymentMethod: '',
  vouchers: [],
  note: '',

  setCheckoutItems: (items, fromCart = false) => set({ items, isFromCart: fromCart }),
  
  setIsFromCart: (isFromCart) => set({ isFromCart }),
  
  setShippingAddress: (address) => set({ shippingAddress: address }),
  
  setPaymentMethod: (method) => set({ paymentMethod: method }),
  
  setVouchers: (vouchers) => set({ vouchers }),
  
  setNote: (note) => set({ note }),
  
  clearCheckout: () => set({
    items: [],
    isFromCart: false,
    shippingAddress: null,
    paymentMethod: '',
    vouchers: [],
    note: '',
  }),
}));
