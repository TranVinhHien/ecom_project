package services

import "time"

// CreateCommentRequest represents the request to create a new comment/review
type CreateCommentRequest struct {
	OrderItemID string  `json:"order_item_id" binding:"required"`
	Comment     string  `json:"comment" binding:"required"`
	Star        int     `json:"star" binding:"required,min=1,max=5"`
	Title       *string `json:"title"`
	ParentID    *string `json:"parent_id"` // Nếu có parent_id thì đây là reply
}

// ListCommentsRequest represents the request to list comments for a product
type ListCommentsRequest struct {
	ProductID string `form:"product_id" binding:"required"`
	PageSize  int32  `form:"page_size" binding:"max=100"`
	Page      int32  `form:"page" binding:"min=0"`
}

// CheckReviewedItemsRequest represents the request to check reviewed order items
type CheckReviewedItemsRequest struct {
	OrderItemIDs []string `json:"order_item_ids" binding:"required,min=1"`
}

// CheckReviewedItemsResponse represents the response with reviewed order item IDs
type CheckReviewedItemsResponse struct {
	ReviewedOrderItemIDs []string `json:"reviewed_order_item_ids"`
}

// CommentResponse represents a single comment with nested replies
type CommentResponse struct {
	CommentID       string            `json:"comment_id"`
	OrderItemID     string            `json:"order_item_id"`
	ProductID       string            `json:"product_id"`
	SkuID           string            `json:"sku_id"`
	UserID          string            `json:"user_id"`
	SkuNameSnapshot *string           `json:"sku_name_snapshot"`
	Rating          int32             `json:"rating"`
	Title           *string           `json:"title"`
	Content         string            `json:"content"`
	Media           interface{}       `json:"media"`
	ParentID        *string           `json:"parent_id"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Children        []CommentResponse `json:"children,omitempty"` // Nested replies
}

// ProductRatingStats represents rating statistics for a product
type ProductRatingStats struct {
	TotalReviews  int64   `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
}

// GetBulkProductRatingStatsRequest represents the request to get rating stats for multiple products
type GetBulkProductRatingStatsRequest struct {
	ProductIDs []string `json:"product_ids" binding:"required,min=1"`
}

// ProductRatingStatsItem represents rating statistics for a single product
type ProductRatingStatsItem struct {
	ProductID     string  `json:"product_id"`
	TotalReviews  int64   `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
}

// GetBulkProductRatingStatsResponse represents the response with rating stats keyed by product_id
type GetBulkProductRatingStatsResponse struct {
	Stats map[string]ProductRatingStatsItem `json:"stats"`
}
