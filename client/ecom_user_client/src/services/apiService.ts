import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from "@tanstack/react-query";
import apiClient from "@/lib/apiClient";
import apiOrderClient from "@/lib/apiOrderService";
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
                customBaseURL:process.env.NEXT_PUBLIC_API_GATEWAY_URL
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
        customBaseURL:process.env.NEXT_PUBLIC_API_GATEWAY_URL

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
 * GET /orders
 */
export const useGetOrders = (params: OrderListParams) => {
  return useQuery<OrderListResponse['result'], Error>({
    queryKey: ['orders', params],
    queryFn: async () => {
      const cleanParams: Record<string, any> = {
        page: params.page || 1,
        limit: params.limit || 10,
      };

      if (params.status) cleanParams.status = params.status;

      const response = await apiOrderClient.get<OrderListResponse>('/orders', {
        params: cleanParams,
      });
      return response.data.result;
    },
    staleTime: 1000 * 60 * 1, // Cache 1 phút
  });
};

/**
 * Hook để lấy danh sách đơn hàng với infinite scroll
 * GET /orders
 */
export const useGetOrdersInfinite = (params: Omit<OrderListParams, 'page'>) => {
  return useInfiniteQuery<OrderListResponse['result'], Error>({
    queryKey: ['orders-infinite', params],
    queryFn: async ({ pageParam = 1 }) => {
      const cleanParams: Record<string, any> = {
        page: pageParam,
        limit: params.limit || 10,
      };

      if (params.status) cleanParams.status = params.status;

      const response = await apiOrderClient.get<OrderListResponse>('/orders', {
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
