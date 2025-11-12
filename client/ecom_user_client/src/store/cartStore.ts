import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { CartItem } from '@/types/cart.types';

interface CartStore {
  items: CartItem[];
  addToCart: (item: CartItem) => void;
  removeFromCart: (sku_id: string) => void;
  updateQuantity: (sku_id: string, quantity: number) => void;
  clearCart: () => void;
  getTotalItems: () => number;
  getTotalPrice: () => number;
}

export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      items: [],

      // Thêm sản phẩm vào giỏ hàng
      addToCart: (item: CartItem) => {
        set((state) => {
          const existingItem = state.items.find((i) => i.sku_id === item.sku_id);

          if (existingItem) {
            // Nếu sản phẩm đã có trong giỏ, tăng số lượng
            return {
              items: state.items.map((i) =>
                i.sku_id === item.sku_id
                  ? { ...i, quantity: i.quantity + item.quantity }
                  : i
              ),
            };
          } else {
            // Nếu sản phẩm chưa có, thêm mới
            return {
              items: [...state.items, item],
            };
          }
        });
      },

      // Xóa sản phẩm khỏi giỏ hàng
      removeFromCart: (sku_id: string) => {
        set((state) => ({
          items: state.items.filter((item) => item.sku_id !== sku_id),
        }));
      },

      // Cập nhật số lượng sản phẩm
      updateQuantity: (sku_id: string, quantity: number) => {
        set((state) => ({
          items: state.items.map((item) =>
            item.sku_id === sku_id ? { ...item, quantity } : item
          ),
        }));
      },

      // Xóa toàn bộ giỏ hàng
      clearCart: () => {
        set({ items: [] });
      },

      // Lấy tổng số lượng sản phẩm
      getTotalItems: () => {
        return get().items.reduce((total, item) => total + item.quantity, 0);
      },

      // Lấy tổng giá trị giỏ hàng
      getTotalPrice: () => {
        return get().items.reduce((total, item) => total + item.price * item.quantity, 0);
      },
    }),
    {
      name: 'cart-storage', // Tên key trong localStorage
      storage: createJSONStorage(() => localStorage), // Sử dụng localStorage
      skipHydration: true, // Prevent automatic hydration to avoid SSR mismatch
    }
  )
);
