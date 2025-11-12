/**
 * Utility functions for complaint page
 */

import ROUTER from "@/assets/configs/routers";

export type ComplaintCategory = "BUG" | "COMPLAINT" | "SUGGESTION" | "OTHER";

interface ComplaintParams {
  category?: ComplaintCategory;
  content?: string;
  phone?: string;
}

/**
 * Generate complaint page URL with pre-filled data
 * Used by Agent to redirect users to complaint form
 */
export const generateComplaintUrl = (params: ComplaintParams): string => {
  const baseUrl = ROUTER.khieunai;
  const searchParams = new URLSearchParams();

  if (params.category) {
    searchParams.append("category", params.category);
  }

  if (params.content) {
    searchParams.append("content", encodeURIComponent(params.content));
  }



  const queryString = searchParams.toString();
  return queryString ? `${baseUrl}?${queryString}` : baseUrl;
};

/**
 * Navigate to complaint page with pre-filled data
 * Can be called from Agent or other components
 */
export const navigateToComplaint = (params: ComplaintParams) => {
  const url = generateComplaintUrl(params);
  if (typeof window !== "undefined") {
    window.location.href = url;
  }
};

/**
 * Open complaint page in new tab with pre-filled data
 */
export const openComplaintInNewTab = (params: ComplaintParams) => {
  const url = generateComplaintUrl(params);
  if (typeof window !== "undefined") {
    window.open(url, "_blank");
  }
};

/**
 * Example usage in Agent:
 * 
 * // Redirect to complaint form with bug report
 * navigateToComplaint({
 *   category: "BUG",
 *   content: "Trang web bị lỗi khi tôi cố gắng thanh toán. Nút 'Xác nhận' không hoạt động."
 * });
 * 
 * // Generate URL only
 * const url = generateComplaintUrl({
 *   category: "COMPLAINT",
 *   content: "Sản phẩm không đúng mô tả"
 * });
 */
