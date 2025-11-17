/**
 * Cart Sync Service
 * 
 * This service handles syncing localStorage cart items to the API cart
 * when a user logs in.
 */

import { useCartStore } from "@/store/cartStore";
import apiCartClient from "./apiCartService";
import API from "@/assets/configs/api";
import { AddToCartPayload, ApiCartResponse } from "@/types/cart.types";

export const cartSyncService = {
  /**
   * Sync localStorage cart items to API cart after login
   * This function should be called after successful login
   */
  syncLocalCartToAPI: async (): Promise<void> => {
    try {
      console.log("🔄 Starting cart sync from localStorage to API...");
      
      // Get items from localStorage cart
      const localCartItems = useCartStore.getState().getItems();
      
      if (localCartItems.length === 0) {
        console.log("ℹ️ No items in localStorage cart to sync");
        return;
      }

      console.log(`📦 Found ${localCartItems.length} items in localStorage cart`);

      // Sync each item to API
      for (const item of localCartItems) {
        try {
          const payload: AddToCartPayload = {
            SkuId: item.sku_id,
            Quantity: item.quantity,
          };

          console.log(`➕ Adding item to API cart:`, payload);
          
          const response = await apiCartClient.post<ApiCartResponse>(
            API.cart.addItem,
            payload
          );

          if (response.data.succeeded) {
            console.log(`✅ Successfully synced item ${item.sku_id}`);
          }
        } catch (error) {
          console.error(`❌ Failed to sync item ${item.sku_id}:`, error);
          // Continue with other items even if one fails
        }
      }

      // Clear localStorage cart after successful sync
      console.log("🧹 Clearing localStorage cart after sync...");
      useCartStore.getState().clearCart();
      
      console.log("✅ Cart sync completed successfully");
    } catch (error) {
      console.error("❌ Error during cart sync:", error);
      throw error;
    }
  },

  /**
   * Check if user has items in localStorage cart
   */
  hasLocalCartItems: (): boolean => {
    const localCartItems = useCartStore.getState().getItems();
    return localCartItems.length > 0;
  },

  /**
   * Get count of items in localStorage cart
   */
  getLocalCartItemsCount: (): number => {
    return useCartStore.getState().getTotalItems();
  },
};
