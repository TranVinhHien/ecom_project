"use client";

import { useState } from 'react';
import { useTokenInfo } from '@/hooks/useTokenRefresh';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Clock, RefreshCw, AlertTriangle } from 'lucide-react';
import { tokenRefreshService } from '@/lib/tokenRefreshService';

/**
 * Token Status Display Component (Optional)
 * Shows token expiration time and allows manual refresh
 * Use this for debugging or admin panels
 */
export default function TokenStatus() {
  const { timeRemaining, isExpiringSoon, formattedTime } = useTokenInfo();
  const [isRefreshing, setIsRefreshing] = useState(false);

  const handleManualRefresh = async () => {
    setIsRefreshing(true);
    await tokenRefreshService.manualRefresh();
    setIsRefreshing(false);
  };

  // Don't show if no token
  if (timeRemaining === 0) return null;

  return (
    <Card className="fixed bottom-20 right-6 p-3 shadow-lg z-40 bg-white/95 backdrop-blur">
      <div className="flex items-center gap-3">
        {/* Status Icon */}
        <div className={`p-2 rounded-full ${isExpiringSoon ? 'bg-orange-100' : 'bg-green-100'}`}>
          {isExpiringSoon ? (
            <AlertTriangle className="h-4 w-4 text-orange-600" />
          ) : (
            <Clock className="h-4 w-4 text-green-600" />
          )}
        </div>

        {/* Token Info */}
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-2">
            <span className="text-xs font-semibold text-gray-700">Token expires in:</span>
            <Badge variant={isExpiringSoon ? "destructive" : "secondary"} className="text-xs">
              {formattedTime}
            </Badge>
          </div>
          
          {isExpiringSoon && (
            <span className="text-xs text-orange-600">
              Auto-refreshing soon...
            </span>
          )}
        </div>

        {/* Manual Refresh Button */}
        <Button
          variant="outline"
          size="sm"
          onClick={handleManualRefresh}
          disabled={isRefreshing}
          className="ml-2"
        >
          {isRefreshing ? (
            <>
              <RefreshCw className="h-3 w-3 mr-1 animate-spin" />
              <span className="text-xs">Refreshing...</span>
            </>
          ) : (
            <>
              <RefreshCw className="h-3 w-3 mr-1" />
              <span className="text-xs">Refresh</span>
            </>
          )}
        </Button>
      </div>
    </Card>
  );
}
