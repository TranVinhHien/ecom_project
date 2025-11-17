# 🛒 Cart Integration Guide

## Overview

The cart system has been completely refactored to support both:
- **Logged-in users**: Cart data stored in backend via API
- **Non-logged users**: Cart data stored in localStorage

### Key Features

✅ Automatic cart sync when user logs in
✅ Seamless switching between localStorage and API
✅ Real-time cart count updates
✅ React Query for efficient data fetching
✅ Optimistic UI updates

---

## Architecture

### 1. **Type Definitions** (`/src/types/cart.types.ts`)

```typescript
// LocalStorage cart item
interface CartItem {
  sku_id: string;
  shop_id: string;
  quantity: number;
  name: string;
  price: number;
  image: string;
  sku_name: string;
}

// API response types
interface ApiCartItem {
  skuId: string;
  productName: string;
  price: number;
  quantity: number;
  isSelected: boolean;
  shopId: string;
  addedDate: string;
}
```

### 2. **API Endpoints** (`/src/assets/configs/api.ts`)

```typescript
cart: {
  getCart: "/Cart",              // GET - Get cart
  addItem: "/Cart/items",        // POST - Add item
  updateItem: "/Cart/items",     // PUT - Update quantity (append /{skuId})
  deleteItem: "/Cart/items",     // DELETE - Remove item (append /{skuId})
  clearCart: "/Cart",            // DELETE - Clear cart
  getCount: "/Cart/count",       // GET - Get item count
}
```

### 3. **Services** (`/src/services/apiService.ts`)

React Query hooks for cart operations:
- `useGetCart()` - Fetch cart data
- `useGetCartCount()` - Fetch cart item count
- `useAddToCart()` - Add item to cart
- `useUpdateCartItem()` - Update item quantity
- `useDeleteCartItem()` - Remove item from cart
- `useClearCart()` - Clear entire cart

### 4. **Zustand Store** (`/src/store/cartStore.ts`)

For localStorage cart (non-logged users only):
```typescript
const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      items: [],
      addToCart: (item) => {...},
      removeFromCart: (sku_id) => {...},
      updateQuantity: (sku_id, quantity) => {...},
      clearCart: () => {...},
      getTotalItems: () => {...},
      getTotalPrice: () => {...},
    }),
    { name: 'cart-storage' }
  )
);
```

### 5. **Cart Sync Service** (`/src/lib/cartSyncService.ts`)

Handles syncing localStorage cart to API when user logs in:
```typescript
cartSyncService.syncLocalCartToAPI() // Called after successful login
```

### 6. **Custom Hook** (`/src/hooks/useCart.ts`)

**Recommended approach** - Abstracts cart logic:
```typescript
const { addToCart, getCartCount, isLoggedIn, isAddingToCart } = useCart();
```

---

## Usage Examples

### Adding Items to Cart

#### Option 1: Using `useCart` Hook (Recommended)

```tsx
import { useCart } from "@/hooks/useCart";

function ProductCard({ product }) {
  const { addToCart, isAddingToCart } = useCart();

  const handleAddToCart = async () => {
    const success = await addToCart({
      sku_id: product.skuId,
      shop_id: product.shopId,
      quantity: 1,
      name: product.name,
      price: product.price,
      image: product.image,
      sku_name: product.skuName,
    });

    if (success) {
      console.log("Item added successfully!");
    }
  };

  return (
    <Button onClick={handleAddToCart} disabled={isAddingToCart}>
      {isAddingToCart ? "Adding..." : "Add to Cart"}
    </Button>
  );
}
```

#### Option 2: Direct API Usage (For logged-in users)

```tsx
import { useAddToCart } from "@/services/apiService";

function ProductPage() {
  const addToCartMutation = useAddToCart();

  const handleAddToCart = async () => {
    try {
      await addToCartMutation.mutateAsync({
        SkuId: "product-sku-id",
        Quantity: 2,
      });
      toast({ title: "Added to cart!" });
    } catch (error) {
      toast({ title: "Error", variant: "destructive" });
    }
  };

  return <Button onClick={handleAddToCart}>Add to Cart</Button>;
}
```

#### Option 3: Direct localStorage Usage (For non-logged users)

```tsx
import { useCartStore } from "@/store/cartStore";

function ProductCard() {
  const addToLocalCart = useCartStore((state) => state.addToCart);

  const handleAddToCart = () => {
    addToLocalCart({
      sku_id: "sku-123",
      shop_id: "shop-456",
      quantity: 1,
      name: "Product Name",
      price: 100000,
      image: "image-url",
      sku_name: "SKU Name",
    });
  };

  return <Button onClick={handleAddToCart}>Add to Cart</Button>;
}
```

---

### Displaying Cart Count in Header

```tsx
import { useCart } from "@/hooks/useCart";

function Header() {
  const { getCartCount, isHydrated } = useCart();
  const count = getCartCount();

  return (
    <Button>
      <ShoppingCart />
      {isHydrated && count > 0 && (
        <span className="badge">{count}</span>
      )}
    </Button>
  );
}
```

---

### Cart Page

The cart page automatically detects login status and displays appropriate data:

```tsx
// /src/app/[locale]/(pages)/gio-hang/page.tsx
import { useGetCart } from "@/services/apiService";
import { useCartStore } from "@/store/cartStore";

export default function CartPage() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  
  // API cart (logged-in users)
  const { data: apiCart, isLoading } = useGetCart();
  
  // LocalStorage cart (non-logged users)
  const localItems = useCartStore((state) => state.items);

  // Determine which cart to display
  const displayItems = isLoggedIn ? (apiCart?.items || []) : localItems;

  return (
    <div>
      {displayItems.map(item => (
        <CartItem key={item.skuId || item.sku_id} item={item} />
      ))}
    </div>
  );
}
```

---

## Login Flow with Cart Sync

When a user logs in, the system automatically syncs localStorage cart to API:

```tsx
// /src/app/[locale]/auth/login/page.tsx
import { cartSyncService } from "@/lib/cartSyncService";

const handleLogin = async () => {
  // 1. Login user
  await loginMutation.mutateAsync(credentials);
  
  // 2. Save token and user info
  cookies.setCookieValues(ACCESS_TOKEN, token);
  localStorage.setItem(INFO_USER, JSON.stringify(userData));
  
  // 3. Sync cart (automatic)
  if (cartSyncService.hasLocalCartItems()) {
    await cartSyncService.syncLocalCartToAPI();
    toast({ title: "Cart synced successfully!" });
  }
  
  // 4. Redirect
  router.push("/");
};
```

---

## API Response Examples

### Get Cart Response

```json
{
  "result": {
    "id": "9acc8610-314d-49cf-973d-31d29f371188",
    "items": [
      {
        "skuId": "019cbb88-0bfa-4389-a01a-98115af5613f",
        "productName": "Default",
        "price": 3190000,
        "quantity": 1,
        "isSelected": true,
        "shopId": "019cbb88-0bfa-4389-a01a-98115af5613f",
        "addedDate": "2025-11-16T07:16:16.780055Z"
      }
    ],
    "totalItems": 1,
    "totalPrice": 3190000,
    "selectedTotalPrice": 3190000
  },
  "succeeded": true,
  "code": 200
}
```

### Get Cart Count Response

```json
{
  "result": 3,
  "succeeded": true,
  "code": 200
}
```

---

## Testing Checklist

### Non-logged User Flow
- [ ] Add product to cart → Saved to localStorage
- [ ] View cart page → Shows localStorage items
- [ ] Header shows correct count from localStorage
- [ ] Navigate away and back → Cart persists

### Logged-in User Flow
- [ ] Add product to cart → Sent to API
- [ ] View cart page → Shows API items
- [ ] Header shows correct count from API
- [ ] Update quantity → Updates via API
- [ ] Remove item → Removes via API
- [ ] Cart syncs across devices/browsers

### Login with Existing Cart
- [ ] Add items to localStorage (not logged in)
- [ ] Login → Items sync to API
- [ ] Toast notification shows sync success
- [ ] localStorage cart is cleared
- [ ] API cart contains all items

### Edge Cases
- [ ] Network error handling
- [ ] Concurrent cart updates
- [ ] Session expiry during cart operations
- [ ] Empty cart states
- [ ] Invalid product IDs

---

## Performance Considerations

1. **React Query Caching**: Cart data is cached for 30 seconds to reduce API calls
2. **Optimistic Updates**: UI updates immediately while API call is in progress
3. **Debouncing**: Quantity updates are debounced to avoid excessive API calls
4. **Lazy Loading**: Cart count fetched only when needed

---

## Troubleshooting

### Cart count not updating
- Check if `useGetCartCount()` is being called
- Verify query invalidation after mutations
- Check network requests in DevTools

### Cart sync fails on login
- Check cartSyncService console logs
- Verify API endpoints are correct
- Check token is valid

### localStorage cart not clearing
- Verify `clearCart()` is called after sync
- Check browser console for errors

---

## Future Enhancements

- [ ] Cart item selection (for checkout)
- [ ] Save for later functionality
- [ ] Cart expiration for non-logged users
- [ ] Real-time cart updates via WebSocket
- [ ] Cart recommendations
- [ ] Guest checkout with cart preservation

---

## Support

For issues or questions, contact the development team or create an issue in the project repository.
