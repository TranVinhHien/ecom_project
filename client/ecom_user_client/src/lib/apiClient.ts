import axios, { AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios';
import { getCookieValues, setCookieValues } from '@/assets/helpers/cookies';
import { ACCESS_TOKEN } from '@/assets/configs/request';
import API from '@/assets/configs/api';
import { jwtDecode } from 'jwt-decode';

// ...existing code...

// ✨ Extend AxiosRequestConfig để thêm customBaseURL
interface CustomAxiosRequestConfig extends AxiosRequestConfig {
  customBaseURL?: string;
}

class ApiClient {
  private client: AxiosInstance;
  private isRefreshing: boolean = false;
  private failedQueue: Array<{
    resolve: (value?: any) => void;
    reject: (reason?: any) => void;
  }> = [];

  constructor(baseURL: string) {
    this.client = axios.create({
      baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request Interceptor
    this.client.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // ✨ Kiểm tra nếu có customBaseURL trong config
        const customConfig = config as InternalAxiosRequestConfig & { customBaseURL?: string };
        
        if (customConfig.customBaseURL) {
          // Thay đổi baseURL cho request này
          config.baseURL = customConfig.customBaseURL;
          console.log(`🔄 Using custom baseURL: ${customConfig.customBaseURL}`);
        }

        // Gắn token vào header
        const token = getCookieValues<string>(ACCESS_TOKEN);
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response Interceptor
    this.client.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & {
          _retry?: boolean;
          customBaseURL?: string;
        };

        // Nếu lỗi 401 và chưa retry
        if (error.response?.status === 401 && !originalRequest._retry) {
          if (this.isRefreshing) {
            // Đang refresh, thêm vào queue
            return new Promise((resolve, reject) => {
              this.failedQueue.push({ resolve, reject });
            })
              .then(() => {
                // ✨ Giữ lại customBaseURL khi retry
                const retryConfig = { ...originalRequest };
                if (originalRequest.customBaseURL) {
                  retryConfig.baseURL = originalRequest.customBaseURL;
                }
                return this.client(retryConfig);
              })
              .catch((err) => Promise.reject(err));
          }

          originalRequest._retry = true;
          this.isRefreshing = true;

          try {
            const oldToken = getCookieValues<string>(ACCESS_TOKEN);
            
            if (!oldToken) {
              throw new Error('No token available');
            }

            // Gọi API refresh token
            const response = await axios.post(
              `${API.base_vinh}${API.user.refresh}`,
              { token: oldToken }
            );

            if (response.data.code === 10000 && response.data.result?.token) {
              const newToken = response.data.result.token;
              
              // Decode để lấy expiry
              const decoded: any = jwtDecode(newToken);
              
              // Lưu token mới
              setCookieValues(ACCESS_TOKEN, newToken, decoded?.exp);

              console.log('✅ Token refreshed successfully');

              // Process queue
              this.failedQueue.forEach((prom) => prom.resolve());
              this.failedQueue = [];

              // ✨ Retry request ban đầu với customBaseURL (nếu có)
              const retryConfig = { ...originalRequest };
              if (originalRequest.customBaseURL) {
                retryConfig.baseURL = originalRequest.customBaseURL;
              }
              return this.client(retryConfig);
            } else {
              throw new Error('Refresh token failed');
            }
          } catch (refreshError) {
            console.error('❌ Refresh token failed:', refreshError);
            
            this.failedQueue.forEach((prom) => prom.reject(refreshError));
            this.failedQueue = [];

            // Clear token và redirect
            document.cookie = `${ACCESS_TOKEN}=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
            localStorage.removeItem('INFO_USER');
            
            if (typeof window !== 'undefined') {
              window.location.href = '/vi/auth/login';

            }

            return Promise.reject(refreshError);
          } finally {
            this.isRefreshing = false;
          }
        }

        return Promise.reject(error);
      }
    );
  }

  // ✨ Các method GET, POST, PUT, DELETE với CustomAxiosRequestConfig
  async get<T = any>(url: string, config?: CustomAxiosRequestConfig) {
    return this.client.get<T>(url, config);
  }

  async post<T = any>(url: string, data?: any, config?: CustomAxiosRequestConfig) {
    return this.client.post<T>(url, data, config);
  }

  async put<T = any>(url: string, data?: any, config?: CustomAxiosRequestConfig) {
    return this.client.put<T>(url, data, config);
  }

  async delete<T = any>(url: string, config?: CustomAxiosRequestConfig) {
    return this.client.delete<T>(url, config);
  }

  async patch<T = any>(url: string, data?: any, config?: CustomAxiosRequestConfig) {
    return this.client.patch<T>(url, data, config);
  }
}

// Export instance
const apiClient = new ApiClient(API.base_vinh);
export { apiClient };
export default apiClient;