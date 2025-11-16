// Dùng cho Trang Danh sách
export interface ProductSummary {
  id: string; // product_id
  name: string;
  key: string; // Dùng cho URL (slug)
  image: string;
  shop_id: string;
  brand_id: string;
  category_id: string;
  min_price: number;
  max_price: number;
  min_price_sku_id: string;
  max_price_sku_id: string;
  description: string;
  short_description: string;
  media: string | null;
  product_is_permission_check: boolean;
  product_is_permission_return: boolean;
  delete_status: string;
  create_date: string;
  update_date: string;
  rating:{
product_id:string;
total_reviews:number;
average_rating:number;
  };
}

export interface PaginatedProductsResponse {
  code: number;
  message: string;
  status: string;
  result: {
    currentPage: number;
    data: ProductSummary[];
    limit: number;
    totalElements: number;
    totalPages: number;
  };
}

// Dùng cho Trang Chi tiết
export interface ProductOptionValue {
  option_value_id: string;
  value: string;
  image?: string | null;
}

export interface ProductOption {
  option_name: string;
  values: ProductOptionValue[];
}

export interface ProductSKU {
  id: string; // SKU ID
  option_value_ids: string[];
  price: number;
  quantity: number; // Tồn kho
  sku_code: string;
  weight: number;
  sku_name:string;
}

export interface ProductBrand {
  brand_id: string;
  code: string;
  name: string;
  image: string | null;
  create_date: string;
  update_date: string;
}

export interface ProductCategory {
  category_id: string;
  key: string;
  name: string;
  image: string | null;
  parent: string | null;
  path: string;
}

export interface ProductInfo {
  id: string;
  key: string;
  name: string;
  description: string;
  short_description: string;
  image: string;
  media: string; // CSV string of media URLs
  shop_id: string;
  brand_id: string;
  category_id: string;
  min_price: number;
  max_price: number;
  min_price_sku_id: string;
  max_price_sku_id: string;
  delete_status: string;
  product_is_permission_check: boolean;
  product_is_permission_return: boolean;
  create_by: string;
  update_by: string;
  create_date: string;
  update_date: string;
}

export interface ProductDetailData {
  brand: ProductBrand;
  category: ProductCategory;
  option: ProductOption[];
  product: ProductInfo;
  sku: ProductSKU[];
}

export interface ProductDetailApiResponse {
  code: number;
  message: string;
  status: string;
  result: {
    data: ProductDetailData;
  };
}

// Params cho tìm kiếm sản phẩm
export interface ProductListParams {
  page?: number;
  limit?: number;
  keywords?: string;
  cate_path?: string;
  brand?: string;
  shop_id?: string;
  price_min?: number;
  price_max?: number;
  sort?: 'price_asc' | 'price_desc' | 'newest' | 'popular';
}
