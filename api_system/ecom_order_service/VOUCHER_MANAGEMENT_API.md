# API Quản Lý Voucher - Chi Tiết Endpoint

## 📋 Endpoint Information

### **GET** `/api/v1/vouchers/management`

**Mô tả**: Lấy danh sách voucher để quản lý
- **Admin** (ROLE_ADMIN) → Chỉ xem voucher PLATFORM (của sàn)
- **Seller** (ROLE_SELLER) → Chỉ xem voucher SHOP (của shop mình)

**Authentication**: Required (Bearer Token)

**Authorization**: ROLE_ADMIN hoặc ROLE_SELLER

---

## 🔑 Request Headers

```http
Authorization: Bearer {your_jwt_token}
Content-Type: application/json
```

---

## 📊 Query Parameters

### **Pagination (Bắt buộc có giá trị hợp lệ)**

| Parameter | Type | Required | Default | Max | Description |
|-----------|------|----------|---------|-----|-------------|
| `page` | int | No | 1 | - | Trang hiện tại (≥ 1) |
| `page_size` | int | No | 20 | 100 | Số voucher mỗi trang |

### **Search Filters (Tùy chọn)**

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `voucher_code` | string | Tìm theo mã voucher (LIKE search) | `SAVE50K` |
| `name` | string | Tìm theo tên voucher (LIKE search) | `Giảm giá` |

### **Attribute Filters (Tùy chọn)**

| Parameter | Type | Allowed Values | Description |
|-----------|------|----------------|-------------|
| `discount_type` | string | `PERCENTAGE`, `FIXED_AMOUNT` | Loại giảm giá |
| `applies_to_type` | string | `ORDER_TOTAL`, `SHIPPING_FEE` | Áp dụng cho |
| `audience_type` | string | `PUBLIC`, `ASSIGNED` | Đối tượng sử dụng |
| `is_active` | bool | `true`, `false` | Trạng thái kích hoạt |

### **Status Filter (Tùy chọn - Tính toán động)**

| Parameter | Type | Allowed Values | Description |
|-----------|------|----------------|-------------|
| `status` | string | `ACTIVE`, `EXPIRED`, `UPCOMING`, `DEPLETED` | Trạng thái voucher |

**Chi tiết Status:**
- `ACTIVE`: Đang hoạt động (is_active=true, trong thời gian hiệu lực, còn số lượng)
- `EXPIRED`: Đã hết hạn (end_date < now)
- `UPCOMING`: Sắp diễn ra (start_date > now)
- `DEPLETED`: Đã hết lượt (used_quantity >= total_quantity)

### **Sorting (Tùy chọn)**

| Parameter | Type | Default | Allowed Values |
|-----------|------|---------|----------------|
| `sort_by` | string | `created_at_desc` | `created_at_desc`, `created_at_asc`, `start_date_desc`, `start_date_asc`, `end_date_desc`, `end_date_asc` |

---

## 📝 Ví Dụ Request URLs

### 1. **Admin - Lấy tất cả voucher PLATFORM đang hoạt động**
```bash
GET /api/v1/vouchers/management?status=ACTIVE&page=1&page_size=20
```

### 2. **Seller - Tìm voucher theo mã**
```bash
GET /api/v1/vouchers/management?voucher_code=FREESHIP&page=1&page_size=10
```

### 3. **Admin - Lọc voucher giảm % còn hiệu lực, sắp xếp theo thời gian kết thúc**
```bash
GET /api/v1/vouchers/management?discount_type=PERCENTAGE&status=ACTIVE&sort_by=end_date_asc&page=1&page_size=50
```

### 4. **Seller - Xem voucher sắp hết hạn của shop**
```bash
GET /api/v1/vouchers/management?status=ACTIVE&sort_by=end_date_asc&page=1&page_size=20
```

### 5. **Admin - Xem voucher đã hết lượt**
```bash
GET /api/v1/vouchers/management?status=DEPLETED&page=1&page_size=20
```

### 6. **Seller - Tìm voucher theo tên, chỉ xem voucher công khai**
```bash
GET /api/v1/vouchers/management?name=Giảm&audience_type=PUBLIC&page=1&page_size=15
```

### 7. **Admin - Lọc voucher miễn phí vận chuyển đang tắt**
```bash
GET /api/v1/vouchers/management?applies_to_type=SHIPPING_FEE&is_active=false&page=1&page_size=20
```

---

## ✅ Response Format

### Success Response (200 OK)

```json
{
  "status": 200,
  "message": "Get vouchers successfully",
  "data": {
    "data": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Giảm 50K cho đơn hàng",
        "voucher_code": "SAVE50K",
        "owner_type": "PLATFORM",
        "owner_id": "admin-uuid-123",
        "discount_type": "FIXED_AMOUNT",
        "discount_value": "50000.00",
        "max_discount_amount": null,
        "applies_to_type": "ORDER_TOTAL",
        "min_purchase_amount": "200000.00",
        "audience_type": "PUBLIC",
        "start_date": "2025-01-01T00:00:00Z",
        "end_date": "2025-12-31T23:59:59Z",
        "total_quantity": 1000,
        "used_quantity": 350,
        "remaining_quantity": 650,
        "max_usage_per_user": 3,
        "is_active": true,
        "status": "ACTIVE",
        "created_at": "2024-12-01T10:00:00Z",
        "updated_at": "2024-12-15T14:30:00Z"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Freeship toàn quốc",
        "voucher_code": "FREESHIP99",
        "owner_type": "SHOP",
        "owner_id": "shop-uuid-456",
        "discount_type": "PERCENTAGE",
        "discount_value": "100.00",
        "max_discount_amount": "30000.00",
        "applies_to_type": "SHIPPING_FEE",
        "min_purchase_amount": "0.00",
        "audience_type": "PUBLIC",
        "start_date": "2025-01-01T00:00:00Z",
        "end_date": "2025-06-30T23:59:59Z",
        "total_quantity": 500,
        "used_quantity": 125,
        "remaining_quantity": 375,
        "max_usage_per_user": 5,
        "is_active": true,
        "status": "ACTIVE",
        "created_at": "2024-11-15T08:30:00Z",
        "updated_at": "2024-12-10T11:20:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 20,
      "total_items": 150,
      "total_pages": 8
    }
  }
}
```

### Error Responses

#### 400 Bad Request - Invalid Parameters
```json
{
  "status": 400,
  "message": "Invalid query parameters: page must be greater than 0"
}
```

#### 400 Bad Request - Invalid Filter Value
```json
{
  "status": 400,
  "message": "discount_type không hợp lệ. Allowed: PERCENTAGE, FIXED_AMOUNT"
}
```

#### 401 Unauthorized - Missing Token
```json
{
  "status": 401,
  "message": "Authorization token required"
}
```

#### 403 Forbidden - Invalid Role
```json
{
  "status": 403,
  "message": "Access denied. Only Admin and Seller can access this endpoint"
}
```

#### 500 Internal Server Error
```json
{
  "status": 500,
  "message": "lỗi khi lấy danh sách voucher: database connection error"
}
```

---

## 🎯 Response Fields Explanation

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | UUID của voucher |
| `name` | string | Tên hiển thị voucher |
| `voucher_code` | string | Mã voucher (unique) |
| `owner_type` | string | PLATFORM (sàn) hoặc SHOP (shop) |
| `owner_id` | string | ID của owner (admin hoặc seller) |
| `discount_type` | string | PERCENTAGE (giảm %) hoặc FIXED_AMOUNT (giảm số tiền cố định) |
| `discount_value` | string | Giá trị giảm (VD: 5 cho 5%, hoặc 50000 cho 50k) |
| `max_discount_amount` | string/null | Số tiền giảm tối đa (chỉ áp dụng cho PERCENTAGE) |
| `applies_to_type` | string | ORDER_TOTAL (tổng đơn) hoặc SHIPPING_FEE (phí ship) |
| `min_purchase_amount` | string | Giá trị đơn hàng tối thiểu |
| `audience_type` | string | PUBLIC (công khai) hoặc ASSIGNED (chỉ định) |
| `start_date` | timestamp | Thời gian bắt đầu |
| `end_date` | timestamp | Thời gian kết thúc |
| `total_quantity` | int | Tổng số lượt có thể sử dụng |
| `used_quantity` | int | Số lượt đã sử dụng |
| `remaining_quantity` | int | **Calculated**: Số lượt còn lại (total - used) |
| `max_usage_per_user` | int | Số lượt tối đa mỗi user |
| `is_active` | bool | Trạng thái kích hoạt |
| `status` | string | **Calculated**: ACTIVE/EXPIRED/UPCOMING/DEPLETED |
| `created_at` | timestamp | Thời gian tạo |
| `updated_at` | timestamp | Thời gian cập nhật |

---

## 💡 Lưu Ý Quan Trọng

### 1. **Phân quyền tự động**
- Không cần truyền `owner_type` trong query params
- Hệ thống tự động xác định dựa trên JWT token:
  - `ROLE_ADMIN` → Chỉ xem voucher `PLATFORM`
  - `ROLE_SELLER` → Chỉ xem voucher `SHOP` của mình

### 2. **Pagination**
- `page_size` tối đa là **100**
- Nếu không truyền, mặc định `page=1`, `page_size=20`
- Luôn kiểm tra `total_pages` trong response để biết có trang tiếp theo không

### 3. **Search với LIKE**
- `voucher_code` và `name` sử dụng LIKE search (không phân biệt hoa thường)
- Ví dụ: `voucher_code=SAVE` sẽ tìm cả `SAVE50K`, `SAVEBIG`, `SUPERSAVE`

### 4. **Status Filter (Quan trọng!)**
- `status` là trường **tính toán động**, không lưu trong database
- Có thể kết hợp với `is_active` để lọc chính xác hơn:
  - `status=ACTIVE` + `is_active=true`: Voucher đang chạy
  - `status=ACTIVE` + `is_active=false`: Voucher còn hạn nhưng đã tắt

### 5. **Sort By**
- Mặc định sắp xếp theo `created_at DESC` (mới nhất trước)
- Để tìm voucher sắp hết hạn: `sort_by=end_date_asc`
- Để tìm voucher mới tạo: `sort_by=created_at_desc`

### 6. **Performance**
- Đề xuất tạo index cho database:
```sql
CREATE INDEX idx_vouchers_owner ON vouchers(owner_id, owner_type);
CREATE INDEX idx_vouchers_dates ON vouchers(start_date, end_date);
CREATE INDEX idx_vouchers_status ON vouchers(is_active, used_quantity, total_quantity);
CREATE INDEX idx_vouchers_code ON vouchers(voucher_code);
CREATE INDEX idx_vouchers_name ON vouchers(name);
```

### 7. **Kết hợp Filters**
Có thể kết hợp nhiều filters cùng lúc:
```bash
GET /api/v1/vouchers/management?discount_type=PERCENTAGE&applies_to_type=ORDER_TOTAL&status=ACTIVE&is_active=true&sort_by=end_date_asc&page=1&page_size=30
```
→ Lấy voucher giảm % cho tổng đơn, đang hoạt động, sắp xếp theo ngày hết hạn

---

## 🧪 Test với cURL

### Admin Test
```bash
curl -X GET "http://localhost:8080/api/v1/vouchers/management?status=ACTIVE&page=1&page_size=20" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### Seller Test
```bash
curl -X GET "http://localhost:8080/api/v1/vouchers/management?voucher_code=SAVE&sort_by=created_at_desc&page=1&page_size=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

---

## 📱 Frontend Integration Example (JavaScript/TypeScript)

```typescript
interface VoucherManagementParams {
  page?: number;
  page_size?: number;
  voucher_code?: string;
  name?: string;
  discount_type?: 'PERCENTAGE' | 'FIXED_AMOUNT';
  applies_to_type?: 'ORDER_TOTAL' | 'SHIPPING_FEE';
  audience_type?: 'PUBLIC' | 'ASSIGNED';
  is_active?: boolean;
  status?: 'ACTIVE' | 'EXPIRED' | 'UPCOMING' | 'DEPLETED';
  sort_by?: 'created_at_desc' | 'created_at_asc' | 'start_date_desc' | 
            'start_date_asc' | 'end_date_desc' | 'end_date_asc';
}

async function getVouchersForManagement(params: VoucherManagementParams) {
  const queryParams = new URLSearchParams();
  
  // Add all non-null params
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      queryParams.append(key, String(value));
    }
  });

  const response = await fetch(
    `${API_BASE_URL}/api/v1/vouchers/management?${queryParams}`,
    {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessToken()}`,
        'Content-Type': 'application/json',
      },
    }
  );

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message);
  }

  return await response.json();
}

// Usage examples:
// 1. Get all active vouchers
const activeVouchers = await getVouchersForManagement({
  status: 'ACTIVE',
  page: 1,
  page_size: 20,
});

// 2. Search by voucher code
const searchResults = await getVouchersForManagement({
  voucher_code: 'SAVE',
  page: 1,
  page_size: 10,
});

// 3. Get expiring vouchers
const expiringVouchers = await getVouchersForManagement({
  status: 'ACTIVE',
  sort_by: 'end_date_asc',
  page: 1,
  page_size: 20,
});
```

---

## 🔍 Troubleshooting

### Vấn đề: Không thấy voucher nào
**Kiểm tra:**
- JWT token có hợp lệ không?
- User có role ADMIN hoặc SELLER không?
- Admin: Có voucher nào với `owner_type=PLATFORM` không?
- Seller: Có voucher nào với `owner_type=SHOP` và `owner_id=<seller_id>` không?

### Vấn đề: Status filter không hoạt động
**Lưu ý:**
- `status` là trường tính toán, không lưu trong DB
- Kiểm tra ngày giờ server có chính xác không
- Kết hợp với `is_active` để lọc chính xác hơn

### Vấn đề: Kết quả không đầy đủ
**Giải pháp:**
- Tăng `page_size` (max 100)
- Kiểm tra `total_items` và `total_pages` trong response
- Gọi API với `page` tiếp theo nếu cần

---

## 📞 Support

Nếu gặp vấn đề, liên hệ:
- Backend Team
- Tạo issue trong repository với tag `voucher-api`
