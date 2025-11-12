"use client";

import { useEffect, useState } from 'react';
import { tokenRefreshService } from '@/lib/tokenRefreshService';

/**
 * Hook to manage automatic token refresh
 * Usage: Add this hook to your root layout or main app component
 */
export function useTokenRefresh() {
  const [isTokenValid, setIsTokenValid] = useState(true);
  const [timeRemaining, setTimeRemaining] = useState<number>(0);

  useEffect(() => {
    // Initialize token refresh service
    tokenRefreshService.initialize();

    // Update time remaining every minute
    const interval = setInterval(() => {
      const remaining = tokenRefreshService.getTokenTimeRemaining();
      setTimeRemaining(remaining);
      setIsTokenValid(remaining > 0);
    }, 60000); // Check every minute

    // Listen for token refresh events
    const handleTokenRefresh = () => {
      console.log('🔔 Token refreshed event received');
      setIsTokenValid(true);
      setTimeRemaining(tokenRefreshService.getTokenTimeRemaining());
    };

    window.addEventListener('tokenRefreshed', handleTokenRefresh);

    // Initial check
    setTimeRemaining(tokenRefreshService.getTokenTimeRemaining());
    setIsTokenValid(tokenRefreshService.getTokenTimeRemaining() > 0);

    // Cleanup
    return () => {
      tokenRefreshService.stopScheduler();
      clearInterval(interval);
      window.removeEventListener('tokenRefreshed', handleTokenRefresh);
    };
  }, []);

  return {
    isTokenValid,
    timeRemaining,
    manualRefresh: () => tokenRefreshService.manualRefresh(),
  };
}

/**
 * Hook to get token expiration info
 */
export function useTokenInfo() {
  const [timeRemaining, setTimeRemaining] = useState<number>(0);
  const [isExpiringSoon, setIsExpiringSoon] = useState(false);

  useEffect(() => {
    const updateTokenInfo = () => {
      const remaining = tokenRefreshService.getTokenTimeRemaining();
      const expiringSoon = tokenRefreshService.isTokenExpiringSoon();
      
      setTimeRemaining(remaining);
      setIsExpiringSoon(expiringSoon);
    };

    // Initial update
    updateTokenInfo();

    // Update every 30 seconds
    const interval = setInterval(updateTokenInfo, 30000);

    // Listen for token refresh
    window.addEventListener('tokenRefreshed', updateTokenInfo);

    return () => {
      clearInterval(interval);
      window.removeEventListener('tokenRefreshed', updateTokenInfo);
    };
  }, []);

  // Format time remaining as human-readable string
  const formatTimeRemaining = () => {
    if (timeRemaining <= 0) return 'Expired';
    
    const hours = Math.floor(timeRemaining / 3600);
    const minutes = Math.floor((timeRemaining % 3600) / 60);
    const seconds = timeRemaining % 60;

    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds}s`;
    } else {
      return `${seconds}s`;
    }
  };

  return {
    timeRemaining,
    isExpiringSoon,
    formattedTime: formatTimeRemaining(),
  };
}
