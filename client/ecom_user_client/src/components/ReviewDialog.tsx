"use client";

import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Star } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import apiClient from "@/lib/apiClient";
import API from "@/assets/configs/api";
import { Loading } from "@/components/ui/loading";

interface ReviewDialogProps {
  isOpen: boolean;
  onClose: () => void;
  orderItemId: string;
  productName: string;
  productImage?: string;
  onSuccess?: () => void;
}

const RATING_TITLES = [
  { star: 1, title: "Rất tệ", color: "text-red-600" },
  { star: 2, title: "Tệ", color: "text-orange-600" },
  { star: 3, title: "Bình thường", color: "text-yellow-600" },
  { star: 4, title: "Tốt", color: "text-blue-600" },
  { star: 5, title: "Cực kì hài lòng", color: "text-green-600" },
];

export default function ReviewDialog({
  isOpen,
  onClose,
  orderItemId,
  productName,
  productImage,
  onSuccess,
}: ReviewDialogProps) {
  const [rating, setRating] = useState<number>(5);
  const [hoveredRating, setHoveredRating] = useState<number>(0);
  const [comment, setComment] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { toast } = useToast();

  const getImageUrl = (imagePath: string | null | undefined) => {
    if (!imagePath) return '/placeholder.png';
    if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
      return imagePath;
    }
    return `http://${imagePath}`;
  };

  const selectedRatingInfo = RATING_TITLES.find(r => r.star === rating) || RATING_TITLES[4];

  const handleSubmit = async () => {
    if (rating === 0) {
      toast({
        title: "Lỗi",
        description: "Vui lòng chọn số sao đánh giá",
        variant: "destructive",
      });
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await apiClient.post(
        "/comments",
        {
          order_item_id: orderItemId,
          comment: comment.trim() || "",
          star: rating,
          title: selectedRatingInfo.title,
        },
        {
          customBaseURL: process.env.NEXT_PUBLIC_API_GATEWAY_URL || API.base_order,
        }
      );

      if (response.data.code === 200 || response.data.status === "success") {
        toast({
          title: "Thành công",
          description: "Đánh giá của bạn đã được gửi thành công!",
        });
        
        // Reset form
        setRating(5);
        setComment("");
        
        // Call success callback
        if (onSuccess) {
          onSuccess();
        }
        
        // Close dialog
        onClose();
      } else {
        throw new Error(response.data.message || "Không thể gửi đánh giá");
      }
    } catch (error: any) {
      console.error("Submit review error:", error);
      toast({
        title: "Lỗi",
        description: error.response.data.message || "Không thể gửi đánh giá. Vui lòng thử lại.",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    if (!isSubmitting) {
      setRating(5);
      setComment("");
      onClose();
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle className="text-2xl font-bold text-center">
            Đánh giá sản phẩm
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-6 py-4">
          {/* Product Info */}
          <div className="flex items-center gap-4 p-4 bg-gray-50 rounded-lg">
            <img
              src={getImageUrl(productImage)}
              alt={productName}
              className="w-16 h-16 object-cover rounded-lg border"
            />
            <div className="flex-1">
              <h3 className="font-medium text-gray-800 line-clamp-2">
                {productName}
              </h3>
            </div>
          </div>

          {/* Rating Stars */}
          <div className="text-center space-y-4">
            <div className="flex justify-center items-center gap-2">
              {[1, 2, 3, 4, 5].map((star) => (
                <button
                  key={star}
                  type="button"
                  onClick={() => setRating(star)}
                  onMouseEnter={() => setHoveredRating(star)}
                  onMouseLeave={() => setHoveredRating(0)}
                  className="transition-transform hover:scale-110 focus:outline-none"
                >
                  <Star
                    className={`w-12 h-12 transition-colors ${
                      star <= (hoveredRating || rating)
                        ? "fill-yellow-400 text-yellow-400"
                        : "text-gray-300"
                    }`}
                  />
                </button>
              ))}
            </div>
            
            {/* Rating Title */}
            <div className={`text-2xl font-bold ${selectedRatingInfo.color}`}>
              {selectedRatingInfo.title}
            </div>
          </div>

          {/* Comment Textarea */}
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">
              Nhận xét của bạn (tùy chọn)
            </label>
            <Textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Chia sẻ trải nghiệm của bạn về sản phẩm này..."
              className="min-h-[120px] resize-none"
              disabled={isSubmitting}
              maxLength={1000}
            />
            <div className="text-xs text-gray-500 text-right">
              {comment.length}/1000 ký tự
            </div>
          </div>
        </div>

        <DialogFooter className="gap-2">
          <Button
            variant="outline"
            onClick={handleClose}
            disabled={isSubmitting}
          >
            Hủy
          </Button>
          <Button
            onClick={handleSubmit}
            disabled={isSubmitting || rating === 0}
            className="bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/0.9)]"
          >
            {isSubmitting ? (
              <>
                <Loading size="sm" variant="default" />
                <span className="ml-2">Đang gửi...</span>
              </>
            ) : (
              "Gửi đánh giá"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
