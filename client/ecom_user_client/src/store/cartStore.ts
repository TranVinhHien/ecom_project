import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { CartItem } from '@/types/cart.types';

interface CartStore {
  items: CartItem[];
  addToCart: (item: CartItem) => void;
  removeFromCart: (sku_id: string) => void;
  updateQuantity: (sku_id: string, quantity: number) => void;
  toggleSelection: (sku_id: string) => void;
  clearCart: () => void;
  getTotalItems: () => number;
  getTotalPrice: () => number;
  getSelectedTotalPrice: () => number;
  getSelectedItems: () => CartItem[];
  getItems: () => CartItem[];
}

/**
 * Cart Store cho người dùng chưa đăng nhập
 * Lưu trữ giỏ hàng trong localStorage
 * 
 * Khi người dùng đã đăng nhập, dữ liệu giỏ hàng sẽ được lấy từ API
 * và store này chỉ dùng làm backup tạm thời
 */
export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      items: [],

      // Thêm sản phẩm vào giỏ hàng (localStorage only - cho user chưa đăng nhập)
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
            // Nếu sản phẩm chưa có, thêm mới với isSelected = true mặc định
            return {
              items: [...state.items, { ...item, isSelected: true }],
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

      // Toggle selection của sản phẩm
      toggleSelection: (sku_id: string) => {
        set((state) => ({
          items: state.items.map((item) =>
            item.sku_id === sku_id ? { ...item, isSelected: !item.isSelected } : item
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

      // Lấy tổng giá trị các sản phẩm đã chọn
      getSelectedTotalPrice: () => {
        return get().items
          .filter(item => item.isSelected)
          .reduce((total, item) => total + item.price * item.quantity, 0);
      },

      // Lấy danh sách sản phẩm đã chọn
      getSelectedItems: () => {
        return get().items.filter(item => item.isSelected);
      },

      // Lấy danh sách items
      getItems: () => {
        return get().items;
      },
    }),
    {
      name: 'cart-storage', // Tên key trong localStorage
      storage: createJSONStorage(() => localStorage),
      skipHydration: true, // Prevent automatic hydration to avoid SSR mismatch
    }
  )
);
