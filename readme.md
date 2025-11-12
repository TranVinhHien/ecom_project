# 🚀 Dự án E-Commerce Le Marché Noble

## 🌟 Giới thiệu Tổng quan

Le Marché Noble là một nền tảng thương mại điện tử đa nhà bán hàng (multi-shop) hoàn chỉnh, được xây dựng trên kiến trúc Microservices. Hệ thống không chỉ bao gồm các phân hệ nghiệp vụ cốt lõi (Sản phẩm, Đơn hàng, Thanh toán) mà còn được trang bị một hệ thống Trí tuệ Nhân tạo (AI Agent) tiên tiến để mang lại trải nghiệm mua sắm thông minh và cá nhân hóa.

Hệ thống được thiết kế với sự tách biệt rõ ràng giữa các lớp:

- **Frontend (Client)**: Giao diện người dùng hiện đại, linh hoạt.
- **Backend (Core Services)**: Các microservice nghiệp vụ (Go) xử lý logic cốt lõi.
- **AI Layer (Agent & Data)**: Lớp dịch vụ AI và dữ liệu ngữ nghĩa.

---

## 🏗️ Kiến trúc Hệ thống Tổng thể

Hệ thống được chia thành ba nhóm dịch vụ chính, giao tiếp với nhau qua API Gateway và hệ thống bus sự kiện (Kafka).

### Lớp Giao diện (Client Layer)

**Le Marché Noble (Client)**: Là ứng dụng Next.js 14 (App Router) cung cấp toàn bộ giao diện người dùng, quản lý trạng thái phía client, và tương tác trực tiếp với API Gateway và AI Agent.

### Lớp Nghiệp vụ Cốt lõi (Core Business Layer)

Được xây dựng chủ yếu bằng Golang, tuân thủ Clean Architecture và sử dụng SQLC để truy vấn CSDL.

- **Product Service (Port 9001)**: Quản lý SPU, SKU, Danh mục, Thương hiệu và Media.
- **Order Service (Port 9002)**: Quản lý vòng đời đơn hàng (tổng và chi tiết shop), xử lý Vouchers, và lắng nghe sự kiện thanh toán.
- **Payment & Transaction Service (Port 9003)**: Quản lý dòng tiền, tích hợp cổng thanh toán (MoMo), quản lý hệ thống Ví (Ledger) nội bộ, và phát sự kiện thanh toán qua Kafka.
- **Các Service Hỗ trợ**: Bao gồm Identity, Profile, Address, Cart, Shop, Banner, Policy, và Analytics (Port 9004) để xử lý các nghiệp vụ phụ trợ.

### Lớp Trí tuệ (AI & Data Layer)

- **AI Agent Service (Port 9000)**: Lõi AI của hệ thống, cung cấp các khả năng tương tác thông minh.
- **Semantic Search Service**: Cơ sở dữ liệu vector (Redis-stack) lưu trữ dữ liệu ngữ nghĩa đã được embedding để phục vụ tìm kiếm.

---

## 🧩 Chi tiết các Microservice

### 1. Frontend: Le Marché Noble (Client)

**Tóm tắt** (từ README3.md): Giao diện người dùng (UI) chính của dự án, được xây dựng bằng công nghệ web hiện đại để mang lại trải nghiệm mượt mà, đẹp mắt và tùy biến cao.

**Công nghệ**: Next.js 14, React 18, TypeScript, Tailwind CSS.

**Quản lý trạng thái**: Zustand (cho giỏ hàng/checkout) và Redux Toolkit, TanStack Query (cho server state).

**Tính năng nổi bật**:

- **Hệ thống Theme**: 9 bảng màu, gradient, Light/Dark mode.
- **Đa ngôn ngữ (i18n)**: Tiếng Việt và Tiếng Anh.
- **Xác thực**: JWT tự động refresh token.
- **Tích hợp**: Kết nối trực tiếp với API Gateway và AI Chatbot.

### 2. Backend: Product Service (Port 9001)

**Tóm tắt** (từ README2.md): Quản lý toàn bộ thông tin liên quan đến sản phẩm, danh mục, và biến thể (SKU).

**Công nghệ**: Go, Gin, MySQL, Redis, SQLC.

**Tính năng chính**:

- Quản lý Sản phẩm (SPU) và Biến thể (SKU) chi tiết.
- Quản lý Danh mục (Category) phân cấp (cha-con).
- Quản lý Thuộc tính (Option Values) và tự động tạo sku_name (ví dụ: Màu Sắc: Đỏ, Size: M).
- Upload và quản lý Media (ảnh/video) cho sản phẩm, danh mục.
- Tìm kiếm và lọc sản phẩm nâng cao.

### 3. Backend: Order Service (Port 9002)

**Tóm tắt** (từ README.md): Xử lý toàn bộ vòng đời đơn hàng, từ lúc tạo cho đến khi hoàn thành.

**Công nghệ**: Go, Gin, MySQL, Redis, Kafka, SQLC.

**Tính năng chính**:

- Tạo đơn hàng tổng (orders) và chia thành các đơn hàng shop (shop_orders).
- Quản lý Voucher: Tạo, kiểm tra điều kiện và áp dụng vào đơn hàng.
- Theo dõi trạng thái đơn hàng chi tiết (AWAITING_PAYMENT, PROCESSING, SHIPPED, etc.).
- Lắng nghe sự kiện (Subscribing) từ Kafka (ví dụ: payment.completed) để tự động cập nhật trạng thái đơn hàng.

### 4. Backend: Payment & Transaction Service (Port 9003)

**Tóm tắt** (từ README1.md): Chịu trách nhiệm cho toàn bộ dòng tiền của hệ thống. Đây là dịch vụ duy nhất được phép xử lý các giao dịch tài chính.

**Công nghệ**: Go, Gin, MySQL, Redis, Kafka, SQLC.

**Tích hợp**: MoMo (Cổng thanh toán), Brevo (Gửi email).

**Tính năng chính**:

- Khởi tạo thanh toán (Online MoMo, Offline COD).
- Xử lý callback (IPN) từ MoMo để xác thực thanh toán.
- Hệ thống Ví (Ledger): Quản lý balance (số dư khả dụng) và pending_balance (số dư chờ) cho Sàn và Shop.
- Phát sự kiện (Publishing): Phát các sự kiện tài chính quan trọng (payment.completed, payment.failed) lên Kafka.

### 5. AI: Semantic Search Service (Vector DB)

**Tóm tắt** (Thông tin bổ sung): Đây là dịch vụ nền tảng dữ liệu cho AI.

**Công nghệ**: Redis-stack (sử dụng khả năng lưu trữ Vector).

**Mô hình Embedding**: dangvantuan/vietnamese-document-embedding.

**Chức năng**:

- Lưu trữ vector embedding của dữ liệu ngữ nghĩa (thông tin sản phẩm, mô tả, chính sách...).
- Cung cấp khả năng tìm kiếm vector (vector similarity search) để AI Agent có thể tìm kiếm thông tin theo ngữ nghĩa thay vì từ khóa chính xác.

### 6. AI: AI Agent Service (Port 9102)

**Tóm tắt** (Thông tin bổ sung): Đây là "bộ não" thông minh của hệ thống, tương tác trực tiếp với người dùng và các service khác.

**Công nghệ**: Google ADK (Agent Development Kit), MCP.

**Chức năng chính**:

- **Tương tác khách hàng**: Cung cấp giao diện chat 24/7.
- **Tìm kiếm thông minh**: Tìm sản phẩm dựa trên gợi ý và ngữ nghĩa (sử dụng Semantic Search Service).
- **Tra cứu nghiệp vụ**: Giao tiếp với các service Go để:
  - Xem danh sách đơn hàng của người dùng.
  - Tìm kiếm voucher người dùng hiện có.
  - Tra cứu chính sách của sàn (đổi trả, bảo mật...).
- **Tóm tắt & Gợi ý**: Phân tích chi tiết sản phẩm và bình luận (từ Analytics Service) để đưa ra tóm tắt và gợi ý người dùng có nên mua sản phẩm hay không.