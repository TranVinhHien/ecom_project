# Controllers Documentation

## Tổng quan
Controllers này implement tất cả các API endpoint theo thiết kế trong `router.go` và `usecase.go`.

## Cấu trúc

### 1. Shop Controllers (`shop_controller.go`)
API dành cho Nhà bán hàng (SHOP)

**Authentication**: Yêu cầu JWT token với role = "SHOP"  
**shop_id**: Tự động lấy từ context thông qua middleware `getShopID()`

#### Nhóm 1: Tổng quan
- `GET /api/v1/shop/overview` - Tổng quan shop
  - Query: `start_date`, `end_date` (YYYY-MM-DD)
  
- `GET /api/v1/shop/wallet/summary` - Tóm tắt ví shop

#### Nhóm 2: Phân tích Đơn hàng
- `GET /api/v1/shop/orders` - Danh sách đơn hàng shop
  - Query: `status`, `start_date`, `end_date`, `limit`, `offset`
  
- `GET /api/v1/shop/orders/:shop_order_id` - Chi tiết đơn hàng
  - Param: `shop_order_id`
  
- `GET /api/v1/shop/order-items` - Danh sách sản phẩm trong đơn
  - Query: `product_id`, `start_date`, `end_date`, `limit`, `offset`

#### Nhóm 3: Phân tích Doanh thu & Dòng tiền
- `GET /api/v1/shop/revenue/timeseries` - Doanh thu theo thời gian
  - Query: `start_date`, `end_date`
  
- `GET /api/v1/shop/wallet/ledger-entries` - Lịch sử giao dịch ví
  - Query: `limit`, `offset`
  
- `GET /api/v1/shop/settlements` - Danh sách đối soát
  - Query: `status`, `start_date`, `end_date`, `limit`, `offset`

#### Nhóm 4: Phân tích Voucher
- `GET /api/v1/shop/vouchers` - Danh sách voucher
  - Query: `is_active` (true/false), `limit`, `offset`
  
- `GET /api/v1/shop/vouchers/performance` - Hiệu suất voucher
  - Query: `start_date`, `end_date`
  
- `GET /api/v1/shop/vouchers/:voucher_id/details` - Chi tiết sử dụng voucher
  - Param: `voucher_id`
  - Query: `limit`, `offset`

#### Nhóm 5: Xếp hạng
- `GET /api/v1/shop/ranking/products` - Xếp hạng sản phẩm
  - Query: `start_date`, `end_date`, `sort_by` (revenue|quantity), `limit`

---

### 2. Platform Controllers (`platform_controller.go`)
API dành cho Nền tảng (ADMIN)

**Authentication**: Yêu cầu JWT token với role = "ADMIN"

#### Nhóm 1: Tổng quan
- `GET /api/v1/platform/overview` - Tổng quan toàn sàn
  - Query: `start_date`, `end_date`

#### Nhóm 2: Quản lý Đơn hàng
- `GET /api/v1/platform/orders` - Danh sách đơn hàng
  - Query: `shop_id`, `user_id`, `status`, `start_date`, `end_date`, `limit`, `offset`
  
- `GET /api/v1/platform/orders/:order_id` - Chi tiết đơn hàng tổng
  - Param: `order_id`

#### Nhóm 3: Quản lý Tài chính
- `GET /api/v1/platform/finance/revenue-timeseries` - Doanh thu theo thời gian
  - Query: `start_date`, `end_date`
  
- `GET /api/v1/platform/finance/transactions` - Danh sách giao dịch
  - Query: `type`, `status`, `start_date`, `end_date`, `limit`, `offset`
  
- `GET /api/v1/platform/finance/settlements` - Danh sách đối soát
  - Query: `status`, `start_date`, `end_date`, `limit`, `offset`
  
- `GET /api/v1/platform/finance/ledgers` - Danh sách sổ cái
  - Query: `owner_type` (SHOP|PLATFORM), `limit`, `offset`
  
- `GET /api/v1/platform/finance/ledgers/:ledger_id/entries` - Các bút toán trong sổ cái
  - Param: `ledger_id`
  - Query: `limit`, `offset`

#### Nhóm 4: Phân tích Voucher
- `GET /api/v1/platform/vouchers` - Danh sách voucher
  - Query: `owner_type` (SHOP|PLATFORM), `is_active` (true/false), `limit`, `offset`
  
- `GET /api/v1/platform/vouchers/performance` - Hiệu suất voucher
  - Query: `start_date`, `end_date`

#### Nhóm 5: Phân tích Shop
- `GET /api/v1/platform/shops` - Danh sách shop
  - Query: `limit`, `offset`
  
- `GET /api/v1/platform/shops/:shop_id/detail` - Chi tiết shop
  - Param: `shop_id`
  - Query: `start_date`, `end_date`

#### Nhóm 6: Xếp hạng Toàn Sàn
- `GET /api/v1/platform/ranking/shops` - Top shop theo GMV
  - Query: `start_date`, `end_date`, `limit`
  
- `GET /api/v1/platform/ranking/products` - Top sản phẩm bán chạy
  - Query: `start_date`, `end_date`, `limit`
  
- `GET /api/v1/platform/ranking/users` - Top khách hàng chi tiêu
  - Query: `start_date`, `end_date`, `limit`
  
- `GET /api/v1/platform/ranking/categories` - Top danh mục
  - Query: `start_date`, `end_date`, `limit`

---

## Quy tắc chung

### Query Parameters
Tất cả các endpoint sử dụng **query parameters** cho việc filter và phân trang:

**Ngày tháng:**
- Format: `YYYY-MM-DD` (ví dụ: `2024-01-15`)
- Default `start_date`: 30 ngày trước
- Default `end_date`: hôm nay

**Phân trang:**
- `limit`: số lượng record (default: 20)
- `offset`: vị trí bắt đầu (default: 0)

**Boolean:**
- `is_active`: `true` hoặc `false`

### Response Format
Tất cả các response đều theo format chuẩn:

**Success Response:**
```json
{
  "code": 200,
  "status": "success",
  "message": "success",
  "result": {
    // Dữ liệu trả về
  }
}
```

**Error Response:**
```json
{
  "code": 400|401|403|404|500,
  "status": "error|authentication|forbidden|notfound",
  "message": "Mô tả lỗi"
}
```

### Middleware Chain

**Shop APIs:**
```
authorization(jwt) -> checkRole("SHOP") -> getShopID() -> controller
```

**Platform APIs:**
```
authorization(jwt) -> checkRole("ADMIN") -> controller
```

---

## Lưu ý khi sử dụng

1. **shop_id cho Shop APIs**: Không cần truyền qua query/param, tự động lấy từ context
2. **Date format**: Luôn sử dụng format `YYYY-MM-DD`
3. **Pagination**: Sử dụng `limit` và `offset` cho phân trang
4. **Filter**: Các filter là optional, không truyền sẽ lấy tất cả
5. **Authorization**: Tất cả API đều yêu cầu JWT token trong header:
   ```
   Authorization: Bearer <token>
   ```

## Testing

Ví dụ test với curl:

```bash
# Shop API - Get overview
curl -X GET "http://localhost:8080/api/v1/shop/overview?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer <shop_token>"

# Platform API - List orders
curl -X GET "http://localhost:8080/api/v1/platform/orders?shop_id=SHOP123&limit=10&offset=0" \
  -H "Authorization: Bearer <admin_token>"

# Platform API - Get order detail
curl -X GET "http://localhost:8080/api/v1/platform/orders/ORDER123" \
  -H "Authorization: Bearer <admin_token>"
```
