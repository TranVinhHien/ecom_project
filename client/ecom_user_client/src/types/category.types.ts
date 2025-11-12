// Dựa trên response của API: /categories/get
export interface Category {
  category_id: string;
  name: string;
  key: string; // Example: "ao-thun-nam"
  path: string; // Example: "/thoi-trang/ao-thun-nam"
  child: {
    data: Category[] | null;
    valid: boolean;
  };
}

export interface CategoryApiResponse {
  result: { categories: Category[] };
}
