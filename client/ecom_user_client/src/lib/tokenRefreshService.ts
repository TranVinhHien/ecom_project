/**
 * Token Refresh Service
 * Automatically refresh access token before it expires
 */

import { jwtDecode } from 'jwt-decode';
import { getCookieValues, setCookieValues, logOut } from '@/assets/helpers/cookies';
import { ACCESS_TOKEN } from '@/assets/configs/request';
import API from '@/assets/configs/api';

interface JWTPayload {
  exp: number;
  iat: number;
  sub: string;
  userId: string;
  email: string;
  scope: string;
  iss: string;
  jti: string;
}

class TokenRefreshService {
  private refreshTimer: NodeJS.Timeout | null = null;
  private isRefreshing: boolean = false;
  private readonly REFRESH_BUFFER_TIME = 5 * 60 * 1000; // 5 minutes before expiry

  /**
   * Initialize token refresh scheduler
   * Call this when app starts or user logs in
   */
  public initialize(): void {
    if (typeof window === 'undefined') return; // Only run on client side

    console.log('🔄 Initializing Token Refresh Service...');
    
    // Clear any existing timer
    this.stopScheduler();
    
    // Check token and schedule refresh
    this.scheduleTokenRefresh();
  }

  /**
   * Stop the refresh scheduler
   * Call this on logout or unmount
   */
  public stopScheduler(): void {
    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer);
      this.refreshTimer = null;
      console.log('⏹️ Token refresh scheduler stopped');
    }
  }

  /**
   * Schedule token refresh based on expiration time
   */
  private scheduleTokenRefresh(): void {
    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      
      if (!token) {
        console.log('⚠️ No token found, skipping refresh schedule');
        return;
      }

      // Decode token to get expiration time
      const decoded = jwtDecode<JWTPayload>(token);
      const expirationTime = decoded.exp * 1000; // Convert to milliseconds
      const currentTime = Date.now();
      const timeUntilExpiry = expirationTime - currentTime;
      const timeUntilRefresh = timeUntilExpiry - this.REFRESH_BUFFER_TIME;

      console.log(`📊 Token Info:
        - Current Time: ${new Date(currentTime).toLocaleString()}
        - Expires At: ${new Date(expirationTime).toLocaleString()}
        - Time Until Expiry: ${Math.floor(timeUntilExpiry / 1000 / 60)} minutes
        - Will Refresh In: ${Math.floor(timeUntilRefresh / 1000 / 60)} minutes
      `);

      // If token expires in less than buffer time, refresh immediately
      if (timeUntilRefresh <= 0) {
        console.log('⚡ Token expiring soon, refreshing immediately...');
        this.refreshToken();
        return;
      }

      // Schedule refresh before expiration
      this.refreshTimer = setTimeout(() => {
        console.log('⏰ Scheduled token refresh triggered');
        this.refreshToken();
      }, timeUntilRefresh);

      console.log(`✅ Token refresh scheduled for ${new Date(currentTime + timeUntilRefresh).toLocaleString()}`);

    } catch (error) {
      console.error('❌ Error scheduling token refresh:', error);
      // If token is invalid, logout
      logOut();
    }
  }

  /**
   * Refresh the access token
   */
  private async refreshToken(): Promise<void> {
    // Prevent multiple simultaneous refresh attempts
    if (this.isRefreshing) {
      console.log('🔒 Token refresh already in progress...');
      return;
    }

    this.isRefreshing = true;

    try {
      const currentToken = getCookieValues<string>(ACCESS_TOKEN);

      if (!currentToken) {
        console.log('❌ No token to refresh, logging out...');
        logOut();
        return;
      }

      console.log('🔄 Refreshing token...');

      // Call refresh API
      const response = await fetch(`${API.base_gateway}${API.user.refresh}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          token: currentToken,
        }),
      });

      const data = await response.json();

      if (data.code === 10000 && data.result?.token) {
        const newToken = data.result.token;

        // Decode new token to get expiration
        const decoded = jwtDecode<JWTPayload>(newToken);

        // Save new token to cookies
        setCookieValues(ACCESS_TOKEN, newToken, decoded.exp);

        console.log(`✅ Token refreshed successfully!
          - New Token: ${newToken.substring(0, 20)}...
          - Expires At: ${new Date(decoded.exp * 1000).toLocaleString()}
        `);

        // Schedule next refresh
        this.scheduleTokenRefresh();

        // Dispatch event for other components
        window.dispatchEvent(new CustomEvent('tokenRefreshed', { 
          detail: { token: newToken } 
        }));

      } else {
        console.error('❌ Token refresh failed:', data);
        throw new Error(data.message || 'Token refresh failed');
      }

    } catch (error) {
      console.error('❌ Error refreshing token:', error);
      
      // If refresh fails, logout user
      console.log('🚪 Logging out due to refresh failure...');
      logOut();

    } finally {
      this.isRefreshing = false;
    }
  }

  /**
   * Manually trigger token refresh
   * Can be called from components if needed
   */
  public async manualRefresh(): Promise<boolean> {
    try {
      await this.refreshToken();
      return true;
    } catch (error) {
      console.error('Manual refresh failed:', error);
      return false;
    }
  }

  /**
   * Check if token is about to expire
   * Returns true if token expires within buffer time
   */
  public isTokenExpiringSoon(): boolean {
    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      if (!token) return true;

      const decoded = jwtDecode<JWTPayload>(token);
      const expirationTime = decoded.exp * 1000;
      const currentTime = Date.now();
      const timeUntilExpiry = expirationTime - currentTime;

      return timeUntilExpiry <= this.REFRESH_BUFFER_TIME;

    } catch (error) {
      console.error('Error checking token expiration:', error);
      return true;
    }
  }

  /**
   * Get remaining time until token expires (in seconds)
   */
  public getTokenTimeRemaining(): number {
    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      if (!token) return 0;

      const decoded = jwtDecode<JWTPayload>(token);
      const expirationTime = decoded.exp * 1000;
      const currentTime = Date.now();
      
      return Math.max(0, Math.floor((expirationTime - currentTime) / 1000));

    } catch (error) {
      console.error('Error getting token time remaining:', error);
      return 0;
    }
  }
}

// Export singleton instance
export const tokenRefreshService = new TokenRefreshService();

// Auto-initialize on client side
if (typeof window !== 'undefined') {
  // Initialize on page load
  window.addEventListener('load', () => {
    tokenRefreshService.initialize();
  });

  // Re-initialize when user becomes active after being idle
  document.addEventListener('visibilitychange', () => {
    if (!document.hidden) {
      console.log('👀 Page visible again, checking token...');
      if (tokenRefreshService.isTokenExpiringSoon()) {
        tokenRefreshService.manualRefresh();
      }
    }
  });
}

export default tokenRefreshService;
