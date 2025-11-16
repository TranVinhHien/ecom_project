# Comment Service API Documentation

Service này xử lý các chức năng liên quan đến đánh giá/bình luận sản phẩm trong hệ thống ecommerce.

## Tổng quan

Comment Service cung cấp 3 API chính:
1. **Tạo đánh giá sản phẩm** - Cho phép người dùng đánh giá sản phẩm đã mua
2. **Lấy danh sách đánh giá** - Hiển thị tất cả đánh giá của sản phẩm
3. **Kiểm tra trạng thái đánh giá** - Kiểm tra order items nào đã được đánh giá

## API Endpoints

### 1. Tạo đánh giá sản phẩm

**Endpoint:** `POST /api/v1/comments`

**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "order_item_id": "uuid-string",
  "comment": "Sản phẩm rất tốt, giao hàng nhanh",
  "star": 5,
  "title": "Tuyệt vời!" // optional
}
```

**Validation Rules:**
- `order_item_id`: Required, phải là UUID hợp lệ
- `comment`: Required
- `star`: Required, phải từ 1-5
- `title`: Optional

**Business Rules:**
- User phải đã mua sản phẩm (order_item_id hợp lệ)
- Đơn hàng phải ở trạng thái COMPLETED
- Mỗi order_item chỉ được đánh giá 1 lần duy nhất
- User_id được lấy từ JWT token

**Response Success (201):**
```json
{
  "status": "success",
  "message": "Đánh giá sản phẩm thành công",
  "data": null
}
```

**Response Error Examples:**

*403 Forbidden - Không có quyền đánh giá:*
```json
{
  "status": "error",
  "message": "bạn không có quyền đánh giá sản phẩm này. Chỉ có thể đánh giá sau khi đơn hàng hoàn thành"
}
```

*409 Conflict - Đã đánh giá rồi:*
```json
{
  "status": "error",
  "message": "bạn đã đánh giá sản phẩm này rồi. Mỗi sản phẩm chỉ được đánh giá 1 lần"
}
```

---

### 2. Lấy danh sách đánh giá

**Endpoint:** `GET /api/v1/comments`

**Authentication:** Not required (Public API)

**Query Parameters:**
- `product_id` (required): UUID của sản phẩm
- `limit` (required): Số lượng comments mỗi trang (1-100)
- `offset` (optional): Vị trí bắt đầu (default: 0)

**Example:**
```
GET /api/v1/comments?product_id=abc-123&limit=20&offset=0
```

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Lấy danh sách đánh giá thành công",
  "data": {
    "comments": [
      {
        "comment_id": "uuid-1",
        "order_item_id": "uuid-order-item",
        "product_id": "abc-123",
        "sku_id": "sku-456",
        "user_id": "user-789",
        "sku_name_snapshot": null,
        "rating": 5,
        "title": "Tuyệt vời!",
        "content": "Sản phẩm rất tốt",
        "media": null,
        "parent_id": null,
        "created_at": "2025-11-13T10:00:00Z",
        "updated_at": "2025-11-13T10:00:00Z",
        "children": [
          {
            "comment_id": "uuid-2",
            "order_item_id": "uuid-order-item-2",
            "product_id": "abc-123",
            "sku_id": "sku-456",
            "user_id": "shop-owner-id",
            "rating": 0,
            "title": null,
            "content": "Cảm ơn bạn đã tin tùng shop!",
            "media": null,
            "parent_id": "uuid-1",
            "created_at": "2025-11-13T11:00:00Z",
            "updated_at": "2025-11-13T11:00:00Z"
          }
        ]
      }
    ],
    "stats": {
      "total_reviews": 150,
      "average_rating": 4.5
    },
    "pagination": {
      "limit": 20,
      "offset": 0
    }
  }
}
```

**Notes:**
- Comments được sắp xếp theo `created_at DESC` (mới nhất trước)
- Children (replies) chỉ có 1 cấp
- Children được sắp xếp theo `created_at ASC` (cũ nhất trước)

---

### 3. Kiểm tra trạng thái đánh giá

**Endpoint:** `POST /api/v1/comments/check-reviewed`

**Authentication:** Not required

**Use Case:** API này dành cho Order Service gọi để kiểm tra xem những order_item nào đã được đánh giá, từ đó hiển thị nút "Đánh giá" hoặc "Đã đánh giá" cho user.

**Request Body:**
```json
{
  "order_item_ids": [
    "uuid-item-1",
    "uuid-item-2",
    "uuid-item-3"
  ]
}
```

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Kiểm tra trạng thái đánh giá thành công",
  "data": {
    "reviewed_order_item_ids": [
      "uuid-item-1",
      "uuid-item-3"
    ]
  }
}
```

**Notes:**
- Chỉ trả về những order_item_id đã có đánh giá
- Nếu order_item chưa được đánh giá, không có trong response
- Empty array nếu không có item nào được đánh giá

---

## Database Schema

### Table: `product_comment`

```sql
CREATE TABLE `product_comment` (
  `comment_id` CHAR(36) PRIMARY KEY,
  `order_item_id` CHAR(36) NOT NULL UNIQUE, -- Đảm bảo mỗi item chỉ review 1 lần
  `product_id` VARCHAR(36) NOT NULL,
  `sku_id` VARCHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `sku_name_snapshot` NVARCHAR(500) DEFAULT NULL,
  `rating` TINYINT NOT NULL, -- 1-5 sao
  `title` NVARCHAR(255) DEFAULT NULL,
  `content` TEXT,
  `media` JSON DEFAULT NULL,
  `parent_id` CHAR(36) DEFAULT NULL, -- Cho phép reply
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  KEY `idx_user_id` (`user_id`),
  KEY `idx_sku_id` (`sku_id`)
);
```

## Security & Business Logic

### Quy trình tạo đánh giá:

1. **Authentication Check**
   - Verify JWT token
   - Extract user_id từ token

2. **Permission Check** (CheckReviewPermission)
   ```
   - order_item_id có tồn tại?
   - order_item này có thuộc về user_id không?
   - Shop order chứa item này có status = 'COMPLETED' không?
   ```

3. **Duplicate Check** (GetCommentByOrderItemID)
   ```
   - order_item_id này đã được review chưa?
   - UNIQUE constraint đảm bảo không thể review 2 lần
   ```

4. **Create Comment**
   - Generate UUID cho comment_id
   - Insert vào database
   - Return success

### Reply Comment (TODO - Chưa implement)

Hiện tại logic reply comment đã được chuẩn bị (parent_id field) nhưng chưa được implement hoàn chỉnh. Khi implement cần:
- Kiểm tra parent_id có tồn tại không
- Chỉ cho phép shop owner reply vào comment của sản phẩm họ bán
- Hoặc cho phép admin reply

## Error Handling

Service sử dụng `ServiceError` struct với 2 fields:
- `Code`: HTTP status code
- `Err`: error object

Common errors:
- 400 Bad Request: Invalid input
- 401 Unauthorized: Missing/invalid token
- 403 Forbidden: Không có quyền đánh giá
- 404 Not Found: Resource không tồn tại
- 409 Conflict: Đã đánh giá rồi
- 500 Internal Server Error: Database/server errors

## Files Structure

```
ecom_order_service/
├── db/
│   ├── query/
│   │   └── product_comment.sql          # SQL queries
│   └── sqlc/
│       └── product_comment.sql.go       # Generated code
├── services/
│   ├── entity/
│   │   └── comment_entity.go            # Request/Response structs
│   ├── interface/
│   │   └── usecase.go                   # Interface definitions
│   ├── comment.service.go               # Business logic
│   └── interface.go                     # Service interface aggregation
└── controllers/
    ├── comment.controllers.go           # HTTP handlers
    └── router.go                        # Route registration
```

## Testing

### Manual Testing với cURL

**1. Tạo đánh giá:**
```bash
curl -X POST http://localhost:8080/api/v1/comments \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_item_id": "uuid-here",
    "comment": "Sản phẩm tốt",
    "star": 5,
    "title": "Excellent"
  }'
```

**2. Lấy danh sách đánh giá:**
```bash
curl -X GET "http://localhost:8080/api/v1/comments?product_id=abc-123&limit=20&offset=0"
```

**3. Kiểm tra trạng thái:**
```bash
curl -X POST http://localhost:8080/api/v1/comments/check-reviewed \
  -H "Content-Type: application/json" \
  -d '{
    "order_item_ids": ["uuid-1", "uuid-2"]
  }'
```

## Future Enhancements

- [ ] Implement reply comment feature (parent_id logic)
- [ ] Add media upload support (images/videos)
- [ ] Add "helpful" voting system (review_likes table)
- [ ] Add admin moderation features
- [ ] Add filter by rating (1-5 stars)
- [ ] Add sort options (most helpful, most recent, etc.)
- [ ] Add pagination metadata (total_count, has_more, etc.)
