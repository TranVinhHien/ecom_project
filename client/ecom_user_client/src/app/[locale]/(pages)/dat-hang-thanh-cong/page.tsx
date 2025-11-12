"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { CheckCircle2, XCircle, AlertCircle, Clock } from "lucide-react";
import { Link } from '@/i18n/routing';
import { useRouter } from "@/i18n/routing";

// MoMo Result Code Error Mapping
const MOMO_ERROR_CODES: Record<string, { 
  description: string; 
  action: string; 
  type: 'success' | 'error' | 'warning' | 'pending';
}> = {
  "0": { 
    description: "Thành công.", 
    action: "", 
    type: 'success' 
  },
  "10": { 
    description: "Hệ thống đang được bảo trì.", 
    action: "Vui lòng quay lại sau khi bảo trì được hoàn tất.", 
    type: 'error' 
  },
  "11": { 
    description: "Truy cập bị từ chối.", 
    action: "Cấu hình tài khoản doanh nghiệp không cho phép truy cập. Vui lòng liên hệ với MoMo để được hỗ trợ.", 
    type: 'error' 
  },
  "12": { 
    description: "Phiên bản API không được hỗ trợ.", 
    action: "Vui lòng nâng cấp lên phiên bản mới nhất.", 
    type: 'error' 
  },
  "13": { 
    description: "Xác thực doanh nghiệp thất bại.", 
    action: "Vui lòng kiểm tra thông tin kết nối.", 
    type: 'error' 
  },
  "20": { 
    description: "Yêu cầu sai định dạng.", 
    action: "Vui lòng kiểm tra định dạng của yêu cầu.", 
    type: 'error' 
  },
  "21": { 
    description: "Số tiền giao dịch không hợp lệ.", 
    action: "Vui lòng kiểm tra số tiền hợp lệ và thực hiện lại.", 
    type: 'error' 
  },
  "22": { 
    description: "Số tiền giao dịch không hợp lệ.", 
    action: "Vui lòng kiểm tra số tiền thanh toán.", 
    type: 'error' 
  },
  "40": { 
    description: "RequestId bị trùng.", 
    action: "Vui lòng thử lại với một requestId khác.", 
    type: 'error' 
  },
  "41": { 
    description: "OrderId bị trùng.", 
    action: "Vui lòng thử lại với một orderId khác.", 
    type: 'error' 
  },
  "42": { 
    description: "OrderId không hợp lệ hoặc không được tìm thấy.", 
    action: "Vui lòng thử lại với một orderId khác.", 
    type: 'error' 
  },
  "43": { 
    description: "Yêu cầu bị từ chối vì xung đột trong quá trình xử lý giao dịch.", 
    action: "Vui lòng kiểm tra và thử lại.", 
    type: 'error' 
  },
  "98": { 
    description: "QR Code tạo không thành công.", 
    action: "Vui lòng thử lại sau.", 
    type: 'error' 
  },
  "99": { 
    description: "Lỗi không xác định.", 
    action: "Vui lòng liên hệ MoMo để biết thêm chi tiết.", 
    type: 'error' 
  },
  "1000": { 
    description: "Giao dịch đã được khởi tạo, chờ người dùng xác nhận thanh toán.", 
    action: "Vui lòng hoàn tất thanh toán trên ví MoMo.", 
    type: 'pending' 
  },
  "1001": { 
    description: "Giao dịch thất bại do tài khoản không đủ tiền.", 
    action: "Vui lòng kiểm tra số dư và thử lại.", 
    type: 'error' 
  },
  "1002": { 
    description: "Giao dịch bị từ chối do nhà phát hành tài khoản thanh toán.", 
    action: "Vui lòng sử dụng phương thức thanh toán khác.", 
    type: 'error' 
  },
  "1003": { 
    description: "Giao dịch đã bị hủy.", 
    action: "Giao dịch đã bị hủy bởi hệ thống hoặc người dùng.", 
    type: 'error' 
  },
  "1004": { 
    description: "Giao dịch thất bại do vượt quá hạn mức thanh toán.", 
    action: "Vui lòng thử lại vào thời gian khác hoặc giảm số tiền giao dịch.", 
    type: 'error' 
  },
  "1005": { 
    description: "Giao dịch thất bại do URL hoặc QR code đã hết hạn.", 
    action: "Vui lòng thực hiện lại giao dịch mới.", 
    type: 'error' 
  },
  "1006": { 
    description: "Người dùng đã từ chối xác nhận thanh toán.", 
    action: "Vui lòng thử lại nếu muốn tiếp tục thanh toán.", 
    type: 'error' 
  },
  "1007": { 
    description: "Tài khoản không tồn tại hoặc đang ngưng hoạt động.", 
    action: "Vui lòng kiểm tra tài khoản hoặc liên hệ MoMo.", 
    type: 'error' 
  },
  "1017": { 
    description: "Giao dịch bị hủy bởi đối tác.", 
    action: "Giao dịch đã bị hủy.", 
    type: 'error' 
  },
  "1026": { 
    description: "Giao dịch bị hạn chế theo thể lệ chương trình khuyến mãi.", 
    action: "Vui lòng liên hệ MoMo để biết thêm chi tiết.", 
    type: 'error' 
  },
  "1080": { 
    description: "Giao dịch hoàn tiền thất bại.", 
    action: "Vui lòng thử lại sau.", 
    type: 'error' 
  },
  "4001": { 
    description: "Tài khoản đang bị hạn chế.", 
    action: "Vui lòng liên hệ MoMo để biết thêm chi tiết.", 
    type: 'error' 
  },
  "4100": { 
    description: "Người dùng không đăng nhập thành công.", 
    action: "Vui lòng thử lại.", 
    type: 'error' 
  },
  "7000": { 
    description: "Giao dịch đang được xử lý.", 
    action: "Vui lòng chờ giao dịch được xử lý hoàn tất.", 
    type: 'pending' 
  },
  "7002": { 
    description: "Giao dịch đang được xử lý bởi nhà cung cấp.", 
    action: "Vui lòng chờ. Kết quả sẽ được thông báo khi hoàn tất.", 
    type: 'pending' 
  },
  "9000": { 
    description: "Giao dịch đã được xác nhận thành công.", 
    action: "Giao dịch đã được xác nhận, vui lòng chờ xử lý.", 
    type: 'success' 
  },
};

export default function OrderSuccessPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  
  const [paymentInfo, setPaymentInfo] = useState<{
    resultCode: string;
    orderId: string;
    amount: string;
    message: string;
    transId: string;
    orderInfo: string;
    partnerCode: string;
  } | null>(null);

  useEffect(() => {
    // Get payment info from URL params
    const resultCode = searchParams.get('resultCode') || '0';
    const orderId = searchParams.get('orderId') || '';
    const amount = searchParams.get('amount') || '';
    const message = searchParams.get('message') || '';
    const transId = searchParams.get('transId') || '';
    const orderInfo = searchParams.get('orderInfo') || '';
    const partnerCode = searchParams.get('partnerCode') || '';

    setPaymentInfo({
      resultCode,
      orderId,
      amount,
      message: decodeURIComponent(message),
      transId,
      orderInfo: decodeURIComponent(orderInfo),
      partnerCode,
    });

    console.log('Payment Result Code:', resultCode);
    console.log('Payment Info:', { orderId, amount, transId });
  }, [searchParams]);

  const formatAmount = (amount: string) => {
    const num = parseInt(amount);
    if (isNaN(num)) return '0 VNĐ';
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
    }).format(num);
  };

  const getErrorInfo = (code: string) => {
    return MOMO_ERROR_CODES[code] || {
      description: "Lỗi không xác định",
      action: "Vui lòng liên hệ với chúng tôi để được hỗ trợ.",
      type: 'error' as const,
    };
  };

  const renderIcon = (type: string) => {
    switch (type) {
      case 'success':
        return <CheckCircle2 className="w-24 h-24 text-green-500" />;
      case 'pending':
        return <Clock className="w-24 h-24 text-yellow-500" />;
      case 'warning':
        return <AlertCircle className="w-24 h-24 text-orange-500" />;
      default:
        return <XCircle className="w-24 h-24 text-red-500" />;
    }
  };

  if (!paymentInfo) {
    return (
      <div className="container mx-auto px-4 py-16">
        <Card className="max-w-2xl mx-auto text-center">
          <CardContent className="pt-12 pb-8">
            <p>Đang tải thông tin...</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const errorInfo = getErrorInfo(paymentInfo.resultCode);
  const isSuccess = paymentInfo.resultCode === '0' || paymentInfo.resultCode === '9000';
  const isPending = errorInfo.type === 'pending';

  return (
    <div className="container mx-auto px-4 py-16">
      <Card className="max-w-2xl mx-auto text-center">
        <CardContent className="pt-12 pb-8">
          <div className="flex justify-center mb-6">
            {renderIcon(errorInfo.type)}
          </div>

          {isSuccess ? (
            <>
              <h1 className="text-3xl font-bold mb-4 text-gray-800">
                Đặt hàng thành công!
              </h1>

              <p className="text-gray-600 mb-2">
                Cảm ơn bạn đã đặt hàng tại cửa hàng của chúng tôi.
              </p>
              <p className="text-gray-600 mb-8">
                Đơn hàng của bạn đã được xác nhận và đang được xử lý.
              </p>
            </>
          ) : isPending ? (
            <>
              <h1 className="text-3xl font-bold mb-4 text-yellow-600">
                Đơn hàng đang chờ xử lý
              </h1>

              <p className="text-gray-600 mb-4">
                {errorInfo.description}
              </p>
              <p className="text-gray-600 mb-8">
                {errorInfo.action}
              </p>
            </>
          ) : (
            <>
              <h1 className="text-3xl font-bold mb-4 text-red-600">
                Thanh toán thất bại
              </h1>

              <p className="text-gray-600 mb-2">
                <strong>Mã lỗi:</strong> {paymentInfo.resultCode}
              </p>
              <p className="text-gray-600 mb-4">
                <strong>Mô tả:</strong> {errorInfo.description}
              </p>
              <p className="text-gray-600 mb-8">
                <strong>Khuyến nghị:</strong> {errorInfo.action}
              </p>
            </>
          )}

          {/* Payment Details */}
          {paymentInfo.partnerCode === 'MOMO' && (
            <div className={`border rounded-lg p-4 mb-8 ${
              isSuccess ? 'bg-green-50 border-green-200' : 
              isPending ? 'bg-yellow-50 border-yellow-200' : 
              'bg-red-50 border-red-200'
            }`}>
              <h3 className="font-semibold mb-3 text-left">Thông tin thanh toán MoMo</h3>
              <div className="space-y-2 text-sm text-left">
                {paymentInfo.orderId && (
                  <div className="flex justify-between">
                    <span className="text-gray-600">Mã đơn hàng:</span>
                    <span className="font-medium">{paymentInfo.orderId}</span>
                  </div>
                )}
                {paymentInfo.transId && (
                  <div className="flex justify-between">
                    <span className="text-gray-600">Mã giao dịch:</span>
                    <span className="font-medium">{paymentInfo.transId}</span>
                  </div>
                )}
                {paymentInfo.amount && (
                  <div className="flex justify-between">
                    <span className="text-gray-600">Số tiền:</span>
                    <span className="font-medium text-primary">{formatAmount(paymentInfo.amount)}</span>
                  </div>
                )}
                {paymentInfo.message && paymentInfo.resultCode !== '0' && (
                  <div className="flex flex-col gap-1 pt-2 border-t">
                    <span className="text-gray-600">Thông báo từ MoMo:</span>
                    <span className="text-xs italic">{paymentInfo.message}</span>
                  </div>
                )}
              </div>
            </div>
          )}

          {isSuccess && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-8">
              <p className="text-sm text-blue-800">
                📧 Chúng tôi đã gửi email xác nhận đơn hàng đến địa chỉ email của bạn.
              </p>
            </div>
          )}

          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link href="/">
              <Button
                variant="outline"
                size="lg"
                className="w-full sm:w-auto"
              >
                Về trang chủ
              </Button>
            </Link>

            {isSuccess ? (
              <Link href="/don-hang">
                <Button
                  className="w-full sm:w-auto bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)]"
                  size="lg"
                >
                  Xem đơn hàng
                </Button>
              </Link>
            ) : (
              <Link href="/gio-hang">
                <Button
                  className="w-full sm:w-auto bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)]"
                  size="lg"
                >
                  Quay lại giỏ hàng
                </Button>
              </Link>
            )}
          </div>

          {isSuccess && (
            <div className="mt-8 pt-8 border-t">
              <h3 className="font-semibold mb-4">Tiếp theo là gì?</h3>
              <ul className="text-sm text-gray-600 space-y-2">
                <li>✅ Đơn hàng của bạn đang được chuẩn bị</li>
                <li>📦 Bạn sẽ nhận được thông báo khi đơn hàng được giao cho đơn vị vận chuyển</li>
                <li>🚚 Thời gian giao hàng dự kiến: 2-5 ngày làm việc</li>
              </ul>
            </div>
          )}

          {!isSuccess && !isPending && (
            <div className="mt-8 pt-8 border-t">
              <h3 className="font-semibold mb-4">Cần hỗ trợ?</h3>
              <p className="text-sm text-gray-600 mb-4">
                Nếu bạn gặp vấn đề, vui lòng liên hệ với chúng tôi:
              </p>
              <div className="text-sm text-gray-600 space-y-1">
                <p>📞 Hotline: 1900 xxxx</p>
                <p>📧 Email: support@example.com</p>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
