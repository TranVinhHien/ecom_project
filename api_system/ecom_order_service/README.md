# 🛒 E-Commerce Order Service

Microservice quản lý đơn hàng cho hệ thống thương mại điện tử, xử lý toàn bộ vòng đời đơn hàng từ tạo, thanh toán đến vận chuyển và hoàn thành.

## 📋 Tổng quan

**Order Service** là một trong những microservice cốt lõi của hệ thống E-Commerce, chịu trách nhiệm:
- Tạo và quản lý đơn hàng
- Quản lý voucher và áp dụng giảm giá
- Xử lý thanh toán online/offline
- Theo dõi trạng thái đơn hàng và vận chuyển
- Tích hợp với Product Service và Transaction Service
- Xử lý events từ Kafka (payment success/failed)

## 🏗️ Kiến trúc

### Tech Stack
- **Language**: Go 1.23
- **Framework**: Gin (HTTP Router)
- **Database**: MySQL 8.0+
- **Cache**: Redis 7+
- **Message Broker**: Kafka
- **Authentication**: JWT
- **ORM**: SQLC (Type-safe SQL)

### Kiến trúc ứng dụng
```
┌─────────────────────────────────────────┐
│          Gin HTTP Server (9002)         │
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
│  - Order Management                     │
│  - Voucher Management                   │
│  - Payment Processing                   │
└───┬────────────────────────────────────┘
    │
┌───▼────────────────────────────────────┐
│         Repository (Data Access)        │
│  - SQLC Generated Code                  │
│  - MySQL Queries                        │
└───┬────────────────────────────────────┘
    │
┌───▼────────────────────────────────────┐
│           External Services             │
│  - MySQL Database                       │
│  - Redis Cache                          │
│  - Kafka (Events)                       │
│  - Product Service (gRPC/HTTP)          │
│  - Transaction Service (gRPC/HTTP)      │
└─────────────────────────────────────────┘
```

### Cấu trúc thư mục
```
ecom_order_service/
├── assets/              # Utilities & helpers
│   ├── api/            # API helpers
│   ├── config/         # Config loader (Viper)
│   ├── token/          # JWT handler
│   └── util/           # Common utilities
├── controllers/         # HTTP handlers
├── services/           # Business logic
│   ├── entity/         # DTOs & Models
│   └── interface/      # Service interfaces
├── db/
│   ├── migration/      # SQL migrations
│   ├── mysql/          # MySQL client
│   ├── query/          # SQLC queries
│   ├── sqlc/           # Generated code
│   └── redis/          # Redis client
├── kafka/              # Kafka producer/consumer
├── server/             # External service clients
└── main.go             # Entry point
```

## 🚀 Chức năng chính

### 1. Quản lý đơn hàng
- ✅ Tạo đơn hàng mới (online/offline payment)
- ✅ Lấy danh sách đơn hàng của user
- ✅ Xem chi tiết đơn hàng
- ✅ Tìm kiếm & lọc đơn hàng (theo trạng thái, ngày, giá trị)
- ✅ Cập nhật trạng thái đơn hàng
- ✅ Xử lý vận chuyển

### 2. Quản lý Voucher
- ✅ Tạo và cập nhật voucher
- ✅ Lấy danh sách voucher (public & assigned)
- ✅ Lọc voucher (theo shop, loại, giá trị)
- ✅ Áp dụng voucher khi đặt hàng
- ✅ Kiểm tra điều kiện voucher
- ✅ Rollback voucher khi hủy đơn

### 3. Xử lý thanh toán
- ✅ Thanh toán online (qua Transaction Service)
- ✅ Thanh toán offline (COD)
- ✅ Callback từ payment gateway
- ✅ Xử lý payment events từ Kafka

### 4. Trạng thái đơn hàng
```
PENDING → AWAITING_PAYMENT → AWAITING_CONFIRMATION 
    → PROCESSING → SHIPPED → COMPLETED
```

## 📦 Dependencies

### Required Services
- **MySQL 8.0+**: Database chính
- **Redis 7+**: Cache & session
- **Kafka**: Message broker cho events

### Optional Services
- **Product Service** (port 9001): Lấy thông tin sản phẩm
- **Transaction Service** (port 9003): Xử lý thanh toán

## 🛠️ Cài đặt & Chạy

### 1. Clone repository
```bash
git clone https://github.com/TranVinhHien/ecom_order_service.git
cd ecom_order_service
```

### 2. Cấu hình môi trường
```bash
cp app.env.example app.env
# Chỉnh sửa app.env với cấu hình của bạn
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

## 🐳 Docker Deployment

### Quick Start
```bash
# 1. Tạo file cấu hình
cp .env.docker.example .env.docker

# 2. Chỉnh sửa cấu hình (nếu cần)
nano .env.docker

# 3. Deploy
./docker-run.sh
```

### Requirements
- Docker 20.0+
- Docker network: `e-commerce-network`
- MySQL container đang chạy
- Redis container đang chạy
- Kafka container đang chạy

### Useful Commands
```bash
# Xem logs
docker logs -f ecom-order-container

# Restart service
docker restart ecom-order-container

# Stop service
docker stop ecom-order-container

# Remove container
docker stop ecom-order-container && docker rm ecom-order-container
```

Chi tiết: [DOCKER_DEPLOYMENT.md](./DOCKER_DEPLOYMENT.md)

## 🔧 Configuration

### Environment Variables
```bash
# Database
DB_SOURCE=root:12345@tcp(localhost:3306)/ecommerce_order_db?parseTime=true

# Server
HTTP_SERVER_ADDRESS=0.0.0.0:9002

# JWT
JWT_SECRET=your-secret-key

# Redis
REDIS_ADDRESS=localhost:6379

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP=ecom-order-service-group

# External Services
URL_PRODUCT_SERVICE=http://localhost:9001
URL_TRANSACTION_SERVICE=http://localhost:9003

# System Token
TOKEN_SYSTEM=your-system-jwt-token
```

## 📡 API Endpoints

### Customer APIs
```
POST   /api/v1/orders                    # Tạo đơn hàng
GET    /api/v1/orders                    # Danh sách đơn hàng
GET    /api/v1/orders/:orderCode         # Chi tiết đơn hàng
GET    /api/v1/orders/search/detail      # Tìm kiếm đơn hàng
GET    /api/v1/vouchers                  # Danh sách voucher
PUT    /api/v1/orders/callback_payment_online/:order_id  # Payment callback
```

### Admin/Shop APIs
```
POST   /api/v1/vouchers                  # Tạo voucher
PUT    /api/v1/vouchers/:voucherID       # Cập nhật voucher
PUT    /api/v1/orders/admin/update_status  # Cập nhật trạng thái
```

## 🧪 Testing

### Run tests
```bash
make test
```

### Run specific test
```bash
go test ./services/test -v
```

## 📊 Database Schema

### Main Tables
- `orders`: Đơn hàng tổng
- `shop_orders`: Đơn hàng theo shop
- `order_items`: Sản phẩm trong đơn
- `vouchers`: Voucher
- `user_vouchers`: Voucher của user
- `voucher_usage_history`: Lịch sử dùng voucher

Chi tiết: [db/migration/](./db/migration/)

## 🔐 Authentication

Service sử dụng JWT cho authentication:
- **User Token**: Cho customer APIs
- **System Token**: Cho giao tiếp giữa các services

## 📝 Development

### Make commands
```bash
make sqlc          # Generate SQLC code
make migrate-up    # Run migrations
make migrate-down  # Rollback migrations
make test          # Run tests
make build         # Build binary
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👥 Authors

- **TranVinhHien** - [GitHub](https://github.com/TranVinhHien)

## 📞 Support

Nếu có vấn đề, vui lòng tạo issue trên GitHub hoặc liên hệ team.
