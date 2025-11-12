import { create } from 'zustand';

export interface CheckoutItem {
  sku_id: string;
  shop_id: string;
  quantity: number;
  // Additional info for display (not sent to API)
  name?: string;
  price?: number;
  image?: string;
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
  
  // Shipping info
  shippingAddress: ShippingAddress | null;
  
  // Payment method (UUID)
  paymentMethod: string;
  
  // Vouchers
  vouchers: string[];
  
  // Note
  note: string;
  
  // Actions
  setCheckoutItems: (items: CheckoutItem[]) => void;
  setShippingAddress: (address: ShippingAddress) => void;
  setPaymentMethod: (method: string) => void;
  setVouchers: (vouchers: string[]) => void;
  setNote: (note: string) => void;
  clearCheckout: () => void;
}

export const useCheckoutStore = create<CheckoutStore>((set) => ({
  items: [],
  shippingAddress: null,
  paymentMethod: '',
  vouchers: [],
  note: '',

  setCheckoutItems: (items) => set({ items }),
  
  setShippingAddress: (address) => set({ shippingAddress: address }),
  
  setPaymentMethod: (method) => set({ paymentMethod: method }),
  
  setVouchers: (vouchers) => set({ vouchers }),
  
  setNote: (note) => set({ note }),
  
  clearCheckout: () => set({
    items: [],
    shippingAddress: null,
    paymentMethod: '',
    vouchers: [],
    note: '',
  }),
}));
