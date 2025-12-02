package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
	assets_services "github.com/TranVinhHien/ecom_order_service/services/assets"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"
	"github.com/google/uuid"
)

// Helper function để xử lý media NULL từ database
func getMediaOrEmpty(media sql.NullString) interface{} {
	if !media.Valid || len(media.String) == 0 {
		return []interface{}{} // Return empty array thay vì null
	}
	var result interface{}
	if err := json.Unmarshal([]byte(media.String), &result); err != nil {
		return []interface{}{} // Nếu unmarshal lỗi, return empty array
	}
	return result
}

// CreateComment xử lý việc tạo bình luận/đánh giá sản phẩm
func (s *service) CreateComment(ctx context.Context, userID string, req services.CreateCommentRequest) *assets_services.ServiceError {
	// Kiểm tra nếu có parent_id (reply comment)
	if req.ParentID != nil && *req.ParentID != "" {
		// TODO: Xử lý logic reply comment ở đây (tạm thời để None)
		// Có thể cần kiểm tra parent comment có tồn tại không, v.v.
		return assets_services.NewError(
			http.StatusNotImplemented,
			errors.New("reply comment feature is not implemented yet"),
		)
	}

	// Bước 1: Check quyền review - Kiểm tra user có quyền review order_item này không
	permission, err := s.repository.CheckReviewPermission(ctx, db.CheckReviewPermissionParams{
		ID:     req.OrderItemID,
		UserID: userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return assets_services.NewError(
				http.StatusForbidden,
				errors.New("bạn không có quyền đánh giá sản phẩm này. Chỉ có thể đánh giá sau khi đơn hàng hoàn thành"),
			)
		}
		return assets_services.NewError(
			http.StatusInternalServerError,
			fmt.Errorf("lỗi khi kiểm tra quyền đánh giá: %w", err),
		)
	}

	// Bước 2: Kiểm tra xem order_item này đã được review chưa
	_, err = s.repository.GetCommentByOrderItemID(ctx, req.OrderItemID)
	if err == nil {
		// Nếu không có lỗi, nghĩa là đã tồn tại review
		return assets_services.NewError(
			http.StatusConflict,
			errors.New("bạn đã đánh giá sản phẩm này rồi. Mỗi sản phẩm chỉ được đánh giá 1 lần"),
		)
	}
	if err != sql.ErrNoRows {
		// Nếu lỗi khác ErrNoRows, có nghĩa là lỗi DB
		return assets_services.NewError(
			http.StatusInternalServerError,
			fmt.Errorf("lỗi khi kiểm tra đánh giá hiện có: %w", err),
		)
	}

	// Bước 3: Tạo comment mới
	commentID := uuid.New().String()

	// Lấy sku_name_snapshot (có thể cần query thêm nếu cần, hiện tại set NULL)
	// var skuNameSnapshot sql.NullString

	params := db.CreateCommentParams{
		CommentID:       commentID,
		OrderItemID:     req.OrderItemID,
		ProductID:       permission.ProductID,
		SkuID:           permission.SkuID,
		UserID:          userID,
		SkuNameSnapshot: permission.SkuAttributesSnapshot,
		Rating:          int8(req.Star),
		Title: sql.NullString{
			String: "",
			Valid:  false,
		},
		Content: sql.NullString{
			String: req.Comment,
			Valid:  true,
		},
		Media: sql.NullString{},
		ParentID: sql.NullString{
			String: "",
			Valid:  false,
		},
	}

	// Set title nếu có
	if req.Title != nil && *req.Title != "" {
		params.Title = sql.NullString{
			String: *req.Title,
			Valid:  true,
		}
	}

	// Thực hiện insert vào DB
	err = s.repository.CreateComment(ctx, params)
	if err != nil {
		return assets_services.NewError(
			http.StatusInternalServerError,
			fmt.Errorf("lỗi khi tạo đánh giá: %w", err),
		)
	}

	return nil
}

// ListComments lấy danh sách bình luận cho một sản phẩm (có phân trang)
// Các comment con (replies) sẽ được nest vào comment cha
func (s *service) ListComments(ctx context.Context, req services.ListCommentsRequest) (map[string]interface{}, *assets_services.ServiceError) {
	// Bước 1: Lấy danh sách comment gốc (parent_id IS NULL)
	comments, err := s.repository.ListCommentsByProduct(ctx, db.ListCommentsByProductParams{
		ProductID: req.ProductID,
		Limit:     req.PageSize,
		Offset:    int32(req.PageSize * (req.Page - 1)),
	})
	if err != nil {
		return nil, assets_services.NewError(
			http.StatusInternalServerError,
			fmt.Errorf("lỗi khi lấy danh sách bình luận: %w", err),
		)
	}

	// Bước 2: Lấy thống kê rating
	stats, err := s.repository.GetProductRatingStats(ctx, req.ProductID)
	if err != nil && err != sql.ErrNoRows {
		return nil, assets_services.NewError(
			http.StatusInternalServerError,
			fmt.Errorf("lỗi khi lấy thống kê đánh giá: %w", err),
		)
	}

	// Bước 3: Build response với nested children
	var commentResponses []services.CommentResponse
	for _, comment := range comments {
		commentResp := services.CommentResponse{
			CommentID:       comment.CommentID,
			OrderItemID:     comment.OrderItemID,
			ProductID:       comment.ProductID,
			SkuID:           comment.SkuID,
			UserID:          comment.UserID,
			SkuNameSnapshot: nil,
			Rating:          int32(comment.Rating),
			Title:           nil,
			Content:         comment.Content.String,
			Media:           getMediaOrEmpty(comment.Media),
			ParentID:        nil,
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
			Children:        []services.CommentResponse{},
		}

		// Set optional fields
		if comment.SkuNameSnapshot.Valid {
			commentResp.SkuNameSnapshot = &comment.SkuNameSnapshot.String
		}
		if comment.Title.Valid {
			commentResp.Title = &comment.Title.String
		}
		if comment.ParentID.Valid {
			commentResp.ParentID = &comment.ParentID.String
		}

		// Lấy các replies (children) cho comment này
		replies, err := s.repository.GetRepliesByCommentID(ctx, sql.NullString{
			String: comment.CommentID,
			Valid:  true,
		})
		if err != nil && err != sql.ErrNoRows {
			// Log error nhưng không fail toàn bộ request
			fmt.Printf("Warning: Lỗi khi lấy replies cho comment %s: %v\n", comment.CommentID, err)
		}

		// Map replies vào children
		if len(replies) > 0 {
			for _, reply := range replies {
				replyResp := services.CommentResponse{
					CommentID:       reply.CommentID,
					OrderItemID:     reply.OrderItemID,
					ProductID:       reply.ProductID,
					SkuID:           reply.SkuID,
					UserID:          reply.UserID,
					SkuNameSnapshot: nil,
					Rating:          int32(reply.Rating),
					Title:           nil,
					Content:         reply.Content.String,
					Media:           getMediaOrEmpty(reply.Media),
					ParentID:        nil,
					CreatedAt:       reply.CreatedAt,
					UpdatedAt:       reply.UpdatedAt,
				}

				// Set optional fields cho reply
				if reply.SkuNameSnapshot.Valid {
					replyResp.SkuNameSnapshot = &reply.SkuNameSnapshot.String
				}
				if reply.Title.Valid {
					replyResp.Title = &reply.Title.String
				}
				if reply.ParentID.Valid {
					replyResp.ParentID = &reply.ParentID.String
				}

				commentResp.Children = append(commentResp.Children, replyResp)
			}
		}

		commentResponses = append(commentResponses, commentResp)
	}

	// Bước 4: Parse average rating từ interface{}
	averageRating := 0.0
	if stats.AverageRating != nil {
		// MySQL AVG() trả về []uint8 hoặc nil
		switch v := stats.AverageRating.(type) {
		case []uint8:
			// Convert []uint8 to string then to float64
			fmt.Sscanf(string(v), "%f", &averageRating)
		case float64:
			averageRating = v
		case float32:
			averageRating = float64(v)
		}
	}

	// Bước 5: Build final response
	// result := map[string]interface{}{
	// 	"data": commentResponses,
	// 	"stats": map[string]interface{}{
	// 		"total_reviews":  stats.TotalReviews,
	// 		"average_rating": averageRating,
	// 	},
	// 	"pagination": map[string]interface{}{
	// 		"limit":  req.Limit,
	// 		"offset": req.Offset,
	// 	},
	// }

	result := map[string]interface{}{}
	result["data"] = commentResponses
	result["currentPage"] = req.Page
	result["totalPages"] = (stats.TotalReviews + int64(req.PageSize) - 1) / int64(req.PageSize)
	result["totalElements"] = stats.TotalReviews
	result["limit"] = req.PageSize

	return result, nil
}

// CheckReviewedItems kiểm tra danh sách order_item_id nào đã được review
// Chỉ trả về những order_item_id đã có bình luận
func (s *service) CheckReviewedItems(ctx context.Context, req services.CheckReviewedItemsRequest) (*services.CheckReviewedItemsResponse, *assets_services.ServiceError) {
	if len(req.OrderItemIDs) == 0 {
		return &services.CheckReviewedItemsResponse{
			ReviewedOrderItemIDs: []string{},
		}, nil
	}

	// Gọi query để check bulk
	reviewedItems, err := s.repository.CheckBulkOrderItemsReviewed(ctx, req.OrderItemIDs)
	if err != nil {
		return nil, assets_services.NewError(
			http.StatusInternalServerError,
			fmt.Errorf("lỗi khi kiểm tra các đánh giá: %w", err),
		)
	}

	return &services.CheckReviewedItemsResponse{
		ReviewedOrderItemIDs: reviewedItems,
	}, nil
}

// GetBulkProductRatingStats lấy thống kê đánh giá cho nhiều sản phẩm cùng lúc
// Trả về map với product_id là key và thông tin rating là value
func (s *service) GetBulkProductRatingStats(ctx context.Context, req services.GetBulkProductRatingStatsRequest) (map[string]interface{}, *assets_services.ServiceError) {
	if len(req.ProductIDs) == 0 {
		return map[string]interface{}{}, nil

	}

	// Gọi query để lấy stats cho nhiều products
	stats, err := s.repository.GetBulkProductRatingStats(ctx, req.ProductIDs)
	if err != nil {
		return nil, assets_services.NewError(
			http.StatusInternalServerError,
			fmt.Errorf("lỗi khi lấy thống kê đánh giá sản phẩm: %w", err),
		)
	}

	// Build map response
	statsMap := []services.ProductRatingStatsItem{}
	for _, stat := range stats {
		// Parse average rating từ interface{}
		averageRating := 0.0
		if stat.AverageRating != nil {
			switch v := stat.AverageRating.(type) {
			case []uint8:
				// MySQL AVG() trả về []uint8
				fmt.Sscanf(string(v), "%f", &averageRating)
			case float64:
				averageRating = v
			case float32:
				averageRating = float64(v)
			}
		}

		statsMap = append(statsMap, services.ProductRatingStatsItem{
			ProductID:     stat.ProductID,
			TotalReviews:  stat.TotalReviews,
			AverageRating: averageRating,
		})
	}
	result := map[string]interface{}{}
	result["data"] = statsMap
	return result, nil
}
