package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/interact"
	assets_services "github.com/TranVinhHien/ecom_analytics_service/services/assets"
	entity "github.com/TranVinhHien/ecom_analytics_service/services/entity"

	"github.com/google/uuid"
)

// =================================================================
// HELPER FUNCTIONS
// =================================================================

// parseDate converts string date (YYYY-MM-DD) to time.Time
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// =================================================================
// MESSAGE RATINGS METHODS
// =================================================================

func (s *service) SubmitMessageRating(ctx context.Context, req *entity.SubmitMessageRatingRequest) *assets_services.ServiceError {
	// Validate rating value
	if req.Rating != 1 && req.Rating != -1 {
		return assets_services.NewError(400, fmt.Errorf("rating phải là 1 (Like) hoặc -1 (Dislike)"))
	}

	// Convert UserID to sql.NullString
	var userIDNull sql.NullString
	if req.UserID != nil && *req.UserID != "" {
		userIDNull = sql.NullString{String: *req.UserID, Valid: true}
	}

	// Convert UserPrompt to sql.NullString
	var userPromptNull sql.NullString
	if req.UserPrompt != nil && *req.UserPrompt != "" {
		userPromptNull = sql.NullString{String: *req.UserPrompt, Valid: true}
	}

	// Convert AgentResponse to sql.NullString
	var agentResponseNull sql.NullString
	if req.AgentResponse != nil && *req.AgentResponse != "" {
		agentResponseNull = sql.NullString{String: *req.AgentResponse, Valid: true}
	}

	// Create rating record
	err := s.interact.CreateMessageRating(ctx, db.CreateMessageRatingParams{
		EventID:       req.EventID,
		SessionID:     req.SessionID,
		UserID:        userIDNull,
		Rating:        int8(req.Rating),
		UserPrompt:    userPromptNull,
		AgentResponse: agentResponseNull,
	})

	if err != nil {
		return assets_services.NewError(500, fmt.Errorf("lỗi khi lưu đánh giá: %w", err))
	}

	return nil
}

func (s *service) GetMessageRatingStats(ctx context.Context, startDate, endDate *string) (*entity.MessageRatingStatsResponse, *assets_services.ServiceError) {
	// Parse dates to time.Time
	var startDateNull, endDateNull sql.NullTime
	if startDate != nil && *startDate != "" {
		if t, err := parseDate(*startDate); err == nil {
			startDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if endDate != nil && *endDate != "" {
		if t, err := parseDate(*endDate); err == nil {
			endDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	stats, err := s.interact.GetMessageRatingStats(ctx, db.GetMessageRatingStatsParams{
		StartDate: startDateNull,
		EndDate:   endDateNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy thống kê: %w", err))
	}

	return &entity.MessageRatingStatsResponse{
		TotalRatings:     stats.TotalRatings,
		LikeCount:        stats.LikeCount,
		DislikeCount:     stats.DislikeCount,
		SatisfactionRate: stats.SatisfactionRate,
	}, nil
}

func (s *service) GetMessageRatingsTimeSeries(ctx context.Context, startDate, endDate *string) ([]entity.MessageRatingTimeSeriesItem, *assets_services.ServiceError) {
	// Parse dates to time.Time
	var startDateNull, endDateNull sql.NullTime
	if startDate != nil && *startDate != "" {
		if t, err := parseDate(*startDate); err == nil {
			startDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if endDate != nil && *endDate != "" {
		if t, err := parseDate(*endDate); err == nil {
			endDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	rows, err := s.interact.GetMessageRatingsTimeSeries(ctx, db.GetMessageRatingsTimeSeriesParams{
		StartDate: startDateNull,
		EndDate:   endDateNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy dữ liệu time series: %w", err))
	}

	result := make([]entity.MessageRatingTimeSeriesItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.MessageRatingTimeSeriesItem{
			ReportDate:       row.ReportDate.Format("2006-01-02"),
			TotalRatings:     row.TotalRatings,
			LikeCount:        row.LikeCount,
			DislikeCount:     row.DislikeCount,
			SatisfactionRate: row.SatisfactionRate,
		})
	}

	return result, nil
}

func (s *service) GetMessageRatingsList(ctx context.Context, req *entity.GetMessageRatingsRequest) ([]entity.MessageRatingDetailItem, *assets_services.ServiceError) {
	// Pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	offset := int32((req.Page - 1) * req.PageSize)
	limit := int32(req.PageSize)

	// Convert filters to sql.Null types
	var sessionIDNull, userIDNull sql.NullString
	var startDateNull, endDateNull sql.NullTime
	var ratingNull sql.NullInt16

	if req.SessionID != nil && *req.SessionID != "" {
		sessionIDNull = sql.NullString{String: *req.SessionID, Valid: true}
	}
	if req.UserID != nil && *req.UserID != "" {
		userIDNull = sql.NullString{String: *req.UserID, Valid: true}
	}
	if req.Rating != nil {
		ratingNull = sql.NullInt16{Int16: int16(*req.Rating), Valid: true}
	}
	if req.StartDate != nil && *req.StartDate != "" {
		if t, err := parseDate(*req.StartDate); err == nil {
			startDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if t, err := parseDate(*req.EndDate); err == nil {
			endDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	rows, err := s.interact.GetMessageRatingsBySession(ctx, db.GetMessageRatingsBySessionParams{
		SessionIDFilter: sessionIDNull,
		UserIDFilter:    userIDNull,
		RatingFilter:    ratingNull,
		StartDate:       startDateNull,
		EndDate:         endDateNull,
		Limit:           limit,
		Offset:          offset,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy danh sách ratings: %w", err))
	}

	result := make([]entity.MessageRatingDetailItem, 0, len(rows))
	for _, row := range rows {
		item := entity.MessageRatingDetailItem{
			ID:        int64(row.ID),
			EventID:   row.EventID,
			SessionID: row.SessionID,
			Rating:    row.Rating,
			CreatedAt: row.CreatedAt,
		}

		if row.UserID.Valid {
			item.UserID = &row.UserID.String
		}
		if row.UserPrompt.Valid {
			item.UserPrompt = &row.UserPrompt.String
		}
		if row.AgentResponse.Valid {
			item.AgentResponse = &row.AgentResponse.String
		}

		result = append(result, item)
	}

	return result, nil
}

// =================================================================
// CUSTOMER FEEDBACK METHODS
// =================================================================

func (s *service) SubmitCustomerFeedback(ctx context.Context, req *entity.SubmitCustomerFeedbackRequest) (string, *assets_services.ServiceError) {
	// Validate category
	validCategories := map[string]bool{"BUG": true, "COMPLAINT": true, "SUGGESTION": true, "OTHER": true}
	if !validCategories[req.Category] {
		return "", assets_services.NewError(400, fmt.Errorf("category không hợp lệ, chỉ chấp nhận: BUG, COMPLAINT, SUGGESTION, OTHER"))
	}

	// Generate UUID
	feedbackID := uuid.New().String()

	// Convert optional fields to sql.Null types
	var userIDNull, emailNull, phoneNull sql.NullString
	if req.UserID != nil && *req.UserID != "" {
		userIDNull = sql.NullString{String: *req.UserID, Valid: true}
	}
	if req.Email != nil && *req.Email != "" {
		emailNull = sql.NullString{String: *req.Email, Valid: true}
	}
	if req.Phone != nil && *req.Phone != "" {
		phoneNull = sql.NullString{String: *req.Phone, Valid: true}
	}

	// Create feedback record
	err := s.interact.CreateCustomerFeedback(ctx, db.CreateCustomerFeedbackParams{
		ID:       feedbackID,
		UserID:   userIDNull,
		Email:    emailNull,
		Phone:    phoneNull,
		Category: req.Category,
		Content:  req.Content,
	})

	if err != nil {
		return "", assets_services.NewError(500, fmt.Errorf("lỗi khi lưu feedback: %w", err))
	}

	return feedbackID, nil
}

func (s *service) GetCustomerFeedbacks(ctx context.Context, req *entity.GetCustomerFeedbacksRequest) ([]entity.CustomerFeedbackItem, *assets_services.ServiceError) {
	// Pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	offset := int32((req.Page - 1) * req.PageSize)
	limit := int32(req.PageSize)

	// Convert filters to sql.Null types
	var categoryNull, userIDNull sql.NullString
	var startDateNull, endDateNull sql.NullTime

	if req.Category != nil && *req.Category != "" {
		categoryNull = sql.NullString{String: *req.Category, Valid: true}
	}
	if req.UserID != nil && *req.UserID != "" {
		userIDNull = sql.NullString{String: *req.UserID, Valid: true}
	}
	if req.StartDate != nil && *req.StartDate != "" {
		if t, err := parseDate(*req.StartDate); err == nil {
			startDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if t, err := parseDate(*req.EndDate); err == nil {
			endDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	rows, err := s.interact.ListCustomerFeedbacks(ctx, db.ListCustomerFeedbacksParams{
		CategoryFilter: categoryNull,
		UserIDFilter:   userIDNull,
		StartDate:      startDateNull,
		EndDate:        endDateNull,
		Limit:          limit,
		Offset:         offset,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy danh sách feedback: %w", err))
	}

	result := make([]entity.CustomerFeedbackItem, 0, len(rows))
	for _, row := range rows {
		item := entity.CustomerFeedbackItem{
			ID:        row.ID,
			Category:  row.Category,
			Content:   row.Content,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}

		if row.UserID.Valid {
			item.UserID = &row.UserID.String
		}
		if row.Email.Valid {
			item.Email = &row.Email.String
		}
		if row.Phone.Valid {
			item.Phone = &row.Phone.String
		}

		result = append(result, item)
	}

	return result, nil
}

func (s *service) GetCustomerFeedbackByID(ctx context.Context, id string) (*entity.CustomerFeedbackItem, *assets_services.ServiceError) {
	row, err := s.interact.GetCustomerFeedbackByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, assets_services.NewError(404, fmt.Errorf("không tìm thấy feedback"))
		}
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy feedback: %w", err))
	}

	item := &entity.CustomerFeedbackItem{
		ID:        row.ID,
		Category:  row.Category,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	if row.UserID.Valid {
		item.UserID = &row.UserID.String
	}
	if row.Email.Valid {
		item.Email = &row.Email.String
	}
	if row.Phone.Valid {
		item.Phone = &row.Phone.String
	}

	return item, nil
}

func (s *service) GetCustomerFeedbackStats(ctx context.Context, startDate, endDate *string) (*entity.CustomerFeedbackStatsResponse, *assets_services.ServiceError) {
	// Parse dates to time.Time
	var startDateNull, endDateNull sql.NullTime
	if startDate != nil && *startDate != "" {
		if t, err := parseDate(*startDate); err == nil {
			startDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if endDate != nil && *endDate != "" {
		if t, err := parseDate(*endDate); err == nil {
			endDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	stats, err := s.interact.GetCustomerFeedbackStats(ctx, db.GetCustomerFeedbackStatsParams{
		StartDate: startDateNull,
		EndDate:   endDateNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy thống kê feedback: %w", err))
	}

	return &entity.CustomerFeedbackStatsResponse{
		TotalFeedbacks:  stats.TotalFeedbacks,
		BugCount:        stats.BugCount,
		ComplaintCount:  stats.ComplaintCount,
		SuggestionCount: stats.SuggestionCount,
		UniqueUsers:     stats.UniqueUsers,
	}, nil
}

func (s *service) GetCustomerFeedbacksByCategory(ctx context.Context, startDate, endDate *string) ([]entity.CustomerFeedbackCategoryStats, *assets_services.ServiceError) {
	// Parse dates to time.Time
	var startDateNull, endDateNull sql.NullTime
	if startDate != nil && *startDate != "" {
		if t, err := parseDate(*startDate); err == nil {
			startDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}
	if endDate != nil && *endDate != "" {
		if t, err := parseDate(*endDate); err == nil {
			endDateNull = sql.NullTime{Time: t, Valid: true}
		}
	}

	rows, err := s.interact.GetCustomerFeedbacksByCategory(ctx, db.GetCustomerFeedbacksByCategoryParams{
		StartDate: startDateNull,
		EndDate:   endDateNull,
	})

	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy thống kê theo category: %w", err))
	}

	result := make([]entity.CustomerFeedbackCategoryStats, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.CustomerFeedbackCategoryStats{
			Category:      row.Category,
			FeedbackCount: row.FeedbackCount,
		})
	}

	return result, nil
}
