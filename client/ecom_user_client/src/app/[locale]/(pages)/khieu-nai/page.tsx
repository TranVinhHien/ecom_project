"use client";

import { useState, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";
import { AlertCircle, CheckCircle2, Loader2 } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { ACCESS_TOKEN, INFO_USER } from "@/assets/configs/request";
import { cookies } from "@/assets/helpers";
import API from "@/assets/configs/api";

type ComplaintCategory = "BUG" | "COMPLAINT" | "SUGGESTION" | "OTHER";

interface ComplaintForm {
  phone: string;
  category: ComplaintCategory;
  content: string;
}

const categoryLabels: Record<ComplaintCategory, string> = {
  BUG: "Báo lỗi kỹ thuật",
  COMPLAINT: "Khiếu nại",
  SUGGESTION: "Đề xuất cải thiện",
  OTHER: "Khác",
};

const categoryDescriptions: Record<ComplaintCategory, string> = {
  BUG: "Báo cáo lỗi hoặc sự cố kỹ thuật trên website",
  COMPLAINT: "Khiếu nại về dịch vụ, sản phẩm hoặc trải nghiệm mua sắm",
  SUGGESTION: "Đề xuất ý tưởng cải thiện sản phẩm hoặc dịch vụ",
  OTHER: "Các vấn đề khác cần hỗ trợ",
};

export default function ComplaintPage() {
  const searchParams = useSearchParams();
  const { toast } = useToast();

  const [isLoading, setIsLoading] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const [feedbackId, setFeedbackId] = useState<string>("");
  const [successMessage, setSuccessMessage] = useState<string>("");

  const [formData, setFormData] = useState<ComplaintForm>({
    phone: "",
    category: "OTHER",
    content: "",
  });

  // Load phone from localStorage and pre-fill from URL params
  useEffect(() => {
    // Get phone from localStorage (INFO_USER)
    const infoUser = localStorage.getItem(INFO_USER);
    if (infoUser) {
      try {
        const userInfo = JSON.parse(infoUser);
        if (userInfo.phone_number) {
          setFormData((prev) => ({ ...prev, phone: userInfo.phone_number }));
        }
      } catch (error) {
        console.error("Failed to parse user info:", error);
      }
    }

    // Pre-fill from URL params if provided (from Agent)
    const categoryParam = searchParams.get("category");
    const contentParam = searchParams.get("content");
    const phoneParam = searchParams.get("phone");

    if (categoryParam && ["BUG", "COMPLAINT", "SUGGESTION", "OTHER"].includes(categoryParam)) {
      setFormData((prev) => ({ ...prev, category: categoryParam as ComplaintCategory }));
    }

    if (contentParam) {
      setFormData((prev) => ({ ...prev, content: decodeURIComponent(contentParam) }));
    }

    if (phoneParam) {
      setFormData((prev) => ({ ...prev, phone: phoneParam }));
    }
  }, [searchParams]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    if (!formData.phone.trim()) {
      toast({
        title: "Lỗi",
        description: "Vui lòng nhập số điện thoại",
        variant: "destructive",
      });
      return;
    }

    if (!formData.content.trim()) {
      toast({
        title: "Lỗi",
        description: "Vui lòng nhập nội dung",
        variant: "destructive",
      });
      return;
    }

    setIsLoading(true);
    setIsSuccess(false);

    try {
      const response = await fetch(`${API.base_analytics}${API.analytics.customerComplaint}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${cookies.getCookieValues(ACCESS_TOKEN)}`,
        },
        body: JSON.stringify({
          phone: formData.phone,
          category: formData.category,
          content: formData.content,
        }),
      });

      const data = await response.json();

      if (data.code === 200 && data.status === "success") {
        setIsSuccess(true);
        setFeedbackId(data.result.feedback_id);
        setSuccessMessage(data.result.message);

        // Reset form
        setFormData((prev) => ({
          phone: prev.phone, // Keep phone
          category: "OTHER",
          content: "",
        }));

        toast({
          title: "Thành công",
          description: data.result.message,
        });
      } else {
        throw new Error(data.message || "Gửi phản hồi thất bại");
      }
    } catch (error: any) {
      console.error("Submit complaint error:", error);
      toast({
        title: "Lỗi",
        description: error.message || "Không thể gửi phản hồi. Vui lòng thử lại.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (field: keyof ComplaintForm, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Gửi khiếu nại & phản hồi</h1>
        <p className="text-gray-600">
          Chúng tôi luôn lắng nghe ý kiến của bạn để cải thiện dịch vụ tốt hơn
        </p>
      </div>

      {/* Success Alert */}
      {isSuccess && (
        <Alert className="mb-6 border-green-500 bg-green-50">
          <CheckCircle2 className="h-5 w-5 text-green-600" />
          <AlertTitle className="text-green-800">Gửi thành công!</AlertTitle>
          <AlertDescription className="text-green-700">
            {successMessage}
            <br />
            <span className="text-sm font-mono mt-1 block">
              Mã phản hồi: {feedbackId}
            </span>
          </AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Thông tin phản hồi</CardTitle>
          <CardDescription>
            Vui lòng điền đầy đủ thông tin để chúng tôi có thể hỗ trợ bạn tốt nhất
          </CardDescription>
        </CardHeader>

        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-6">
            {/* Phone Number */}
            <div className="space-y-2">
              <Label htmlFor="phone">
                Số điện thoại <span className="text-red-500">*</span>
              </Label>
              <Input
                id="phone"
                type="tel"
                placeholder="Nhập số điện thoại của bạn"
                value={formData.phone}
                onChange={(e) => handleChange("phone", e.target.value)}
                disabled={isLoading}
                required
              />
              <p className="text-xs text-gray-500">
                Chúng tôi sẽ liên hệ với bạn qua số điện thoại này
              </p>
            </div>

            {/* Category */}
            <div className="space-y-2">
              <Label htmlFor="category">
                Loại phản hồi <span className="text-red-500">*</span>
              </Label>
              <Select
                value={formData.category}
                onValueChange={(value) => handleChange("category", value)}
                disabled={isLoading}
              >
                <SelectTrigger id="category">
                  <SelectValue placeholder="Chọn loại phản hồi" />
                </SelectTrigger>
                <SelectContent>
                  {Object.entries(categoryLabels).map(([key, label]) => (
                    <SelectItem key={key} value={key}>
                      <div className="flex flex-col items-start">
                        <span className="font-medium">{label}</span>
                        <span className="text-xs text-gray-500">
                          {categoryDescriptions[key as ComplaintCategory]}
                        </span>
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Content */}
            <div className="space-y-2">
              <Label htmlFor="content">
                Nội dung <span className="text-red-500">*</span>
              </Label>
              <Textarea
                id="content"
                placeholder="Vui lòng mô tả chi tiết vấn đề hoặc ý kiến của bạn..."
                value={formData.content}
                onChange={(e) => handleChange("content", e.target.value)}
                disabled={isLoading}
                required
                rows={6}
                className="resize-none"
              />
              <p className="text-xs text-gray-500">
                Tối thiểu 10 ký tự. Hãy mô tả chi tiết để chúng tôi hiểu rõ vấn đề của bạn.
              </p>
            </div>

            {/* Info Alert */}
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Lưu ý</AlertTitle>
              <AlertDescription className="text-sm">
                <ul className="list-disc list-inside space-y-1 mt-2">
                  <li>Thông tin của bạn sẽ được bảo mật tuyệt đối</li>
                  <li>Chúng tôi sẽ phản hồi trong vòng 24-48 giờ</li>
                  <li>Vui lòng cung cấp thông tin chính xác để được hỗ trợ nhanh chóng</li>
                </ul>
              </AlertDescription>
            </Alert>
          </CardContent>

          <CardFooter className="flex gap-3">
            <Button
              type="submit"
              disabled={isLoading || !formData.phone.trim() || !formData.content.trim()}
              className="bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)] flex-1"
            >
              {isLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Đang gửi...
                </>
              ) : (
                "Gửi phản hồi"
              )}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setFormData((prev) => ({
                  phone: prev.phone,
                  category: "OTHER",
                  content: "",
                }));
                setIsSuccess(false);
              }}
              disabled={isLoading}
            >
              Làm mới
            </Button>
          </CardFooter>
        </form>
      </Card>

      {/* Help Section */}
      <div className="mt-8 p-6 bg-gray-50 rounded-lg">
        <h3 className="font-semibold mb-3">Cần hỗ trợ khẩn cấp?</h3>
        <div className="space-y-2 text-sm text-gray-600">
          <p>📞 Hotline: <span className="font-semibold text-[hsl(var(--primary))]">1900 xxxx</span></p>
          <p>📧 Email: <span className="font-semibold text-[hsl(var(--primary))]">support@example.com</span></p>
          <p>⏰ Thời gian làm việc: 8:00 - 22:00 (Tất cả các ngày trong tuần)</p>
        </div>
      </div>
    </div>
  );
}
