# 💳 E-Commerce Payment Service

Microservice quản lý thanh toán cho hệ thống thương mại điện tử, xử lý toàn bộ quy trình thanh toán từ khởi tạo, xác thực đến hoàn tất giao dịch.

## 📋 Tổng quan

**Payment Service** là một trong những microservice cốt lõi của hệ thống E-Commerce, chịu trách nhiệm:
- Khởi tạo và xử lý thanh toán online/offline
- Tích hợp cổng thanh toán MoMo
- Quản lý giao dịch (transactions) và ledger
- Xử lý callback từ payment gateway
- Gửi email xác nhận thanh toán
- Phát sự kiện thanh toán qua Kafka
- Theo dõi trạng thái giao dịch

## 🏗️ Kiến trúc

### Tech Stack
- **Language**: Go 1.23
- **Framework**: Gin (HTTP Router)
- **Database**: MySQL 8.0+
- **Cache**: Redis 7+
- **Message Broker**: Kafka
- **Payment Gateway**: MoMo
- **Email Service**: Brevo API
- **Authentication**: JWT
- **ORM**: SQLC (Type-safe SQL)

### Kiến trúc ứng dụng
```
┌─────────────────────────────────────────┐
│          Gin HTTP Server (9003)         │
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
│  - Payment Processing                   │
│  - Transaction Management               │
│  - Ledger Management                    │
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
│  - Kafka Producer/Consumer              │
│  - MoMo Payment Gateway                 │
│  - Brevo Email Service                  │
│  - Order Service (HTTP)                 │
│  - Product Service (HTTP)               │
└─────────────────────────────────────────┘
```

### Cấu trúc thư mục
```
ecom_payment_service/
├── assets/              # Utilities & helpers
│   ├── config/         # Config loader (Viper)
│   ├── email/          # Email templates & sender
│   ├── jobs/           # Background jobs scheduler
│   └── token/          # JWT handler
├── controllers/         # HTTP handlers
├── services/           # Business logic
│   ├── entity/         # DTOs & Models
│   └── assets/         # Service utilities
├── db/
│   ├── migration/      # SQL migrations
│   ├── mysql/          # MySQL client & store
│   ├── query/          # SQLC queries
│   ├── sqlc/           # Generated code
│   └── redis/          # Redis client
├── kafka/              # Kafka producer/consumer
│   ├── kafka.go        # Main Kafka client
│   ├── producer.go     # Message producer
│   ├── consumer.go     # Message consumer
│   ├── events.go       # Event publisher
│   └── topics.go       # Topic definitions
├── server/             # External service clients
└── main.go             # Entry point
```

## 🚀 Chức năng chính

### 1. Quản lý thanh toán
- ✅ Khởi tạo thanh toán (online/COD)
- ✅ Tích hợp MoMo payment gateway
- ✅ Xử lý callback từ MoMo
- ✅ Lấy URL thanh toán lại
- ✅ Kiểm tra trạng thái thanh toán
- ✅ Xử lý thanh toán thất bại

### 2. Quản lý giao dịch (Transactions)
- ✅ Tạo transaction mới
- ✅ Cập nhật trạng thái transaction
- ✅ Lưu thông tin gateway transaction
- ✅ Theo dõi lịch sử giao dịch
- ✅ Xử lý pending/success/failed states

### 3. Quản lý Ledger
- ✅ Tạo ledger entries (DEBIT/CREDIT)
- ✅ Cập nhật balance & pending balance
- ✅ Theo dõi dòng tiền (cash flow)
- ✅ Đối chiếu giao dịch
- ✅ Platform ledger management

### 4. Kafka Event System
- ✅ Publish \`payment.completed\` events
- ✅ Publish \`payment.failed\` events
- ✅ Publish \`transaction.created\` events
- ✅ Consumer để lắng nghe events từ services khác
- ✅ Worker pool để xử lý concurrent messages

### 5. Email Notifications
- ✅ Gửi email xác nhận thanh toán thành công
- ✅ Template HTML responsive
- ✅ Tích hợp Brevo API
- ✅ Xử lý retry khi gửi thất bại

### 6. Payment Methods
- ✅ Quản lý danh sách phương thức thanh toán
- ✅ Chi tiết phương thức thanh toán
- ✅ Cấu hình payment gateway

### 7. Trạng thái giao dịch
```
PENDING → AWAITING_PAYMENT → SUCCESS
              ↓
            FAILED
```

## 📦 Dependencies

### Required Services
- **MySQL 8.0+**: Database chính (ecommerce_transacion_db)
- **Redis 7+**: Cache transaction data
- **Kafka**: Message broker cho events
- **MoMo**: Payment gateway

### Optional Services
- **Order Service** (port 9002): Lấy thông tin đơn hàng
- **Product Service** (port 9001): Lấy thông tin sản phẩm
- **Brevo API**: Email service

## 🛠️ Cài đặt & Chạy

### 1. Clone repository
```bash
git clone https://github.com/TranVinhHien/ecom_payment_service.git
cd ecom_payment_service
```

### 2. Cấu hình môi trường
```bash
cp app.env.example app.env
# Chỉnh sửa app.env với cấu hình của bạn
```

Cấu hình quan trọng:
```bash
# MoMo Gateway
ACCESS_KEY_MOMO=your-access-key
SECRET_KEY_MOMO=your-secret-key
ENDPOINT_MOMO=https://test-payment.momo.vn/v2/gateway/api/create
IPNURL=https://your-domain.ngrok-free.app/v1/transaction/callback

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP=ecom-payment-service-group

# Email
BREVO_API_KEY=your-brevo-api-key
SENDER_EMAIL=your-email@example.com
```

### 3. Cài đặt dependencies
```bash
go mod tidy
```

### 4. Cài đặt Kafka client
```bash
go get github.com/IBM/sarama
```

### 5. Chạy migrations
```bash
make createtb
```

### 6. Generate SQLC code
```bash
make sqlc
```

### 7. Chạy ứng dụng

#### Development
```bash
make run
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

# 2. Chỉnh sửa cấu hình (quan trọng: IPNURL)
nano .env.docker

# 3. Tạo network (nếu chưa có)
make create-network

# 4. Build image
make docker-build

# 5. Deploy
make docker-run

# Hoặc dùng script tự động
./docker-run.sh
```

### Requirements
- Docker 20.0+
- Docker network: \`e-commerce-network\`
- MySQL container đang chạy
- Redis container đang chạy
- Kafka container đang chạy

### Useful Commands
```bash
# Xem logs
make docker-logs

# Xem logs realtime
make docker-logs-tail

# Restart service
make docker-restart

# Stop service
make docker-stop

# Rebuild và deploy lại
make docker-rebuild

# Vào shell container
make docker-exec
```

## 🔧 Configuration

### Environment Variables
```bash
# Database
DB_SOURCE=root:101204@tcp(localhost:3306)/ecommerce_transacion_db?parseTime=true

# Server
HTTP_SERVER_ADDRESS=0.0.0.0:9003

# JWT
JWT_SECRET=bv-T"-u6@-WR?SHiHQ7yQ]CK*dd9(@jM9BI)|g;zq)ur-Z.Jw/u5HyJHgg,KS.fa

# Client
CLIENT_IP=http://localhost:9999

# Redis
REDIS_ADDRESS=localhost:6379

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP=ecom-payment-service-group

# MoMo Payment Gateway
ACCESS_KEY_MOMO=F8BBA842ECF85
SECRET_KEY_MOMO=K951B6PE1waDMi640xX08PD3vg6EkVlz
ENDPOINT_MOMO=https://test-payment.momo.vn/v2/gateway/api/create
REDIRECTURL=http://localhost:9999/vi/dat-hang-thanh-cong
IPNURL=https://your-domain.ngrok-free.app/v1/transaction/callback
PUBLIC_ID=https://your-domain.ngrok-free.app/v1

# Email Service (Brevo)
BREVO_API_KEY=xkeysib-your-api-key
SENDER_EMAIL=your-email@gmail.com
SENDER_NAME=lemarchenoble

# Platform
PLATFORM_ID=111111111111111111111111111111111111
ORDER_DURATION=90m

# External Services
URL_PRODUCT_SERVICE=http://172.26.127.95:9001
URL_ORDER_SERVICE=http://172.26.127.95:9002
```

## 📡 API Endpoints

### Payment APIs
```
GET    /v1/payment-method                     # Danh sách phương thức thanh toán
GET    /v1/payment-method/:id                 # Chi tiết phương thức
POST   /v1/payment/init                       # Khởi tạo thanh toán
GET    /v1/payment/get-url-again              # Lấy URL thanh toán lại
```

### Transaction APIs
```
POST   /v1/transaction/callback               # MoMo callback (IPN)
```

### Admin APIs
```
GET    /v1/admin/transactions                 # Danh sách giao dịch
GET    /v1/admin/ledger                       # Ledger entries
GET    /v1/admin/statistics                   # Thống kê thanh toán
```

## 🔄 Kafka Integration

### Topics
```go
// Published Events
payment.completed           // Khi thanh toán thành công
payment.failed              // Khi thanh toán thất bại
transaction.created         // Khi tạo transaction mới
transaction.timeout         // Khi transaction timeout
order.payment.received      // Thông báo đến Order Service

// Consumed Events (nếu có)
order.created               // Từ Order Service
order.cancelled             // Từ Order Service
```

### Sử dụng Kafka trong Service
```go
// Gửi event thanh toán thành công
eventData := map[string]interface{}{
    "transaction_id": transactionID,
    "order_id":       orderID,
    "amount":         amount,
    "status":         "SUCCESS",
}

err := s.producer.PaymentCompleted(ctx, orderID, eventData)
```

Chi tiết: [kafka/KAFKA_GUIDE.md](kafka/KAFKA_GUIDE.md)

## 💳 MoMo Integration

### Flow thanh toán MoMo
```
1. User click "Thanh toán MoMo"
2. Service gọi POST /v1/payment/init
3. Service tạo payload và gọi MoMo API
4. MoMo trả về payUrl (QR code/deeplink)
5. User scan QR hoặc mở MoMo app
6. User xác nhận thanh toán trong MoMo
7. MoMo gọi IPN callback → POST /v1/transaction/callback
8. Service xử lý và cập nhật trạng thái
9. Service gửi email xác nhận
10. Service publish Kafka event
11. Service callback đến Order Service
```

### Test MoMo Sandbox
```bash
# Credentials test
ACCESS_KEY: F8BBA842ECF85
SECRET_KEY: K951B6PE1waDMi640xX08PD3vg6EkVlz

# Endpoint test
https://test-payment.momo.vn/v2/gateway/api/create

# Test QR Code Payment
Amount: Bất kỳ (min 1000 VND)
OTP: Nhập bất kỳ 6 số
```

## 📧 Email Templates

Service tự động gửi email xác nhận khi:
- ✅ Thanh toán thành công
- ✅ Đơn hàng được tạo (COD)
- ⚠️ Thanh toán thất bại (optional)

Template: [assets/email/payment_success.go](assets/email/payment_success.go)

## 🧪 Testing

### Run tests
```bash
make test
```

### Test MoMo callback manually
```bash
curl -X POST http://localhost:9003/v1/transaction/callback \
  -H "Content-Type: application/json" \
  -d '{
    "partnerCode": "MOMO",
    "orderId": "order-123",
    "requestId": "transaction-456",
    "amount": 100000,
    "orderInfo": "Test payment",
    "orderType": "momo_wallet",
    "transId": 123456789,
    "resultCode": 0,
    "message": "Success",
    "payType": "qr",
    "responseTime": 1234567890,
    "extraData": "",
    "signature": "..."
  }'
```

### Test Kafka events
```bash
# Subscribe to topic
kafka-console-consumer --bootstrap-server localhost:9092 \
  --topic payment.completed \
  --from-beginning

# Check consumer group
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --group ecom-payment-service-group --describe
```

## 📊 Database Schema

### Main Tables
- **\`transactions\`**: Giao dịch thanh toán
  - id, transaction_code, order_id, payment_method_id
  - amount, currency, type, status
  - gateway_transaction_id, notes
  - created_at, processed_at

- **\`payment_methods\`**: Phương thức thanh toán
  - id, code, name, type (ONLINE/OFFLINE)
  - description, is_active

- **\`account_ledgers\`**: Sổ cái platform
  - id, account_number, account_name
  - balance, pending_balance
  - currency, status

- **\`ledger_entries\`**: Bút toán kế toán
  - id, ledger_id, transaction_id
  - amount, type (DEBIT/CREDIT)
  - description, created_at

- **\`order_platform_costs\`**: Chi phí platform
  - order_id, payment_transaction_id
  - site_order_voucher_discount_amount
  - site_promotion_discount_amount
  - site_shipping_discount_amount
  - total_site_funded_product_discount

- **\`shop_order_settlements\`**: Đối soát với shop
  - id, shop_order_id, order_transaction_id
  - status, order_subtotal
  - shop_funded_product_discount
  - site_funded_product_discount
  - shop_voucher_discount, shipping_fee
  - commission_fee, net_settled_amount

Chi tiết: [db/migration/](db/migration/)

## 🔐 Authentication

Service sử dụng JWT cho authentication:
- **User Token**: Cho customer APIs (payment init)
- **System Token**: Cho giao tiếp giữa các services
- **No Auth**: Cho MoMo callback (verify bằng signature)

## 📝 Development

### Make commands
```bash
make run                 # Chạy ứng dụng
make sqlc                # Generate SQLC code
make createtb            # Run migrations
make droptb              # Rollback migrations
make docker-build        # Build Docker image
make docker-run          # Run container
make docker-rebuild      # Rebuild & run
make docker-logs         # Xem logs
make create-network      # Tạo Docker network
```

## 🐛 Troubleshooting

### Lỗi thường gặp

**1. Panic: send on closed channel**
```bash
# Nguyên nhân: Worker pool bị đóng khi đang xử lý message
# Giải pháp: Đã fix bằng mutex và check stopped state trong consumer.go
# Code đã có sẵn handle case này
```

**2. MoMo callback không về**
```bash
# Kiểm tra IPNURL có public không (dùng ngrok)
ngrok http 9003

# Cập nhật IPNURL trong .env
IPNURL=https://your-subdomain.ngrok-free.app/v1/transaction/callback

# Kiểm tra signature có đúng không
# Xem logs MoMo: https://developers.momo.vn/
```

**3. Kafka connection refused**
```bash
# Đảm bảo Kafka đang chạy
docker ps | grep kafka

# Kiểm tra port
netstat -tulpn | grep 9092

# Start Kafka nếu chưa chạy
make startkafka
```

**4. Redis cache miss**
```bash
# Transaction data được lưu trong Redis với TTL (ORDER_DURATION)
# Nếu quá lâu mới callback, data có thể bị xóa
# Kiểm tra ORDER_DURATION trong .env (mặc định 90m)

# Check Redis
redis-cli
> KEYS transaction:online:*
> TTL transaction:online:order-123
```

**5. Email không gửi được**
```bash
# Kiểm tra Brevo API key
curl -X GET "https://api.brevo.com/v3/account" \
  -H "api-key: your-api-key"

# Kiểm tra sender email đã verify chưa
# Vào Brevo dashboard: https://app.brevo.com/
```

**6. Database connection error**
```bash
# Kiểm tra MySQL đang chạy
docker ps | grep mysql

# Test connection
mysql -h 172.26.127.95 -u root -p101204 -e "SHOW DATABASES;"

# Kiểm tra database tồn tại
mysql -h 172.26.127.95 -u root -p101204 -e "USE ecommerce_transacion_db; SHOW TABLES;"
```

## 🔍 Monitoring & Logs

### View logs
```bash
# Container logs
make docker-logs

# Realtime logs (100 dòng cuối)
make docker-logs-tail

# Grep specific error
docker logs ecom-payment-container 2>&1 | grep "ERROR"

# Grep MoMo callback
docker logs ecom-payment-container 2>&1 | grep "CallBackMoMo"
```

### Health check
```bash
# Check service health
curl http://localhost:9003/health

# Check Kafka connection
docker exec -it ecom-payment-container sh
# Trong container
ps aux | grep kafka
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (\`git checkout -b feature/AmazingFeature\`)
3. Commit your changes (\`git commit -m 'Add some AmazingFeature'\`)
4. Push to the branch (\`git push origin feature/AmazingFeature\`)
5. Open a Pull Request

### Coding Standards
- Follow Go best practices
- Use SQLC for database queries
- Write tests for new features
- Update documentation
- Use meaningful commit messages

## 📄 License

This project is licensed under the MIT License.

## 👥 Authors

- **TranVinhHien** - [GitHub](https://github.com/TranVinhHien)

## 📞 Support

Nếu có vấn đề, vui lòng tạo issue trên GitHub hoặc liên hệ team.

---

## 🎯 Roadmap

- [ ] Tích hợp thêm payment gateway (VNPay, ZaloPay)
- [ ] Retry mechanism cho failed payments
- [ ] Payment analytics dashboard
- [ ] Webhook cho external systems
- [ ] Unit tests coverage 80%+
- [ ] Load testing với k6
- [ ] OpenTelemetry tracing
- [ ] Circuit breaker cho external calls
- [ ] Dead letter queue cho failed events
- [ ] Payment reconciliation tool

## 📚 Documentation

- [Kafka Integration Guide](kafka/KAFKA_GUIDE.md)
- [API Documentation](docs/API.md) (coming soon)
- [Database Schema](docs/DATABASE.md) (coming soon)
- [Deployment Guide](docs/DEPLOYMENT.md) (coming soon)

## 🏆 Best Practices Implemented

✅ **Clean Architecture**: Separation of concerns (controllers, services, repository)  
✅ **Type Safety**: SQLC for type-safe SQL queries  
✅ **Event-Driven**: Kafka for async communication  
✅ **Caching**: Redis for performance optimization  
✅ **Transaction Management**: ACID compliance with MySQL transactions  
✅ **Error Handling**: Comprehensive error handling & logging  
✅ **Security**: JWT authentication, signature verification  
✅ **Scalability**: Worker pool, Kafka consumer groups  
✅ **Observability**: Structured logging with zerolog  
✅ **Configuration**: Environment-based config with Viper  

---

**Version**: 1.0.0  
**Last Updated**: November 12, 2025  
**Status**: Production Ready ✅
