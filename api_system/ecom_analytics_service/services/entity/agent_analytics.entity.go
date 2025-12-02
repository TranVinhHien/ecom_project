package services

// =================================================================
// 1. DASHBOARD STATS ENTITIES
// =================================================================

// DashboardStatsResponse - Response thống kê tổng quan dashboard
type DashboardStatsResponse struct {
	TotalSessions      int64 `json:"total_sessions"`
	TotalUserMessages  int64 `json:"total_user_messages"`
	TotalAgentMessages int64 `json:"total_agent_messages"`
}

// =================================================================
// 2. MESSAGE VOLUME BY HOUR ENTITIES
// =================================================================

// MessageVolumeByHourItem - Một điểm dữ liệu theo giờ
type MessageVolumeByHourItem struct {
	HourOfDay    int32 `json:"hour_of_day"` // 0-23
	MessageCount int64 `json:"message_count"`
}

// =================================================================
// 3. TOP ACTIVE USERS ENTITIES
// =================================================================

// TopActiveUserItem - Một user tích cực
type TopActiveUserItem struct {
	UserID       string `json:"user_id"`
	MessageCount int64  `json:"message_count"`
}

// GetTopActiveUsersRequest - Request lấy top users
type GetTopActiveUsersRequest struct {
	StartDate *string `form:"start_date,omitempty"` // YYYY-MM-DD
	EndDate   *string `form:"end_date,omitempty"`   // YYYY-MM-DD
	Limit     int32   `form:"limit,default=10"`     // Default top 10
}

// =================================================================
// 4. TOPIC STATS ENTITIES
// =================================================================

// TopicStatsItem - Thống kê một chủ đề
type TopicStatsItem struct {
	Topic      string  `json:"topic"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// =================================================================
// 5. PURCHASE INTENT STATS ENTITIES
// =================================================================

// PurchaseIntentStatsItem - Thống kê ý định mua hàng
type PurchaseIntentStatsItem struct {
	PurchaseIntent string `json:"purchase_intent"` // High, Medium, Low
	Count          int64  `json:"count"`
}

// =================================================================
// 6. TOP MENTIONED CATEGORIES ENTITIES
// =================================================================

// CustomMetadata - Cấu trúc JSON của custom_metadata
type CustomMetadata struct {
	Entities []EntityItem `json:"entities"`
}

// EntityItem - Một entity trong metadata
type EntityItem struct {
	Type    string `json:"type"`    // category_name, product_name, etc.
	Context string `json:"context"` // Tên danh mục/sản phẩm
}

// CategoryMentionItem - Một danh mục được nhắc đến
type CategoryMentionItem struct {
	CategoryName string `json:"category_name"`
	MentionCount int    `json:"mention_count"`
}

// GetTopMentionedCategoriesRequest - Request lấy top categories
type GetTopMentionedCategoriesRequest struct {
	StartDate *string `form:"start_date,omitempty"` // YYYY-MM-DD
	EndDate   *string `form:"end_date,omitempty"`   // YYYY-MM-DD
	Limit     int     `form:"limit,default=10"`     // Default top 10
}

// =================================================================
// COMMON REQUEST FOR TIME FILTERING
// =================================================================

// TimeRangeRequest - Request chung cho các query có filter thời gian
type TimeRangeRequest struct {
	StartDate *string `form:"start_date,omitempty"` // YYYY-MM-DD
	EndDate   *string `form:"end_date,omitempty"`   // YYYY-MM-DD
}
