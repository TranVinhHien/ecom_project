"use client"

import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Star, ChevronDown, ChevronUp } from 'lucide-react';
import { useTranslations } from 'next-intl';
import apiClient from '@/lib/apiClient';
import API from '@/assets/configs/api';
import { Loading } from '@/components/ui/loading';
import { getImageUrl } from '@/assets/helpers/convert_tool';
import ImageGalleryModal from '@/components/ImageGalleryModal';



interface Comment {
  comment_id: string;
  order_item_id: string;
  product_id: string;
  sku_id: string;
  user_id: string;
  sku_name_snapshot: string;
  rating: number;
  title: string;
  content: string;
  media: string[];
  parent_id: string | null;
  created_at: string;
  updated_at: string;
}

interface CommentsResponse {
  code: number;
  message: string;
  status: string;
  result: {
    currentPage: number;
    data: Comment[];
    limit: number;
    totalElements: number;
    totalPages: number;
  };
}

interface ProductCommentsProps {
  productId: string;
}

export default function ProductComments({ productId }: ProductCommentsProps) {
  const t = useTranslations("System");
  const [currentPage, setCurrentPage] = useState(1);
  const [allComments, setAllComments] = useState<Comment[]>([]);
  const [isExpanded, setIsExpanded] = useState(false);
  const [hasInitialLoad, setHasInitialLoad] = useState(false);
  const [isGalleryOpen, setIsGalleryOpen] = useState(false);
  const [galleryImages, setGalleryImages] = useState<string[]>([]);
  const [galleryNotes, setGalleryNotes] = useState<string[]>([]);
  const [initialGalleryIndex, setInitialGalleryIndex] = useState(0);
  const pageSize = 5;

  // Fetch comments
  const { data, isLoading, error, refetch } = useQuery<CommentsResponse>({
    queryKey: ['product-comments', productId, currentPage],
    queryFn: async () => {
      const response = await apiClient.get(
        `/comments?product_id=${productId}&page=${currentPage}&page_size=${pageSize}`,
        {
          customBaseURL: process.env.NEXT_PUBLIC_API_GATEWAY_URL || API.base_order
        }
      );
      return response.data;
    },
    enabled: isExpanded, // Fetch when expanded
  });
  console.log("Comments Data:", data,productId);
  // Update allComments when data changes
  React.useEffect(() => {
    if (data?.result?.data) {
      if (currentPage === 1 && !hasInitialLoad) {
        setAllComments(data.result.data);
        setHasInitialLoad(true);
      } else if (currentPage > 1) {
        // Only append if it's a new page load
        setAllComments(prev => {
          const existingIds = new Set(prev.map(c => c.comment_id));
          const newComments = data.result.data.filter(c => !existingIds.has(c.comment_id));
          return [...prev, ...newComments];
        });
      }
    }
  }, [data, currentPage, hasInitialLoad]);

  const totalComments = data?.result?.totalElements || 0;
  const totalPages = data?.result?.totalPages || 0;
  const hasMoreComments = currentPage < totalPages;

  // Load more comments
  const handleLoadMore = () => {
    setCurrentPage(prev => prev + 1);
  };

  // Collapse comments to show only first 5
  const handleCollapse = () => {
    setAllComments(prev => prev.slice(0, pageSize));
    setCurrentPage(1);
  };

  // Toggle expand/collapse
  const handleToggle = () => {
    if (!isExpanded) {
      setIsExpanded(true);
    } else {
      setIsExpanded(false);
      // Reset to first page when closing
      setCurrentPage(1);
    }
  };

  // Handle image click - collect all images from all comments
  const handleImageClick = (clickedImage: string, comment: Comment) => {
    // Collect all images from all comments with their notes
    const allImages: string[] = [];
    const allNotes: string[] = [];
    
    allComments.forEach((c) => {
      if (c.media && c.media.length > 0) {
        c.media.forEach((img) => {
          allImages.push(img);
          // Create note with user info and comment
          const note = `Bình luận của User ${c.user_id}${c.content ? ': ' + c.content : ''}`;
          allNotes.push(note);
        });
      }
    });

    // Find the index of the clicked image
    const clickedIndex = allImages.findIndex(img => img === clickedImage);
    
    setGalleryImages(allImages);
    setGalleryNotes(allNotes);
    setInitialGalleryIndex(clickedIndex >= 0 ? clickedIndex : 0);
    setIsGalleryOpen(true);
  };

  // Render star rating
  const renderStars = (rating: number) => {
    return (
      <div className="flex gap-1">
        {[1, 2, 3, 4, 5].map((star) => (
          <Star
            key={star}
            className={`w-4 h-4 ${
              star <= rating
                ? 'fill-yellow-400 text-yellow-400'
                : 'text-gray-300'
            }`}
          />
        ))}
      </div>
    );
  };

  // Format date
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('vi-VN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return (
    <div className="border-t pt-6 mt-6">
      {/* Header with toggle */}
      <Button
        variant="ghost"
        className="w-full flex items-center justify-between p-4 hover:bg-gray-50"
        onClick={handleToggle}
      >
        <h2 className="text-xl font-bold">
          {t("danh_gia_san_pham")} {totalComments > 0 && `(${totalComments})`}
        </h2>
        {isExpanded ? (
          <ChevronUp className="w-5 h-5" />
        ) : (
          <ChevronDown className="w-5 h-5" />
        )}
      </Button>

      {/* Comments List */}
      {isExpanded && (
        <div className="mt-4 space-y-4">
          {isLoading && !hasInitialLoad ? (
            <div className="flex justify-center py-8">
              <Loading size="lg" variant="primary" />
            </div>
          ) : error ? (
            <div className="text-center py-8 text-red-500">
              {t("co_loi_xay_ra_khi_tai_du_lieu")}
            </div>
          ) : allComments.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              {t("chua_co_danh_gia_nao")}
            </div>
          ) : (
            <>
              {allComments.map((comment) => (
                <Card key={comment.comment_id} className="p-4">
                  <div className="flex items-start gap-4">
                    {/* User Avatar Placeholder */}
                    <div className="w-10 h-10 rounded-full bg-primary text-white flex items-center justify-center font-semibold">
                      {comment.user_id.slice(0, 2).toUpperCase()}
                    </div>

                    <div className="flex-1">
                      {/* User ID & Date */}
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-semibold text-sm">
                          User {comment.user_id}
                        </span>
                        <span className="text-xs text-gray-500">
                          {formatDate(comment.created_at)}
                        </span>
                      </div>

                      {/* Rating */}
                      <div className="mb-2">
                        {renderStars(comment.rating)}
                      </div>

                      {/* SKU Info */}
                      {comment.sku_name_snapshot && (
                        <div className="text-xs text-gray-600 mb-2">
                          {t("phan_loai")}: {comment.sku_name_snapshot}
                        </div>
                      )}

                      {/* Title */}
                      {comment.title && (
                        <div className="font-semibold mb-1">
                          {comment.title}
                        </div>
                      )}

                      {/* Content */}
                      {comment.content && (
                        <div className="text-gray-700 mb-2">
                          {comment.content}
                        </div>
                      )}

                      {/* Media */}
                      {comment.media && comment.media.length > 0 && (
                        <div className="flex gap-2 flex-wrap">
                          {comment.media.map((mediaUrl, idx) => (
                            <button
                              key={idx}
                              onClick={() => handleImageClick(mediaUrl, comment)}
                              className="relative group cursor-pointer"
                            >
                              <img
                                src={getImageUrl(mediaUrl)}
                                alt={`Review ${idx + 1}`}
                                className="w-20 h-20 object-cover rounded border transition-transform group-hover:scale-105"
                                onError={(e) => {
                                  (e.target as HTMLImageElement).src = '/placeholder.png';
                                }}
                              />
                              {/* Hover overlay */}
                              <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-all rounded flex items-center justify-center">
                                <svg 
                                  className="w-6 h-6 text-white opacity-0 group-hover:opacity-100 transition-opacity"
                                  fill="none" 
                                  stroke="currentColor" 
                                  viewBox="0 0 24 24"
                                >
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                                </svg>
                              </div>
                            </button>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                </Card>
              ))}

              {/* Action Buttons */}
              <div className="flex gap-3 justify-center pt-4">
                {hasMoreComments && (
                  <Button
                    variant="outline"
                    onClick={handleLoadMore}
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <Loading size="sm" variant="primary" />
                    ) : (
                      <>
                        <ChevronDown className="w-4 h-4 mr-2" />
                        {t("xem_them")} ({totalComments - allComments.length} {t("con_lai")})
                      </>
                    )}
                  </Button>
                )}

                {allComments.length > pageSize && (
                  <Button
                    variant="outline"
                    onClick={handleCollapse}
                  >
                    <ChevronUp className="w-4 h-4 mr-2" />
                    {t("an_bot")}
                  </Button>
                )}
              </div>

              {/* Collapse All Button at Bottom */}
              <div className="flex justify-center pt-4 pb-2 border-t mt-4">
                <Button
                  variant="outline"
                  onClick={handleToggle}
                  className="gap-2"
                >
                  <ChevronUp className="w-4 h-4" />
                  {t("dong_lai")}
                </Button>
              </div>
            </>
          )}
        </div>
      )}

      {/* Image Gallery Modal */}
      <ImageGalleryModal
        images={galleryImages}
        initialIndex={initialGalleryIndex}
        isOpen={isGalleryOpen}
        onClose={() => setIsGalleryOpen(false)}
        notes={galleryNotes}
      />
    </div>
  );
}
