import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from "@tanstack/react-query";
import apiClient from "@/lib/apiClient";
import apiOrderClient from "@/lib/apiOrderService";
import apiCartClient from "@/lib/apiCartService";
import { cookies } from "@/assets/helpers";
import { ACCESS_TOKEN } from "@/assets/configs/request";
import { 
  ProductListParams, 
  PaginatedProductsResponse, 
  ProductDetailData, 
  ProductDetailApiResponse 
} from "@/types/product.types";
import { Category, CategoryApiResponse } from "@/types/category.types";
import { 
  CreateOrderPayload, 
  CreateOrderSuccessResponse, 
  OrderListParams, 
  OrderListResponse,
  OrderDetailResponse
} from "@/types/order.types";
import { VoucherApiResponse } from "@/types/voucher.types";
import { BannerType, BannerApiResponse } from "@/types/shop.types";
import API from "@/assets/configs/api";
import { 
  ApiCartResponse, 
  ApiCartCountResponse, 
  AddToCartPayload, 
  UpdateCartItemPayload 
} from "@/types/cart.types";

// ================ CATEGORIES ================

/**
 * Hook để lấy danh sách categories
 * GET /categories/get
 */
export const useGetCategories = () => {
  return useQuery<Category[], Error>({
    queryKey: ['categories'],
    queryFn: async () => {
      const response = await apiClient.get<CategoryApiResponse>('/categories/get',{
                customBaseURL:API.base_product
      });
      return response.data.result.categories;
    },
    staleTime: 1000 * 60 * 5, // Cache 5 phút
  });
};

// ================ PRODUCTS ================

/**
 * Hook để lấy danh sách sản phẩm với phân trang và filters
 * GET /product/getall
 */
export const useGetProducts = (params: ProductListParams) => {
  return useQuery<PaginatedProductsResponse['result'], Error>({
    queryKey: ['products', params],
    queryFn: async () => {
      // Loại bỏ các params undefined/null
      const cleanParams: Record<string, any> = {
        page: params.page || 1,
        limit: params.limit || 20,
      };

      if (params.keywords) cleanParams.keywords = params.keywords;
      if (params.cate_path) cleanParams.cate_path = params.cate_path;
      if (params.brand) cleanParams.brand = params.brand;
      if (params.shop_id) cleanParams.shop_id = params.shop_id;
      if (params.price_min !== undefined) cleanParams.price_min = params.price_min;
      if (params.price_max !== undefined) cleanParams.price_max = params.price_max;
      if (params.sort) cleanParams.sort = params.sort;

      const response = await apiClient.get<PaginatedProductsResponse>('/product/getall', {
        params: cleanParams,
        customBaseURL:API.base_product

      },);
      return response.data.result;
    },
    staleTime: 1000 * 60 * 2, // Cache 2 phút
  });
};

/**
 * Hook để lấy chi tiết sản phẩm
 * GET /product/getdetail/:key
 */
export const useGetProductDetail = (key: string) => {
  return useQuery<ProductDetailData, Error>({
    queryKey: ['product', key],
    queryFn: async () => {
      const response = await apiClient.get<ProductDetailApiResponse>(
        `/product/getdetail/${key}`
      );
      return response.data.result.data;
    },
    enabled: !!key, // Chỉ fetch khi có key
    staleTime: 1000 * 60 * 3, // Cache 3 phút
  });
};

// ================ ORDERS ================

/**
 * Hook để tạo đơn hàng
 * POST /orders
 */
export const useCreateOrder = () => {
  return useMutation<CreateOrderSuccessResponse, Error, CreateOrderPayload>({
    mutationFn: async (payload: CreateOrderPayload) => {

      const response = await apiOrderClient.post<CreateOrderSuccessResponse>('/orders', payload);
      return response.data;
    },
  });
};

/**
 * Hook để lấy danh sách đơn hàng
 * GET /orders/search/detail
 */
export const useGetOrders = (params: OrderListParams) => {
  return useQuery<OrderListResponse['result'], Error>({
    queryKey: ['orders', params],
    queryFn: async () => {
      const cleanParams: Record<string, any> = {
        page: params.page || 1,
        page_size: params.limit || 10,
      };

      if (params.status) cleanParams.status = params.status;

      const response = await apiOrderClient.get<OrderListResponse>('/orders/search/detail', {
        params: cleanParams,
      });
      return response.data.result;
    },
    staleTime: 1000 * 60 * 1, // Cache 1 phút
  });
};

/**
 * Hook để lấy danh sách đơn hàng với infinite scroll
 * GET /orders/search/detail
 */
export const useGetOrdersInfinite = (params: Omit<OrderListParams, 'page'>) => {
  return useInfiniteQuery<OrderListResponse['result'], Error>({
    queryKey: ['orders-infinite', params],
    queryFn: async ({ pageParam = 1 }) => {
      const cleanParams: Record<string, any> = {
        page: pageParam,
        page_size: params.limit || 10,
      };

      if (params.status) cleanParams.status = params.status;

      const response = await apiOrderClient.get<OrderListResponse>('/orders/search/detail', {
        params: cleanParams,
      });
      return response.data.result;
    },
    getNextPageParam: (lastPage, allPages) => {
      // Nếu còn trang tiếp theo thì trả về số trang kế tiếp
      if (lastPage.currentPage < lastPage.totalPages) {
        return lastPage.currentPage + 1;
      }
      return undefined; // Không còn trang nào nữa
    },
    initialPageParam: 1,
    staleTime: 1000 * 60 * 1, // Cache 1 phút
  });
};

/**
 * Hook để lấy chi tiết đơn hàng
 * GET /orders/{orderId}
 */
export const useGetOrderDetail = (orderId: string) => {
  return useQuery<OrderDetailResponse['result'], Error>({
    queryKey: ['order-detail', orderId],
    queryFn: async () => {
      const response = await apiOrderClient.get<OrderDetailResponse>(`/orders/${orderId}`);
      return response.data.result;
    },
    enabled: !!orderId,
    staleTime: 1000 * 60 * 2, // Cache 2 phút
  });
};

// ================ VOUCHERS ================

/**
 * Hook để lấy danh sách vouchers
 * GET /vouchers
 */
export const useGetVouchers = () => {
  return useQuery<VoucherApiResponse['result'], Error>({
    queryKey: ['vouchers'],
    queryFn: async () => {
      const response = await apiOrderClient.get<VoucherApiResponse>('/vouchers');
      return response.data.result;
    },
    staleTime: 1000 * 60 * 5, // Cache 5 phút
  });
};

// ================ CART ================

/**
 * Hook để lấy giỏ hàng của người dùng đã đăng nhập
 * GET /Cart
 * Nếu không có token (chưa đăng nhập), sẽ skip API call
 */
export const useGetCart = () => {
  // Kiểm tra token trước khi gọi API
  const hasToken = () => {
    if (typeof window === 'undefined') return false;
    return !!cookies.getCookieValues<string>(ACCESS_TOKEN);
  };
  
  return useQuery<ApiCartResponse['result'], Error>({
    queryKey: ['cart'],
    queryFn: async () => {
      const response = await apiCartClient.get<ApiCartResponse>(API.cart.getCart);
      return response.data.result;
    },
    enabled: hasToken(), // Chỉ gọi API nếu có token
    staleTime: 0, // Không cache, luôn lấy dữ liệu mới
    retry: false, // Không retry nếu lỗi
  });
};

/**
 * Hook để lấy số lượng sản phẩm trong giỏ hàng
 * GET /Cart/count
 * Nếu không có token, sẽ skip API call
 */
export const useGetCartCount = () => {
  // Kiểm tra token trước khi gọi API
  const hasToken = () => {
    if (typeof window === 'undefined') return false;
    return !!cookies.getCookieValues<string>(ACCESS_TOKEN);
  };
  
  return useQuery<number, Error>({
    queryKey: ['cart-count'],
    queryFn: async () => {
      const response = await apiCartClient.get<ApiCartCountResponse>(API.cart.getCount);
      return response.data.result;
    },
    enabled: hasToken(), // Chỉ gọi API nếu có token
    // staleTime: 1000 * 30, // Cache 30 giây
    retry: false,
  });
};

/**
 * Hook để thêm sản phẩm vào giỏ hàng
 * POST /Cart/items
 */
export const useAddToCart = () => {
  const queryClient = useQueryClient();
  
  return useMutation<ApiCartResponse['result'], Error, AddToCartPayload>({
    mutationFn: async (payload: AddToCartPayload) => {
      const response = await apiCartClient.post<ApiCartResponse>(API.cart.addItem, payload);
      return response.data.result;
    },
    onSuccess: () => {
      // Invalidate cart queries để refetch
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      queryClient.invalidateQueries({ queryKey: ['cart-count'] });
    },
  });
};

/**
 * Hook để cập nhật số lượng sản phẩm trong giỏ
 * PUT /Cart/items/{skuId}
 */
export const useUpdateCartItem = () => {
  const queryClient = useQueryClient();
  
  return useMutation<ApiCartResponse['result'], Error, { skuId: string; payload: UpdateCartItemPayload }>({
    mutationFn: async ({ skuId, payload }) => {
      const response = await apiCartClient.put<ApiCartResponse>(
        `${API.cart.updateItem}/${skuId}`, 
        payload
      );
      return response.data.result;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      queryClient.invalidateQueries({ queryKey: ['cart-count'] });
    },
  });
};

/**
 * Hook để xóa sản phẩm khỏi giỏ hàng
 * DELETE /Cart/items/{skuId}
 */
export const useDeleteCartItem = () => {
  const queryClient = useQueryClient();
  
  return useMutation<ApiCartResponse['result'], Error, string>({
    mutationFn: async (skuId: string) => {
      const response = await apiCartClient.delete<ApiCartResponse>(
        `${API.cart.deleteItem}/${skuId}`
      );
      return response.data.result;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      queryClient.invalidateQueries({ queryKey: ['cart-count'] });
    },
  });
};

/**
 * Hook để xóa toàn bộ giỏ hàng
 * DELETE /Cart
 */
export const useClearCart = () => {
  const queryClient = useQueryClient();
  
  return useMutation<ApiCartResponse['result'], Error, void>({
    mutationFn: async () => {
      const response = await apiCartClient.delete<ApiCartResponse>(API.cart.clearCart);
      return response.data.result;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      queryClient.invalidateQueries({ queryKey: ['cart-count'] });
    },
  });
};

// ================ BANNERS ================

/**
 * Hook để lấy danh sách banners theo loại
 * GET /Banners/active
 */
export const useGetActiveBanners = (bannerType?: BannerType) => {
  return useQuery<BannerApiResponse['result'], Error>({
    queryKey: ['banners', bannerType],
    queryFn: async () => {
      const params: Record<string, any> = {};
      if (bannerType) params.bannerType = bannerType;

      const response = await apiClient.get<BannerApiResponse>('/Banners/active', {
        params,
        customBaseURL: 'http://localhost:8000/api'
      });
      return response.data.result;
    },
    staleTime: 1000 * 60 * 5, // Cache 5 phút
  });
};
