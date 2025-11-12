# 🛍️ E-Commerce Product Service

Microservice quản lý sản phẩm cho hệ thống thương mại điện tử, xử lý toàn bộ vòng đời sản phẩm từ tạo, cập nhật, quản lý danh mục đến upload media.

## 📋 Tổng quan

**Product Service** là một trong những microservice cốt lõi của hệ thống E-Commerce, chịu trách nhiệm:
- Quản lý sản phẩm (SPU) và biến thể sản phẩm (SKU)
- Quản lý danh mục sản phẩm phân cấp
- Quản lý thương hiệu (Brand)
- Quản lý thuộc tính sản phẩm (Option Values: Màu sắc, Size, ...)
- Upload và quản lý media (ảnh, video)
- Tìm kiếm và lọc sản phẩm nâng cao
- Tích hợp Redis cache để tối ưu hiệu suất

## 🏗️ Kiến trúc

### Tech Stack
- **Language**: Go 1.24.6
- **Framework**: Gin (HTTP Router)
- **Database**: MySQL 8.0+
- **Cache**: Redis 7+
- **Authentication**: JWT
- **ORM**: SQLC (Type-safe SQL)
- **Media Storage**: Local filesystem / Cloud (configurable)

### Kiến trúc ứng dụng
```
┌─────────────────────────────────────────┐
│          Gin HTTP Server (9001)         │
└─────────────────┬───────────────────────┘
                  │
    ┌─────────────┼─────────────┐
    │             │             │
┌───▼───┐   ┌─────▼─────┐   ┌──▼──────┐
│Controllers│  │Middleware│  │  Router │
└───┬───┘   └───────────┘   └─────────┘
    │
┌───▼────────────────────────────────────┐
│         Services (Business Logic)       │
│  - Product Management                   │
│  - Category Management                  │
│  - Media Management                     │
│  - SKU & Option Management              │
└───┬────────────────────────────────────┘
    │
┌───▼────────────────────────────────────┐
│         Repository (Data Access)        │
│  - SQLC Generated Code                  │
│  - MySQL Queries                        │
│  - Transaction Management               │
└───┬────────────────────────────────────┘
    │
┌───▼────────────────────────────────────┐
│           External Services             │
│  - MySQL Database                       │
│  - Redis Cache                          │
│  - Firebase (Optional)                  │
└─────────────────────────────────────────┘
```

### Cấu trúc thư mục
```
ecom_product_service/
├── assets/              # Utilities & helpers
│   ├── api/            # API helpers
│   ├── config/         # Config loader (Viper)
│   ├── token/          # JWT handler
│   ├── util/           # Common utilities
│   └── fire-base/      # Firebase integration
├── controllers/         # HTTP handlers
│   ├── models/         # Request/Response models
│   └── assets/         # Controller helpers
├── services/           # Business logic
│   ├── entity/         # DTOs & Models
│   ├── interface/      # Service interfaces
│   └── assets/         # Service helpers
├── db/
│   ├── migration/      # SQL migrations
│   ├── mysql/          # MySQL client & store
│   ├── query/          # SQLC queries
│   ├── sqlc/           # Generated code
│   └── redis/          # Redis client
├── server/             # Server configuration
├── test/               # Test suites
│   ├── create_product_test.go
│   └── update_product_test.go
├── images/             # Local media storage
├── Dockerfile          # Docker build config
├── docker-run.sh       # Docker deployment script
└── main.go             # Entry point
```


## 🚀 Chức năng chính

### 1. Quản lý sản phẩm (Product)
- ✅ Tạo sản phẩm với nhiều biến thể (SKU)
- ✅ Cập nhật thông tin sản phẩm
- ✅ Xóa mềm sản phẩm (soft delete)
- ✅ Lấy danh sách sản phẩm (phân trang, lọc, sắp xếp)
- ✅ Xem chi tiết sản phẩm (theo ID hoặc Key)
- ✅ Tìm kiếm sản phẩm theo từ khóa
- ✅ Lọc theo danh mục, thương hiệu, khoảng giá

### 2. Quản lý SKU (Product Variants)
- ✅ Tự động tạo SKU name từ option values
- ✅ Quản lý tồn kho (quantity, quantity_reserver)
- ✅ Cập nhật số lượng SKU (HOLD/COMMIT/ROLLBACK)
- ✅ Liên kết SKU với option values
- ✅ Quản lý giá, trọng lượng từng SKU

### 3. Quản lý danh mục (Category)
- ✅ CRUD danh mục sản phẩm
- ✅ Cấu trúc phân cấp (parent-child)
- ✅ Upload ảnh danh mục
- ✅ Lấy danh mục con theo parent ID

### 4. Quản lý thuộc tính (Option Values)
- ✅ Tạo option values (Màu sắc, Size, ...)
- ✅ Upload ảnh cho từng option value
- ✅ Liên kết option values với SKU
- ✅ Tự động tạo SKU name từ options

### 5. Quản lý Media
- ✅ Upload ảnh/video sản phẩm
- ✅ Upload ảnh danh mục
- ✅ Upload ảnh option values
- ✅ Serve media files (local hoặc cloud URL)
- ✅ Xóa media khi cập nhật/xóa sản phẩm

### 6. Tính năng nâng cao
- ✅ Redis caching
- ✅ Transaction rollback khi có lỗi
- ✅ Logging chi tiết (tiếng Việt)
- ✅ JWT authentication
- ✅ CORS configuration
- ✅ Firebase integration (optional)

## � Dependencies

### Required Services
- **MySQL 8.0+**: Database chính
- **Redis 7+**: Cache & session

### Optional Services
- **Firebase**: Cloud storage & authentication (optional)

## 🛠️ Cài đặt & Chạy

### 1. Clone repository
```bash
git clone https://github.com/TranVinhHien/ecom_product_service.git
cd ecom_product_service
```

### 2. Cấu hình môi trường
```bash
cp app.env.example app.env
# Chỉnh sửa app.env với cấu hình của bạn
```

**Biến môi trường cần thiết:**
```bash
DB_SOURCE=root:101204@tcp(localhost:3306)/ecommerce_product_db?parseTime=true
REDIS_ADDRESS=localhost:6379
HTTP_SERVER_ADDRESS=0.0.0.0:9001
JWT_SECRET=your-secret-key
CLIENT_IP=http://localhost:9999,http://localhost:8989
IMAGE_PATH=./images/
```

### 3. Cài đặt dependencies
```bash
go mod tidy
```

### 4. Chạy migrations
```bash
make migrate-up
```

### 5. Generate SQLC code
```bash
make sqlc
```

### 6. Chạy ứng dụng

#### Development
```bash
go run main.go
```

#### Production
```bash
go build -o main main.go
./main
```

Service sẽ chạy tại: **http://localhost:9001**

## 🐳 Docker Deployment

### Quick Start
```bash
# 1. Sửa file cấu hình
nano .env.docker

# 2. Thay đổi IP (thay 172.26.127.95 bằng IP máy của bạn)
DB_SOURCE=root:101204@tcp(<YOUR_IP>:3306)/ecommerce_product_db?parseTime=true
REDIS_ADDRESS=<YOUR_IP>:6379

# 3. Cấp quyền và chạy
chmod +x docker-run.sh
./docker-run.sh
```

### Requirements
- Docker 20.0+
- MySQL container đang chạy trên host
- Redis container đang chạy trên host

### Cách lấy IP máy local
```bash
# Linux/Mac
ip addr show | grep inet

# Hoặc
hostname -I

# Windows (PowerShell)
ipconfig
```

### Useful Commands
```bash
# Xem logs
docker logs -f ecom-product-container

# Restart service
docker restart ecom-product-container

# Stop service
docker stop ecom-product-container

# Remove container (không xóa image)
docker stop ecom-product-container && docker rm ecom-product-container

# Cập nhật biến môi trường (không build lại image)
# 1. Sửa .env.docker
nano .env.docker

# 2. Chạy lại script
./docker-run.sh

# Xem resource usage
docker stats ecom-product-container
```

### Docker Image Info
- **Base image**: golang:1.24-alpine (build stage)
- **Runtime image**: alpine:3.20
- **Size**: ~50-80MB (optimized với multi-stage build)
- **User**: non-root user (bảo mật)

## 🔧 Configuration

### Environment Variables
```bash
# Database MySQL
DB_SOURCE=root:101204@tcp(172.26.127.95:3306)/ecommerce_product_db?parseTime=true

# Redis Cache
REDIS_ADDRESS=172.26.127.95:6379

# HTTP Server (nên giữ nguyên 0.0.0.0:9001)
HTTP_SERVER_ADDRESS=0.0.0.0:9001

# JWT Secret Key
JWT_SECRET=bv-T"-u6@-WR?SHiHQ7yQ]CK*dd9(@jM9BI)|g;zq)ur-Z.Jw/u5HyJHgg,KS.fa

# CORS Origins (phân cách bằng dấu phẩy)
CLIENT_IP=http://localhost:9999,http://localhost:8989

# Image Storage Path
IMAGE_PATH=./images/
```

### Makefile Commands
```bash
make sqlc          # Generate SQLC code
make migrate-up    # Run migrations
make migrate-down  # Rollback migrations
make test          # Run tests
make build         # Build binary
```


## � API Documentation

Base URL: `http://172.26.127.95:9001/v1`

### Categories API

#### GET `/categories/get` - Lấy danh sách danh mục
```bash
# Lấy tất cả danh mục
curl http://172.26.127.95:9001/v1/categories/get

# Lấy danh mục con theo parent ID
curl "http://172.26.127.95:9001/v1/categories/get?cate_id=<parent_id>"
```

#### POST `/categories/create` - Tạo danh mục (ADMIN)
```bash
curl -X POST http://172.26.127.95:9001/v1/categories/create \
  -H "Authorization: Bearer <token>" \
  -F "name=Áo thời trang" \
  -F "parent=<parent_id>" \
  -F "media=@/path/to/image.jpg"
```

#### PUT `/categories/update` - Cập nhật danh mục (ADMIN)
```bash
curl -X PUT http://172.26.127.95:9001/v1/categories/update \
  -H "Authorization: Bearer <token>" \
  -F "cate_id=<category_id>" \
  -F "name=Quần dài nam" \
  -F "media=@/path/to/image.jpg"
```

**Lưu ý:** Chỉ truyền các trường cần cập nhật kèm theo `cate_id`

#### DELETE `/categories/delete/:id` - Xóa danh mục (ADMIN)
```bash
curl -X DELETE "http://172.26.127.95:9001/v1/categories/delete/<id>" \
  -H "Authorization: Bearer <token>"
```

---

### Products API

#### GET `/product/getall` - Lấy danh sách sản phẩm

**Query Parameters:**

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `page` | int | Trang hiện tại (default: 1) | `1` |
| `page_size` | int | Số sản phẩm/trang (default: 10) | `20` |
| `sort` | string | Sắp xếp | `price_asc`, `price_desc`, `name_asc`, `name_desc` |
| `category_path` | string | Lọc theo path danh mục | `/fashion/women/tops` |
| `brand_code` | string | Lọc theo mã thương hiệu | `b001` |
| `shop_id` | string | Lọc theo shop | `shop001` |
| `min_price` | float | Giá tối thiểu | `50000` |
| `max_price` | float | Giá tối đa | `500000` |
| `keywords` | string | Từ khóa tìm kiếm | `áo thun` |

```bash
curl "http://172.26.127.95:9001/v1/product/getall?page=1&page_size=20&sort=price_asc&min_price=100000&max_price=500000"
```

#### GET `/product/detail` - Lấy chi tiết sản phẩm

```bash
# Theo ID
curl "http://172.26.127.95:9001/v1/product/detail?id=<product_id>"

# Theo Key (slug)
curl "http://172.26.127.95:9001/v1/product/detail?key=ao-thun-nam"
```

**Response bao gồm:**
- Thông tin sản phẩm (SPU)
- Danh sách SKU variants
- Option values và ảnh
- Thông tin thương hiệu
- Thông tin danh mục

#### POST `/product/create` - Tạo sản phẩm

**⚠️ Quan trọng:** Chỉ sử dụng `multipart/form-data`

**Form-data fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `product` | JSON string | ✅ | Thông tin sản phẩm (xem cấu trúc bên dưới) |
| `image` | File | ✅ | Ảnh chính sản phẩm |
| `media` | File[] | ❌ | Danh sách ảnh/video phụ (nhiều files) |
| `option_value_images[0..n]` | File | ❌ | Ảnh cho từng option value (theo thứ tự) |

**Product JSON Structure:**

```json
{
  "name": "Áo thun nam cao cấp",
  "key": "ao-thun-nam-cao-cap",
  "description": "Áo thun nam chất liệu cotton 100%, thoáng mát, co giãn tốt...",
  "short_description": "Áo thun nam cao cấp, cotton 100%",
  "brand_id": "b001",
  "category_id": "cat001",
  "shop_id": "shop001",
  "product_is_permission_return": true,
  "product_is_permission_check": true,
  "option_value": [
    {"option_name": "Màu Sắc", "value": "Đỏ"},
    {"option_name": "Màu Sắc", "value": "Xanh"},
    {"option_name": "Size", "value": "M"},
    {"option_name": "Size", "value": "L"}
  ],
  "product_sku": [
    {
      "sku_code": "SKU-DO-M",
      "price": 199000,
      "quantity": 100,
      "weight": 0.3,
      "option_value": [
        {"option_name": "Màu Sắc", "value": "Đỏ"},
        {"option_name": "Size", "value": "M"}
      ]
    },
    {
      "sku_code": "SKU-DO-L",
      "price": 199000,
      "quantity": 50,
      "weight": 0.35,
      "option_value": [
        {"option_name": "Màu Sắc", "value": "Đỏ"},
        {"option_name": "Size", "value": "L"}
      ]
    },
    {
      "sku_code": "SKU-XANH-M",
      "price": 199000,
      "quantity": 80,
      "weight": 0.3,
      "option_value": [
        {"option_name": "Màu Sắc", "value": "Xanh"},
        {"option_name": "Size", "value": "M"}
      ]
    },
    {
      "sku_code": "SKU-XANH-L",
      "price": 199000,
      "quantity": 60,
      "weight": 0.35,
      "option_value": [
        {"option_name": "Màu Sắc", "value": "Xanh"},
        {"option_name": "Size", "value": "L"}
      ]
    }
  ]
}
```

**📌 Quy tắc quan trọng:**

1. **Option Value Images**: 
   - Sắp xếp theo đúng thứ tự với `option_value` trong JSON
   - `option_value_images[0]` → Ảnh cho option_value **đầu tiên** (Đỏ)
   - `option_value_images[1]` → Ảnh cho option_value **thứ hai** (Xanh)
   - Ví dụ: 2 màu có ảnh → truyền 2 files theo thứ tự

2. **Product SKU - Tổ hợp đầy đủ**: 
   - **Phải tạo tất cả** tổ hợp option values
   - Ví dụ: 2 màu × 2 size = **4 SKUs**:
     ```
     Đỏ  + M  → SKU-DO-M
     Đỏ  + L  → SKU-DO-L
     Xanh + M  → SKU-XANH-M
     Xanh + L  → SKU-XANH-L
     ```

3. **SKU Name tự động**:
   - Hệ thống tự động tạo `sku_name` từ option values
   - Format: `Màu Sắc: Đỏ, Size: M`

**Example Request:**

```bash
curl -X POST http://172.26.127.95:9001/v1/product/create \
  -H "Authorization: Bearer <token>" \
  -F "product=$(cat product.json)" \
  -F "image=@main_image.jpg" \
  -F "media=@gallery_1.jpg" \
  -F "media=@gallery_2.jpg" \
  -F "media=@gallery_3.jpg" \
  -F "option_value_images[0]=@red_color.jpg" \
  -F "option_value_images[1]=@blue_color.jpg"
```

#### PUT `/product/update` - Cập nhật sản phẩm

**Updatable Fields:**

```go
type ProductUpdate struct {
    Name                      *string       `json:"name"`                          // Optional
    Key                       *string       `json:"key"`                           // Optional
    Description               *string       `json:"description"`                   // Optional
    ShortDescription          *string       `json:"short_description"`             // Optional
    ProductIsPermissionReturn *bool         `json:"product_is_permission_return"`  // Optional
    ProductIsPermissionCheck  *bool         `json:"product_is_permission_check"`   // Optional
    DeleteStatus              *bool         `json:"delete_status"`                 // Optional (Xóa mềm)
    ProductSKU                []ProductSku  `json:"product_sku"`                   // Optional
    OptionValue               []OptionValue `json:"option_value"`                  // Optional
    KeepMediaURLs             []string      `json:"keep_media_urls"`               // Giữ lại media URLs
    RemoveMediaURLs           []string      `json:"remove_media_urls"`             // Xóa media URLs
    RemoveMainImage           *bool         `json:"remove_main_image"`             // Xóa ảnh chính
}
```

**⚠️ Hạn chế:**
- ❌ **KHÔNG được tạo thêm** `option_value` mới
- ❌ **KHÔNG được sửa** field `option_name`
- ✅ **CHỈ được sửa** field `value` của `option_value` đã tồn tại
- ✅ **SKU**: Phải có `id`, có thể sửa `sku_code`, `price`, `quantity`, `weight`
- ✅ Để **xóa sản phẩm**: set `delete_status: true`

**Example:**

```bash
curl -X PUT http://172.26.127.95:9001/v1/product/update \
  -H "Authorization: Bearer <token>" \
  -F "product_id=<product_id>" \
  -F 'product={"name":"Áo thun nam NEW","price":250000}' \
  -F "image=@new_main_image.jpg"
```

---

### Media API

#### GET `/media/:filename` - Lấy ảnh/video

Service trả về URL có **2 dạng**:

**1. Cloud URL** (nếu đã upload cloud)
```
https://cdn.example.com/images/product.jpg
→ Dùng trực tiếp: <img src="https://..." />
```

**2. Local filename**
```
anhthe.png-7ff26be0-87d1-4400-bc31-e5121a4289ad.png
→ Cần gắn base URL
```

**Cách sử dụng trong Frontend:**

```javascript
// Function để format image URL
function getImageUrl(imageUrl) {
  if (imageUrl.startsWith('http://') || imageUrl.startsWith('https://')) {
    return imageUrl; // Cloud URL
  }
  return `http://172.26.127.95:9001/v1/media/${imageUrl}`; // Local URL
}

// Usage
<img src={getImageUrl(product.image)} />
```

**Example HTML:**

```html
<!-- Cloud URL -->
<img src="https://cdn.example.com/images/product.jpg" alt="Product" />

<!-- Local URL -->
<img src="http://172.26.127.95:9001/v1/media/QR.jpg-6385b218-8136-43c0-8b84-038d9f492d94.jpg" alt="Product" />
```

## 🧪 Testing

Service có bộ test đầy đủ cho các chức năng chính.

### Run all tests
```bash
cd test

# Chạy create product tests
./run_tests.sh

# Chạy update product tests
./run_update_tests.sh
```

### Run specific test
```bash
# Test tạo sản phẩm
go test -v -run TestCreateProduct

# Test cập nhật sản phẩm
go test -v -run TestUpdateProduct

# Test với option images
go test -v -run TestCreateProductWithOptionImages
```

### Test Structure
```
test/
├── create_product_test.go       # Test cases cho create
├── update_product_test.go       # Test cases cho update
├── run_tests.sh                 # Script chạy create tests
├── run_update_tests.sh          # Script chạy update tests
├── README_TEST.md               # Hướng dẫn test create
├── UPDATE_TEST_GUIDE.md         # Hướng dẫn test update
└── QUICKSTART_UPDATE.md         # Quick start guide
```

### Test Coverage
- ✅ Tạo sản phẩm cơ bản
- ✅ Tạo sản phẩm với media files
- ✅ Tạo sản phẩm với option images
- ✅ Cập nhật thông tin sản phẩm
- ✅ Cập nhật SKU prices & quantities
- ✅ Cập nhật option values
- ✅ Upload và xóa media

Chi tiết: [`test/README_TEST.md`](test/README_TEST.md)

## 📊 Database Schema

### Main Tables

**1. product** - Sản phẩm (SPU)
```sql
- id (PK)
- name
- key (unique slug)
- description
- short_description
- brand_id (FK)
- category_id (FK)
- shop_id
- image (main image URL)
- media (JSON array of URLs)
- product_is_permission_return
- product_is_permission_check
- delete_status
- create_date, update_date
```

**2. product_sku** - Biến thể sản phẩm (SKU)
```sql
- id (PK)
- product_id (FK)
- sku_code
- price
- quantity
- quantity_reserver (số lượng đã đặt)
- sku_name (auto-generated)
- weight
- create_date, update_date
```

**3. option_value** - Thuộc tính sản phẩm
```sql
- id (PK)
- product_id (FK)
- option_name (Màu Sắc, Size, ...)
- value (Đỏ, M, ...)
- image (option image URL)
```

**4. sku_attr** - Liên kết SKU với Options
```sql
- sku_id (FK)
- product_id (FK)
- option_value_id (FK)
```

**5. category** - Danh mục
```sql
- id (PK)
- name
- parent_id (self FK)
- path (hierarchical path)
- media
- delete_status
```

**6. brand** - Thương hiệu
```sql
- id (PK)
- code (unique)
- name
```

### Database Triggers

Service sử dụng MySQL trigger để tự động tạo `sku_name`:

```sql
-- Trigger tự động tạo sku_name khi insert/update sku_attr
CREATE TRIGGER generate_sku_name_after_insert ...
CREATE TRIGGER generate_sku_name_after_update ...
```

Chi tiết migrations: [`db/migration/`](db/migration/)

## 🔐 Authentication

Service sử dụng JWT cho authentication:

### JWT Token Format
```
Authorization: Bearer <jwt_token>
```

### Token Claims
```go
{
  "user_id": "uuid",
  "username": "string",
  "role": "user|admin",
  "exp": timestamp
}
```

### Protected Endpoints
- `POST /categories/create` - Admin only
- `PUT /categories/update` - Admin only
- `DELETE /categories/delete/:id` - Admin only
- `POST /product/create` - Authenticated users
- `PUT /product/update` - Authenticated users

### Public Endpoints
- `GET /categories/get`
- `GET /product/getall`
- `GET /product/detail`
- `GET /media/:filename`

## 🐛 Troubleshooting

### 1. Container không kết nối được MySQL/Redis

**Triệu chứng:** Container start nhưng không connect được database

**Giải pháp:**
```bash
# 1. Kiểm tra IP máy local
ip addr show | grep inet
# hoặc
hostname -I

# 2. Sửa file .env.docker với IP đúng
nano .env.docker
# Thay đổi:
DB_SOURCE=root:101204@tcp(<YOUR_IP>:3306)/ecommerce_product_db?parseTime=true
REDIS_ADDRESS=<YOUR_IP>:6379

# 3. Chạy lại container
./docker-run.sh
```

### 2. Permission denied khi chạy script

```bash
chmod +x docker-run.sh
./docker-run.sh
```

### 3. Image build failed - Go version mismatch

**Lỗi:** `go.mod requires go >= 1.24.6 (running go 1.23.x)`

**Giải pháp:** Dockerfile đã được cấu hình với `golang:1.24-alpine`, nếu vẫn lỗi:
```bash
# Xóa image cũ
docker rmi ecom-product-service:latest

# Build lại
./docker-run.sh
```

### 4. Port 9001 already in use

```bash
# Tìm process đang dùng port
lsof -i :9001

# Kill process
kill -9 <PID>

# Hoặc đổi port trong .env.docker
HTTP_SERVER_ADDRESS=0.0.0.0:9002
```

### 5. Xem logs để debug

```bash
# Real-time logs
docker logs -f ecom-product-container

# Logs 100 dòng cuối
docker logs --tail 100 ecom-product-container

# Logs với timestamp
docker logs -t ecom-product-container
```

### 6. Database migration issues

```bash
# Kiểm tra migrations đã chạy
mysql -u root -p ecommerce_product_db -e "SHOW TABLES;"

# Rollback migrations
make migrate-down

# Chạy lại migrations
make migrate-up
```

### 7. Cập nhật biến môi trường mà không build lại

```bash
# 1. Sửa .env.docker
nano .env.docker

# 2. Stop và remove container (không xóa image)
docker stop ecom-product-container
docker rm ecom-product-container

# 3. Chạy lại với env mới
docker run -d \
    --name ecom-product-container \
    --env-file .env.docker \
    -p 9001:9001 \
    -v $(pwd)/images:/app/images \
    --restart unless-stopped \
    ecom-product-service:latest
```

## � Monitoring & Performance

### Resource Usage

Kiểm tra tài nguyên container:
```bash
# Real-time stats
docker stats ecom-product-container

# Chi tiết
docker inspect ecom-product-container
```

**Typical Resource Usage:**
- **CPU**: 0.5-2% (idle), 5-15% (under load)
- **Memory**: 50-150MB (Go app rất nhẹ!)
- **Disk I/O**: Low (chủ yếu read)
- **Network**: Tùy traffic

### Performance Tips

1. **Redis Caching**
   - Enable Redis để cache product queries
   - TTL: 5-15 phút cho product list
   - Invalidate cache khi update product

2. **Database Indexing**
   - Index trên `category_id`, `brand_id`, `shop_id`
   - Index trên `key` (unique)
   - Composite index cho filtering

3. **Image Optimization**
   - Resize images trước khi upload
   - Sử dụng CDN cho production
   - Lazy loading cho frontend

4. **Query Optimization**
   - Sử dụng pagination
   - Limit số lượng joins
   - Use SQLC's compile-time query validation

## 🎯 Best Practices

### Khi Tạo Sản Phẩm

1. ✅ Validate dữ liệu trước khi gửi
2. ✅ Tạo đủ số lượng SKU (tổ hợp tất cả options)
3. ✅ Sắp xếp `option_value_images` đúng thứ tự
4. ✅ Tối ưu kích thước ảnh (< 2MB)
5. ✅ Sử dụng `key` (slug) duy nhất và SEO-friendly
6. ✅ Test với Postman trước khi integrate

### Khi Cập Nhật Sản Phẩm

1. ✅ Chỉ truyền các trường cần thay đổi
2. ❌ Không tạo thêm `option_value` mới
3. ❌ Không sửa `option_name`
4. ✅ Dùng `delete_status: true` để xóa mềm
5. ✅ Backup dữ liệu quan trọng trước khi update
6. ✅ Test trên staging environment trước

### Docker Deployment

1. ✅ Luôn sửa IP trong `.env.docker` trước khi deploy
2. ✅ Backup database trước khi update service
3. ✅ Monitor logs sau khi deploy
4. ✅ Test API endpoints sau khi container chạy
5. ✅ Sử dụng volume cho `/app/images`
6. ✅ Set `--restart unless-stopped` cho auto-restart

### Code Development

1. ✅ Sử dụng SQLC để generate type-safe queries
2. ✅ Luôn wrap DB operations trong transaction
3. ✅ Log chi tiết với context (đã có Vietnamese logging)
4. ✅ Handle errors properly và return ServiceError
5. ✅ Write tests cho business logic quan trọng
6. ✅ Document API changes trong README

## 📝 Changelog

### v1.0.0 (Current)
- ✅ CRUD Products với SKU variants
- ✅ CRUD Categories (hierarchical)
- ✅ Upload & serve media files
- ✅ Auto-generate SKU names từ options
- ✅ Advanced filtering & search
- ✅ Docker support với multi-stage build
- ✅ Redis caching integration
- ✅ JWT authentication
- ✅ Comprehensive test suite
- ✅ Vietnamese logging
- ✅ Transaction rollback on errors

### Planned Features (v1.1.0)
- 🔄 Elasticsearch integration cho full-text search
- 🔄 Cloud storage (S3/GCS) cho media
- 🔄 GraphQL API support
- 🔄 Product reviews & ratings
- 🔄 Inventory alerts
- 🔄 Bulk import/export

## 🤝 Contributing

Chúng tôi hoan nghênh mọi đóng góp!

### How to Contribute

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Coding Guidelines

- Follow Go best practices
- Write tests for new features
- Update documentation
- Use meaningful commit messages
- Keep functions small and focused

## 📞 Support

Nếu gặp vấn đề hoặc có câu hỏi:

- **GitHub Issues**: [Create an issue](https://github.com/TranVinhHien/ecom_product_service/issues)
- **Pull Requests**: [Submit a PR](https://github.com/TranVinhHien/ecom_product_service/pulls)
- **Documentation**: Xem trong các thư mục `test/` và `db/migration/`

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **TranVinhHien** - *Initial work* - [GitHub](https://github.com/TranVinhHien)

---

**Made with ❤️ for E-Commerce Platform**

*Last updated: November 2025*