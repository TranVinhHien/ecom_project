# 🛍️ Le Marché Noble - E-Commerce Client

<div align="center">

![Next.js](https://img.shields.io/badge/Next.js-14.2.21-black?style=for-the-badge&logo=next.js)
![React](https://img.shields.io/badge/React-18-blue?style=for-the-badge&logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5-blue?style=for-the-badge&logo=typescript)
![Tailwind CSS](https://img.shields.io/badge/Tailwind-3.4-38bdf8?style=for-the-badge&logo=tailwind-css)

Nền tảng thương mại điện tử hiện đại với giao diện người dùng đẹp mắt và trải nghiệm mua sắm mượt mà.

[Tính Năng](#-tính-năng-chính) • [Công Nghệ](#-công-nghệ-sử-dụng) • [Cài Đặt](#-cài-đặt) • [API](#-kết-nối-api) • [Triển Khai](#-triển-khai)

</div>

---

## 📋 Mục Lục

- [Giới Thiệu](#-giới-thiệu)
- [Tính Năng Chính](#-tính-năng-chính)
- [Công Nghệ Sử Dụng](#-công-nghệ-sử-dụng)
- [Kiến Trúc Hệ Thống](#-kiến-trúc-hệ-thống)
- [Kết Nối API](#-kết-nối-api)
- [Cài Đặt](#-cài-đặt)
- [Chạy Dự Án](#-chạy-dự-án)
- [Build & Deploy](#-build--deploy)
- [Cấu Trúc Thư Mục](#-cấu-trúc-thư-mục)
- [Theme System](#-theme-system)
- [Đóng Góp](#-đóng-góp)

---

## 🌟 Giới Thiệu

**Le Marché Noble** là một nền tảng thương mại điện tử cao cấp được xây dựng với Next.js 14, cung cấp trải nghiệm mua sắm trực tuyến mượt mà và chuyên nghiệp. Dự án tập trung vào:

- ✨ **Giao diện đẹp mắt** với hệ thống theme động (9 bảng màu + gradient)
- 🌍 **Đa ngôn ngữ** (Tiếng Việt & English)
- 🎨 **UX/UI hiện đại** với animations và transitions mượt mà
- 🚀 **Performance tối ưu** với Next.js 14 App Router
- 🔐 **Bảo mật** với JWT Authentication & Auto Token Refresh
- 💬 **AI Chatbot** hỗ trợ khách hàng tự động

---

## 🎯 Tính Năng Chính

### 🛒 **Shopping Experience**
- ✅ Duyệt sản phẩm với lọc & tìm kiếm nâng cao
- ✅ Chi tiết sản phẩm với hình ảnh, options, và reviews
- ✅ Giỏ hàng thông minh với cập nhật real-time
- ✅ Hệ thống voucher & khuyến mãi
- ✅ Thanh toán đơn giản với nhiều phương thức

### 👤 **User Management**
- ✅ Đăng ký / Đăng nhập với JWT
- ✅ Quản lý profile & địa chỉ giao hàng
- ✅ Lịch sử đơn hàng & tracking
- ✅ Wishlist & sản phẩm yêu thích

### 🎨 **Theme & Customization**
- ✅ 9 bảng màu: Orange, Blue, Green, Rose, Zinc, Purple, Cyan, Yellow, Teal
- ✅ Light & Dark mode
- ✅ Gradient backgrounds với 2-3 màu pha trộn
- ✅ Smooth hover effects & animations
- ✅ Responsive design cho mọi thiết bị

### 💬 **AI Support**
- ✅ AI Chatbot hỗ trợ 24/7
- ✅ Tự động trả lời câu hỏi thường gặp
- ✅ Gợi ý sản phẩm thông minh
- ✅ Feedback & rating system

### 📱 **Additional Features**
- ✅ Multi-language support (i18n)
- ✅ Search với filters nâng cao
- ✅ Image optimization với Next.js Image
- ✅ SEO friendly
- ✅ Analytics & tracking

---

## 🛠️ Công Nghệ Sử Dụng

### **Core Framework**
```json
{
  "next": "14.2.21",           // React Framework với App Router
  "react": "18",               // UI Library
  "typescript": "5"            // Type Safety
}
```

### **UI & Styling**
- **Tailwind CSS 3.4** - Utility-first CSS framework
- **Radix UI** - Headless UI components
  - Dialog, Dropdown, Popover, Toast, Tabs, etc.
- **Framer Motion 12** - Animation library
- **Lucide React** - Beautiful icons
- **Class Variance Authority** - Component variants

### **State Management**
- **Zustand 5** - Simple & powerful state management
- **Redux Toolkit 2.5** - Complex state với middleware
- **TanStack Query 5** - Server state & caching

### **Form & Validation**
- **React Hook Form 7** - Form handling
- **Zod 4** - Schema validation
- **@hookform/resolvers** - Form validation integration

### **HTTP & API**
- **Axios 1.7** - HTTP client với interceptors
- **JWT Decode** - Token parsing
- **Cookies Next** - Cookie management

### **Internationalization**
- **Next Intl 3** - i18n for Next.js
- **i18next 24** - Translation framework

### **Theming**
- **Next Themes 0.4** - Dark mode support
- Custom theme system với 9 color palettes

### **Development Tools**
- **ESLint** - Code linting
- **PostCSS** - CSS processing
- **Tailwind Animate** - Animation utilities

---

## 🏗️ Kiến Trúc Hệ Thống

```
┌─────────────────────────────────────────────────┐
│         Next.js 14 Client Application           │
│              (Port: 9999)                       │
└─────────────────────────────────────────────────┘
                      │
                      ├─────────────────────────────┐
                      │                             │
        ┌─────────────▼──────────────┐   ┌─────────▼────────────┐
        │   Gateway Service          │   │   AI Agent Service   │
        │   (lemarchenoble.id.vn)    │   │   (localhost:9000)   │
        │   - Identity               │   │   - Chat Session     │
        │   - Profile                │   │   - AI Responses     │
        │   - Address                │   └──────────────────────┘
        │   - Categories             │
        └────────────┬───────────────┘
                     │
        ┌────────────┼────────────────────────┐
        │            │                        │
┌───────▼──────┐ ┌──▼────────┐ ┌────────────▼───────┐
│ Product      │ │ Order     │ │ Analytics          │
│ Service      │ │ Service   │ │ Service            │
│ (Port 9001)  │ │ (Port     │ │ (Port 9004)        │
│              │ │ 9002)     │ │ - Reviews          │
│ - Products   │ │ - Orders  │ │ - Complaints       │
│ - Details    │ │ - Cart    │ │ - Statistics       │
│              │ │ - Checkout│ │                    │
└──────────────┘ └───────────┘ └────────────────────┘
```

---

## 🔌 Kết Nối API

Ứng dụng kết nối với nhiều microservices:

### **1. Gateway Service** 
**Base URL:** `https://lemarchenoble.id.vn/api/v1`

#### Identity Service
```typescript
POST   /identity/auth/login           // Đăng nhập
POST   /identity/auth/refresh         // Refresh token
POST   /identity/users/register       // Đăng ký
```

#### Profile Service
```typescript
GET    /profile/users/profiles/get-my-profile           // Lấy profile
PUT    /profile/users/profiles/update                   // Cập nhật profile
GET    /profile/users/profiles/profile-subs/get-all-my-sub-profile  // Lấy địa chỉ
POST   /profile/users/profiles/profile-subs/insert      // Thêm địa chỉ
PUT    /profile/users/profiles/profile-subs/update/:id  // Sửa địa chỉ
DELETE /profile/users/profiles/profile-subs/delete/:id  // Xóa địa chỉ
```

#### Address Service
```typescript
GET    /address/provinces/get-all     // Lấy tỉnh/thành
GET    /address/districts/get-all     // Lấy quận/huyện
GET    /address/wards/get-all         // Lấy phường/xã
```

#### Category Service
```typescript
GET    /categories/get                // Lấy danh mục
```

### **2. Product Service**
**Base URL:** `http://172.26.127.95:9001/v1`

```typescript
GET    /product/getall                // Danh sách sản phẩm
GET    /product/getdetail/:id         // Chi tiết sản phẩm
GET    /media/products                // Hình ảnh sản phẩm
```

### **3. Order Service**
**Base URL:** `http://172.26.127.95:9002/v1`

```typescript
GET    /orders                        // Danh sách đơn hàng
POST   /orders                        // Tạo đơn hàng
GET    /orders/:id                    // Chi tiết đơn hàng
PUT    /orders/:id                    // Cập nhật đơn hàng
```

### **4. AI Agent Service**
**Base URL:** `http://localhost:9000/api`

```typescript
POST   /session                       // Tạo chat session
POST   /message                       // Gửi tin nhắn cho AI
```

### **5. Analytics Service**
**Base URL:** `http://172.26.127.95:9004/v1`

```typescript
POST   /public/chatbox/review         // Review chatbot
POST   /public/customer-support/complaint  // Gửi khiếu nại
```

### **Authentication**
Tất cả API được bảo vệ bằng JWT Token:

```typescript
headers: {
  'Authorization': 'Bearer <access_token>',
  'Content-Type': 'application/json'
}
```

**Auto Token Refresh:**
- Tự động refresh khi token hết hạn (401)
- Queue các request thất bại và retry sau khi refresh
- Redirect về login nếu refresh thất bại

---

## 📥 Cài Đặt

### **Yêu Cầu Hệ Thống**
- Node.js >= 18.x
- npm >= 9.x hoặc yarn >= 1.22.x
- Git

### **Bước 1: Clone Repository**
```bash
git clone https://github.com/TranVinhHien/ecom_user_client.git
cd ecom_user_client
```

### **Bước 2: Cài Đặt Dependencies**
```bash
npm install
# hoặc
yarn install
```

### **Bước 3: Cấu Hình Environment**
Tạo file `.env.local` trong thư mục root:

```env
# API Configuration
NEXT_PUBLIC_API_BASE_URL=https://lemarchenoble.id.vn/api/v1
NEXT_PUBLIC_API_AGENT_URL=http://localhost:9000/api
NEXT_PUBLIC_API_PRODUCT_URL=http://172.26.127.95:9001/v1
NEXT_PUBLIC_API_ORDER_URL=http://172.26.127.95:9002/v1
NEXT_PUBLIC_API_ANALYTICS_URL=http://172.26.127.95:9004/v1

# Application
NEXT_PUBLIC_APP_URL=http://localhost:9999
NEXT_PUBLIC_PORT=9999

# Other configs...
```

---

## 🚀 Chạy Dự Án

### **Development Mode**
Chạy môi trường phát triển với hot reload:

```bash
npm run dev
```

Ứng dụng sẽ chạy tại: **http://localhost:9999**

### **Các Commands Hữu Ích**

```bash
# Chạy development server
npm run dev

# Build production
npm run build

# Chạy production server
npm run start

# Lint code
npm run lint

# Type check
npx tsc --noEmit
```

### **Testing Features**

#### 1. **Test Theme System**
```bash
# Mở http://localhost:9999
# Click vào color selector ở header
# Chọn từng màu để xem gradient effects
```

#### 2. **Test Authentication**
```bash
# Đăng ký tài khoản mới: /vi/auth/register
# Đăng nhập: /vi/auth/login
# Token sẽ tự động refresh khi hết hạn
```

#### 3. **Test Shopping Flow**
```bash
# Duyệt sản phẩm: /vi
# Thêm vào giỏ hàng
# Thanh toán: /vi/thanh-toan
# Xem đơn hàng: /vi/profile/don-hang
```

#### 4. **Test AI Chatbot**
```bash
# Click vào icon chat ở góc dưới phải
# Nhập câu hỏi và nhận phản hồi từ AI
```

---

## 📦 Build & Deploy

### **Build cho Production**

```bash
# Build ứng dụng
npm run build

# Output sẽ được tạo trong thư mục .next/
```

### **Docker Deployment**

Dự án đã được cấu hình với `output: "standalone"` trong `next.config.mjs` để dễ dàng deploy với Docker.

**Dockerfile:**
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV production
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 9999
ENV PORT 9999
CMD ["node", "server.js"]
```

**Build & Run:**
```bash
# Build Docker image
docker build -t ecom-client:latest .

# Run container
docker run -p 9999:9999 ecom-client:latest
```

### **Deploy lên Vercel**

```bash
# Cài Vercel CLI
npm i -g vercel

# Deploy
vercel

# Deploy production
vercel --prod
```

### **Deploy lên Nginx**

```bash
# Build
npm run build

# Copy files
cp -r .next/standalone/* /var/www/ecom-client/
cp -r public /var/www/ecom-client/
cp -r .next/static /var/www/ecom-client/.next/

# Start with PM2
pm2 start node --name ecom-client -- server.js
```

**Nginx Config:**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:9999;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

---

## 📁 Cấu Trúc Thư Mục

```
ecom_client_user/
├── public/                      # Static files
│   ├── fonts/                   # Custom fonts
│   └── images/                  # Images
│
├── src/
│   ├── app/                     # Next.js 14 App Router
│   │   ├── [locale]/           # i18n routes
│   │   │   ├── (pages)/        # Main pages
│   │   │   │   ├── page.tsx    # Homepage
│   │   │   │   ├── product/    # Product pages
│   │   │   │   ├── gio-hang/   # Cart
│   │   │   │   ├── thanh-toan/ # Checkout
│   │   │   │   ├── profile/    # User profile
│   │   │   │   └── ...
│   │   │   ├── auth/           # Authentication
│   │   │   ├── layout.tsx      # Root layout
│   │   │   └── globals.css     # Global styles
│   │   └── fonts/              # Font definitions
│   │
│   ├── assets/
│   │   ├── configs/            # Configuration files
│   │   │   ├── api.ts          # API endpoints
│   │   │   ├── request.ts      # Request configs
│   │   │   ├── theme-color.ts  # Theme colors
│   │   │   └── routers.ts      # Route configs
│   │   ├── helpers/            # Helper functions
│   │   │   ├── cookies.ts      # Cookie utils
│   │   │   ├── request.ts      # API request utils
│   │   │   └── string.ts       # String utils
│   │   ├── hooks/              # Custom hooks
│   │   │   └── useDebounce.ts
│   │   ├── interface/          # TypeScript interfaces
│   │   ├── locales/            # i18n translations
│   │   │   ├── en.json
│   │   │   └── vi.json
│   │   ├── middleware/         # Middleware
│   │   └── types/              # Type definitions
│   │
│   ├── components/             # React components
│   │   ├── ui/                 # UI components
│   │   │   ├── button.tsx
│   │   │   ├── badge.tsx
│   │   │   ├── dialog.tsx
│   │   │   └── ...
│   │   ├── AddressDialog.tsx   # Address management
│   │   ├── ChatBox.tsx         # AI Chatbot
│   │   ├── VoucherSelector.tsx # Voucher component
│   │   ├── theme-provider.tsx  # Theme provider
│   │   └── ...
│   │
│   ├── hooks/                  # Custom React hooks
│   │   ├── use-toast.ts
│   │   └── useTokenRefresh.ts
│   │
│   ├── i18n/                   # Internationalization
│   │   ├── request.ts
│   │   └── routing.ts
│   │
│   ├── lib/                    # Utility libraries
│   │   ├── apiClient.ts        # Axios client with interceptors
│   │   ├── apiOrderService.ts  # Order API service
│   │   ├── theme-color.ts      # Theme color utils
│   │   ├── tokenRefreshService.ts
│   │   └── utils.ts            # General utils
│   │
│   ├── resources/              # Resources
│   │   ├── components_thuongdung/  # Common components
│   │   └── layout/             # Layout components
│   │       ├── header.tsx
│   │       ├── footer.tsx
│   │       └── theme-color-toggle.tsx
│   │
│   ├── services/               # API Services
│   │   └── apiService.ts       # TanStack Query hooks
│   │
│   ├── store/                  # State management
│   │   ├── cartStore.ts        # Zustand cart store
│   │   └── checkoutStore.ts    # Checkout store
│   │
│   ├── types/                  # TypeScript types
│   │   ├── address.types.ts
│   │   ├── cart.types.ts
│   │   ├── product.types.ts
│   │   └── ...
│   │
│   └── middleware.ts           # Next.js middleware
│
├── .env.local                  # Environment variables (not in git)
├── next.config.mjs             # Next.js configuration
├── tailwind.config.ts          # Tailwind configuration
├── tsconfig.json               # TypeScript configuration
├── package.json                # Dependencies
├── Dockerfile                  # Docker configuration
└── README.md                   # This file
```

---

## 🎨 Theme System

Ứng dụng có hệ thống theme mạnh mẽ với **9 bảng màu** và **gradient effects**.

### **Available Themes**
1. **Orange** - Năng động, phù hợp thương mại điện tử
2. **Blue** - Tin cậy, công nghệ
3. **Green** - Tự nhiên, sức khỏe
4. **Rose** - Thời trang, làm đẹp
5. **Zinc** - Chuyên nghiệp, tối giản
6. **Purple** - Sáng tạo, sang trọng
7. **Cyan** - Hiện đại, tươi mới
8. **Yellow** - Tích cực, năng lượng
9. **Teal** - Cân bằng, hài hòa

### **Features**
- ✅ **Gradient backgrounds** - 2-3 màu pha trộn hài hòa
- ✅ **Light & Dark mode** - Tự động hoặc thủ công
- ✅ **Smooth transitions** - 300-400ms
- ✅ **No white overlay** - Text luôn rõ ràng
- ✅ **Custom hover effects** - 6 loại hiệu ứng
- ✅ **Border radius** - Bo tròn mềm mại 12-16px

### **Usage**
```tsx
import { useThemeContext } from "@/components/theme-color-provider";

function MyComponent() {
  const { themeColor, setThemeColor } = useThemeContext();
  
  return (
    <button onClick={() => setThemeColor("Purple")}>
      Change to Purple Theme
    </button>
  );
}
```

### **Documentation**
- 📖 [COLOR_SYSTEM_GUIDE.md](./COLOR_SYSTEM_GUIDE.md) - Hướng dẫn chi tiết
- 📖 [GRADIENT_IMPROVEMENTS_SUMMARY.md](./GRADIENT_IMPROVEMENTS_SUMMARY.md) - Tổng quan cải tiến
- 📖 [COLOR_QUICK_REFERENCE.md](./COLOR_QUICK_REFERENCE.md) - Tham khảo nhanh

---

## 🧪 Testing

### **Manual Testing**
```bash
# Test các tính năng chính
1. Homepage - Load products
2. Search - Tìm kiếm sản phẩm
3. Product Detail - Xem chi tiết
4. Add to Cart - Thêm giỏ hàng
5. Checkout - Thanh toán
6. Profile - Quản lý tài khoản
7. Theme Switch - Đổi màu theme
8. Language Switch - Đổi ngôn ngữ
9. AI Chatbot - Chat với AI
```

### **API Testing**
Sử dụng Postman hoặc cURL để test APIs:

```bash
# Test login
curl -X POST https://lemarchenoble.id.vn/api/v1/identity/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test@example.com","password":"password123"}'

# Test get products
curl -X GET http://172.26.127.95:9001/v1/product/getall \
  -H "Authorization: Bearer <token>"
```

---

## 🔧 Troubleshooting

### **Common Issues**

#### 1. **Port 9999 đã được sử dụng**
```bash
# Tìm process đang dùng port
lsof -i :9999

# Kill process
kill -9 <PID>

# Hoặc đổi port trong package.json
"dev": "next dev -p 3000"
```

#### 2. **Dependencies lỗi**
```bash
# Xóa node_modules và reinstall
rm -rf node_modules package-lock.json
npm install
```

#### 3. **API Connection Failed**
```bash
# Check API services đang chạy
curl http://localhost:9000/api/session

# Check network config
# Đảm bảo không bị firewall block
```

#### 4. **Build Error**
```bash
# Clear .next cache
rm -rf .next

# Rebuild
npm run build
```

---

## 📊 Performance

### **Lighthouse Scores** (Target)
- 🟢 Performance: 90+
- 🟢 Accessibility: 95+
- 🟢 Best Practices: 90+
- 🟢 SEO: 95+

### **Optimizations**
- ✅ Next.js Image optimization
- ✅ Code splitting & lazy loading
- ✅ CSS optimization với Tailwind
- ✅ API caching với TanStack Query
- ✅ Static page generation
- ✅ Compression & minification

---

## 🔐 Security

### **Implemented**
- ✅ JWT Authentication
- ✅ Auto token refresh
- ✅ HTTP-only cookies
- ✅ CSRF protection
- ✅ Input validation với Zod
- ✅ Sanitized user inputs
- ✅ Secure headers

### **Best Practices**
```typescript
// Không lưu sensitive data trong localStorage
// Dùng HTTP-only cookies cho tokens
// Validate tất cả inputs
// Sanitize HTML content
// Use HTTPS in production
```

---

## 🤝 Đóng Góp

Chúng tôi hoan nghênh mọi đóng góp! Để đóng góp:

### **Steps**
1. Fork repository
2. Tạo branch mới (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Mở Pull Request

### **Coding Standards**
- TypeScript strict mode
- ESLint rules
- Prettier formatting
- Meaningful commit messages
- Component documentation

---

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## 👥 Team

- **Product Owner** - [TranVinhHien](https://github.com/TranVinhHien)
- **Developer** - Full Stack Team
- **Designer** - UI/UX Team

---

## 📞 Contact & Support

- **Repository:** https://github.com/TranVinhHien/ecom_user_client
- **Issues:** https://github.com/TranVinhHien/ecom_user_client/issues
- **Email:** support@lemarchenoble.id.vn
- **Website:** https://lemarchenoble.id.vn

---

## 🎉 Acknowledgments

- Next.js Team
- Radix UI
- Tailwind CSS
- Vercel
- All contributors

---

<div align="center">

### ⭐ If you like this project, give it a star!

Made with ❤️ by Le Marché Noble Team

</div>
