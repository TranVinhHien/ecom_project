package services

import "time"

// =================================================================
// MESSAGE RATINGS ENTITIES
// =================================================================

// SubmitMessageRatingRequest - Request khi user đánh giá message
type SubmitMessageRatingRequest struct {
	EventID       string  `json:"event_id" binding:"required"`          // ID của event cần đánh giá
	SessionID     string  `json:"session_id" binding:"required"`        // Session ID của hội thoại
	UserID        *string `json:"user_id,omitempty"`                    // User ID (nếu đã đăng nhập)
	Rating        int8    `json:"rating" binding:"required,oneof=1 -1"` // 1 = Like, -1 = Dislike
	UserPrompt    *string `json:"user_prompt,omitempty"`                // Snapshot câu hỏi
	AgentResponse *string `json:"agent_response,omitempty"`             // Snapshot câu trả lời
}

// MessageRatingStatsResponse - Response thống kê tổng quan ratings
type MessageRatingStatsResponse struct {
	TotalRatings     int64   `json:"total_ratings"`
	LikeCount        int64   `json:"like_count"`
	DislikeCount     int64   `json:"dislike_count"`
	SatisfactionRate float64 `json:"satisfaction_rate"` // % Like
}

// MessageRatingTimeSeriesItem - Một điểm dữ liệu theo thời gian
type MessageRatingTimeSeriesItem struct {
	ReportDate       string  `json:"report_date"` // YYYY-MM-DD
	TotalRatings     int64   `json:"total_ratings"`
	LikeCount        int64   `json:"like_count"`
	DislikeCount     int64   `json:"dislike_count"`
	SatisfactionRate float64 `json:"satisfaction_rate"`
}

// MessageRatingDetailItem - Chi tiết một rating
type MessageRatingDetailItem struct {
	ID            int64     `json:"id"`
	EventID       string    `json:"event_id"`
	SessionID     string    `json:"session_id"`
	UserID        *string   `json:"user_id,omitempty"`
	Rating        int8      `json:"rating"`
	UserPrompt    *string   `json:"user_prompt,omitempty"`
	AgentResponse *string   `json:"agent_response,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetMessageRatingsRequest - Request lấy danh sách ratings (for admin)
type GetMessageRatingsRequest struct {
	SessionID *string `form:"session_id,omitempty"`
	UserID    *string `form:"user_id,omitempty"`
	Rating    *int8   `form:"rating,omitempty"`
	StartDate *string `form:"start_date,omitempty"` // YYYY-MM-DD
	EndDate   *string `form:"end_date,omitempty"`   // YYYY-MM-DD
	Page      int     `form:"page,default=1"`
	PageSize  int     `form:"page_size,default=20"`
}

// =================================================================
// CUSTOMER FEEDBACK ENTITIES
// =================================================================

// SubmitCustomerFeedbackRequest - Request khi user gửi feedback/complaint
type SubmitCustomerFeedbackRequest struct {
	UserID   *string `json:"user_id,omitempty"`                                          // User ID (nếu đã đăng nhập)
	Email    *string `json:"email,omitempty"`                                            // Email liên hệ
	Phone    *string `json:"phone,omitempty"`                                            // SĐT liên hệ
	Category string  `json:"category" binding:"required,oneof=BUG COMPLAINT SUGGESTION"` // Phân loại
	Content  string  `json:"content" binding:"required,min=10"`                          // Nội dung chi tiết
}

// CustomerFeedbackItem - Một item feedback
type CustomerFeedbackItem struct {
	ID        string    `json:"id"`
	UserID    *string   `json:"user_id,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Category  string    `json:"category"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetCustomerFeedbacksRequest - Request lấy danh sách feedback (for admin)
type GetCustomerFeedbacksRequest struct {
	Category  *string `form:"category,omitempty"`
	UserID    *string `form:"user_id,omitempty"`
	StartDate *string `form:"start_date,omitempty"` // YYYY-MM-DD
	EndDate   *string `form:"end_date,omitempty"`   // YYYY-MM-DD
	Page      int     `form:"page,default=1"`
	PageSize  int     `form:"page_size,default=20"`
}

// CustomerFeedbackStatsResponse - Response thống kê feedback
type CustomerFeedbackStatsResponse struct {
	TotalFeedbacks  int64 `json:"total_feedbacks"`
	BugCount        int64 `json:"bug_count"`
	ComplaintCount  int64 `json:"complaint_count"`
	SuggestionCount int64 `json:"suggestion_count"`
	UniqueUsers     int64 `json:"unique_users"`
}

// CustomerFeedbackCategoryStats - Thống kê theo category
type CustomerFeedbackCategoryStats struct {
	Category      string `json:"category"`
	FeedbackCount int64  `json:"feedback_count"`
}
