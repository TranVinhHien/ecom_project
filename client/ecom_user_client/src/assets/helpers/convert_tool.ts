
  // Helper để xử lý image URL
  const getImageUrl = (imageUrl: string | null | undefined) => {
    if (!imageUrl) return "/placeholder.png";
    // Nếu là URL đầy đủ (http/https) thì giữ nguyên
    if (imageUrl.startsWith("http://") || imageUrl.startsWith("https://")) {
      return imageUrl;
    }
    // Nếu không có protocol, thêm http:// vào đầu
    return `http://${imageUrl}`;
  };
  // Format giá
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };
export { getImageUrl, formatPrice };
