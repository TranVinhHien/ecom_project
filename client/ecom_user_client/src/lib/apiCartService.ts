import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { cookies } from "@/assets/helpers"
import { ACCESS_TOKEN } from '@/assets/configs/request';
import { jwtDecode } from 'jwt-decode';
import API from '@/assets/configs/api';

// Biến để theo dõi việc refresh token đang diễn ra
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: any) => void;
  reject: (reason?: any) => void;
}> = [];

const processQueue = (error: AxiosError | null, token: string | null = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  
  failedQueue = [];
};

// Tạo instance axios cho Cart & Shop Service
const apiCartClient = axios.create({
  baseURL: API.base_cart_shop,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - Tự động đính kèm Bearer Token nếu có
apiCartClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Lấy token từ cookies
    const token = cookies.getCookieValues<string>(ACCESS_TOKEN);
    
    console.log("🔑 API CART SERVICE - Token:", token ? "✅ Có token" : "❌ Không có token");
    
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    } else {
      // Không có token - Đánh dấu request này là từ guest user
      if (config.headers) {
        config.headers['X-Guest-Request'] = 'true';
      }
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor - Xử lý lỗi và refresh token tự động
apiCartClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // Xử lý lỗi 401 - Token hết hạn
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // Nếu đang refresh, đợi và thử lại với token mới
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${token}`;
          }
          return apiCartClient(originalRequest);
        }).catch(err => {
          return Promise.reject(err);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const currentToken = cookies.getCookieValues<string>(ACCESS_TOKEN);

      if (!currentToken) {
        // ❌ Không có token - Người dùng chưa đăng nhập
        // KHÔNG redirect! Component sẽ tự động dùng localStorage
        console.log('⚠️ Guest user - API call failed, will use localStorage');
        isRefreshing = false;
        return Promise.reject(error);
      }

      try {
        // Gọi API refresh token
        const refreshResponse = await axios.post(
          `${API.base_gateway}${API.user.refresh}`,
          {},
          {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${currentToken}`,
            },
          }
        );

        if (refreshResponse.data.code === 10000 && refreshResponse.data.result?.token) {
          const newToken = refreshResponse.data.result.token;
          
          // Decode token để lấy thời gian hết hạn
          const decoded: any = jwtDecode(newToken);
          
          // Lưu token mới vào cookies
          cookies.setCookieValues(ACCESS_TOKEN, newToken, decoded?.exp);
          
          console.log('✅ Token refreshed successfully (Cart Service)');
          
          // Cập nhật token trong request ban đầu
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${newToken}`;
          }
          
          // Xử lý các request đang đợi
          processQueue(null, newToken);
          isRefreshing = false;
          
          // Thử lại request ban đầu
          return apiCartClient(originalRequest);
        } else {
          throw new Error('Invalid refresh response');
        }
      } catch (refreshError) {
        console.error('❌ Failed to refresh token:', refreshError);
        processQueue(refreshError as AxiosError, null);
        isRefreshing = false;
        
        // Chỉ logout và redirect nếu token đã tồn tại trước đó (user đã đăng nhập)
        // Điều này có nghĩa là token thật sự hết hạn, không phải guest user
        if (currentToken && typeof window !== 'undefined') {
          console.log('🔐 Token expired - Logging out...');
          cookies.logOut();
          window.location.href = '/vi/auth/login';
        }
        
        return Promise.reject(refreshError);
      }
    }
    
    return Promise.reject(error);
  }
);

export default apiCartClient;
