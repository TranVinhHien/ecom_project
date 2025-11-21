package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	db "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/agent_ai_db"
	assets_services "github.com/TranVinhHien/ecom_analytics_service/services/assets"
	entity "github.com/TranVinhHien/ecom_analytics_service/services/entity"
)

// =================================================================
// HELPER FUNCTIONS
// =================================================================

// parseDate converts string date (YYYY-MM-DD) to time.Time
func parseDateAgent(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// convertToInt64 safely converts interface{} to int64
func convertToInt64(val interface{}) int64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}

// =================================================================
// 1. DASHBOARD STATS
// =================================================================

func (s *service) GetDashboardStats(ctx context.Context, startDate, endDate *string) (*entity.DashboardStatsResponse, *assets_services.ServiceError) {
	// Parse dates to time.Time
	var startTimeNull, endTimeNull sql.NullTime
	if startDate != nil && *startDate != "" {
		if t, err := parseDateAgent(*startDate); err == nil {
			startTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if endDate != nil && *endDate != "" {
		if t, err := parseDateAgent(*endDate); err == nil {
			endTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	stats, err := s.db_agent_ai_db.GetDashboardStats(ctx, db.GetDashboardStatsParams{
		StartTime: startTimeNull,
		EndTime:   endTimeNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy thống kê dashboard: %w", err))
	}

	return &entity.DashboardStatsResponse{
		TotalSessions:      stats.TotalSessions,
		TotalUserMessages:  convertToInt64(stats.TotalUserMessages),
		TotalAgentMessages: convertToInt64(stats.TotalAgentMessages),
	}, nil
}

// =================================================================
// 2. MESSAGE VOLUME BY HOUR
// =================================================================

func (s *service) GetMessageVolumeByHour(ctx context.Context, startDate, endDate *string) ([]entity.MessageVolumeByHourItem, *assets_services.ServiceError) {
	// Parse dates to time.Time
	var startTimeNull, endTimeNull sql.NullTime
	if startDate != nil && *startDate != "" {
		if t, err := parseDateAgent(*startDate); err == nil {
			startTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if endDate != nil && *endDate != "" {
		if t, err := parseDateAgent(*endDate); err == nil {
			endTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	rows, err := s.db_agent_ai_db.GetMessageVolumeByHour(ctx, db.GetMessageVolumeByHourParams{
		StartTime: startTimeNull,
		EndTime:   endTimeNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy dữ liệu message volume: %w", err))
	}

	result := make([]entity.MessageVolumeByHourItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.MessageVolumeByHourItem{
			HourOfDay:    row.HourOfDay,
			MessageCount: row.MessageCount,
		})
	}

	return result, nil
}

// =================================================================
// 3. TOP ACTIVE USERS
// =================================================================

func (s *service) GetTopActiveUsers(ctx context.Context, req *entity.GetTopActiveUsersRequest) ([]entity.TopActiveUserItem, *assets_services.ServiceError) {
	// Validate and set default limit
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	// Parse dates to time.Time
	var startTimeNull, endTimeNull sql.NullTime
	if req.StartDate != nil && *req.StartDate != "" {
		if t, err := parseDateAgent(*req.StartDate); err == nil {
			startTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if t, err := parseDateAgent(*req.EndDate); err == nil {
			endTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	rows, err := s.db_agent_ai_db.GetTopActiveUsers(ctx, db.GetTopActiveUsersParams{
		StartTime: startTimeNull,
		EndTime:   endTimeNull,
		Limit:     req.Limit,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy top active users: %w", err))
	}

	result := make([]entity.TopActiveUserItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.TopActiveUserItem{
			UserID:       row.UserID,
			MessageCount: row.MessageCount,
		})
	}

	return result, nil
}

// =================================================================
// 4. TOPIC STATS
// =================================================================

func (s *service) GetTopicStats(ctx context.Context, startDate, endDate *string) ([]entity.TopicStatsItem, *assets_services.ServiceError) {
	// Parse dates to time.Time
	var startTimeNull, endTimeNull sql.NullTime
	if startDate != nil && *startDate != "" {
		if t, err := parseDateAgent(*startDate); err == nil {
			startTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if endDate != nil && *endDate != "" {
		if t, err := parseDateAgent(*endDate); err == nil {
			endTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	rows, err := s.db_agent_ai_db.GetTopicStats(ctx, db.GetTopicStatsParams{
		StartTime: startTimeNull,
		EndTime:   endTimeNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy thống kê topic: %w", err))
	}

	result := make([]entity.TopicStatsItem, 0, len(rows))
	for _, row := range rows {
		var topic string
		if err := json.Unmarshal(row.Topic, &topic); err != nil {
			topic = string(row.Topic)
		}

		result = append(result, entity.TopicStatsItem{
			Topic:      topic,
			Count:      row.Count,
			Percentage: row.Percentage,
		})
	}

	return result, nil
}

// =================================================================
// 5. PURCHASE INTENT STATS
// =================================================================

func (s *service) GetPurchaseIntentStats(ctx context.Context, startDate, endDate *string) ([]entity.PurchaseIntentStatsItem, *assets_services.ServiceError) {
	// Parse dates to time.Time
	var startTimeNull, endTimeNull sql.NullTime
	if startDate != nil && *startDate != "" {
		if t, err := parseDateAgent(*startDate); err == nil {
			startTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if endDate != nil && *endDate != "" {
		if t, err := parseDateAgent(*endDate); err == nil {
			endTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	rows, err := s.db_agent_ai_db.GetPurchaseIntentStats(ctx, db.GetPurchaseIntentStatsParams{
		StartTime: startTimeNull,
		EndTime:   endTimeNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy thống kê purchase intent: %w", err))
	}

	result := make([]entity.PurchaseIntentStatsItem, 0, len(rows))
	for _, row := range rows {
		var intent string
		if err := json.Unmarshal(row.PurchaseIntent, &intent); err != nil {
			intent = string(row.PurchaseIntent)
		}

		result = append(result, entity.PurchaseIntentStatsItem{
			PurchaseIntent: intent,
			Count:          row.Count,
		})
	}

	return result, nil
}

// =================================================================
// 6. TOP MENTIONED CATEGORIES
// =================================================================

func (s *service) GetTopMentionedCategories(ctx context.Context, req *entity.GetTopMentionedCategoriesRequest) ([]entity.CategoryMentionItem, *assets_services.ServiceError) {
	// Validate and set default limit
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	// Parse dates to time.Time
	var startTimeNull, endTimeNull sql.NullTime
	if req.StartDate != nil && *req.StartDate != "" {
		if t, err := parseDateAgent(*req.StartDate); err == nil {
			startTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if t, err := parseDateAgent(*req.EndDate); err == nil {
			endTimeNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	// Get raw JSON strings from database
	rows, err := s.db_agent_ai_db.GetTopMentionedCategories(ctx, db.GetTopMentionedCategoriesParams{
		StartTime: startTimeNull,
		EndTime:   endTimeNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy dữ liệu categories: %w", err))
	}

	// Map to count category mentions
	categoryCount := make(map[string]int)

	// Process each JSON string
	for _, row := range rows {
		if !row.Valid || row.String == "" {
			continue
		}

		// Unmarshal JSON to struct
		var metadata entity.CustomMetadata
		if err := json.Unmarshal([]byte(row.String), &metadata); err != nil {
			// Skip invalid JSON
			continue
		}

		// Process entities array
		for _, entity := range metadata.Entities {
			// Filter only category_name type
			if entity.Type == "category_name" && entity.Context != "" {
				categoryCount[entity.Context]++
			}
		}
	}

	// Convert map to slice for sorting
	type categoryPair struct {
		name  string
		count int
	}
	pairs := make([]categoryPair, 0, len(categoryCount))
	for name, count := range categoryCount {
		pairs = append(pairs, categoryPair{name: name, count: count})
	}

	// Sort by count descending
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Limit results
	limit := req.Limit
	if limit > len(pairs) {
		limit = len(pairs)
	}

	// Build result
	result := make([]entity.CategoryMentionItem, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, entity.CategoryMentionItem{
			CategoryName: pairs[i].name,
			MentionCount: pairs[i].count,
		})
	}

	return result, nil
}
