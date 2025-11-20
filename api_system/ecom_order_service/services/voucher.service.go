package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
	assets_services "github.com/TranVinhHien/ecom_order_service/services/assets"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"
	"github.com/google/uuid"
)

func (s *service) CreateVoucher(ctx context.Context, req services.CreateVoucherRequest, shop_id, user_type string) *assets_services.ServiceError {
	// 1. Validate toàn bộ dữ liệu đầu vào
	if err := s.validateCreateVoucherRequest(req); err != nil {
		return &assets_services.ServiceError{
			Code: 400,
			Err:  err,
		}
	}

	// 2. Tạo UUID
	voucherID := uuid.New().String()
	owner := shop_id
	max_discount_amount := 0.0
	if req.MaxDiscountAmount != nil {
		max_discount_amount = *req.MaxDiscountAmount
	}
	var ownerType db.VouchersOwnerType
	if user_type == "ROLE_ADMIN" {
		ownerType = db.VouchersOwnerTypePLATFORM
		owner = s.env.PlatformOwnerID

	} else if user_type == "ROLE_SELLER" {
		ownerType = db.VouchersOwnerTypeSHOP
	} else {
		return &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user_type không hợp lệ để tạo voucher"),
		}
	}

	// 3. Map DTO (Request) sang Params (sqlc)
	// Lưu ý: Giả định sqlc.yaml của bạn map ENUM của MySQL thành String của Go
	params := db.CreateVoucherParams{
		ID:                voucherID,
		Name:              req.Name,
		VoucherCode:       req.VoucherCode,
		OwnerType:         ownerType,
		OwnerID:           owner,
		DiscountType:      db.VouchersDiscountType(req.DiscountType),
		DiscountValue:     fmt.Sprintf("%.2f", req.DiscountValue),
		MaxDiscountAmount: sql.NullString{String: fmt.Sprintf("%.2f", max_discount_amount), Valid: db.VouchersDiscountType(req.DiscountType) == db.VouchersDiscountTypePERCENTAGE},
		AppliesToType:     db.VouchersAppliesToType(req.AppliesToType),
		MinPurchaseAmount: fmt.Sprintf("%.2f", req.MinPurchaseAmount),
		AudienceType:      db.VouchersAudienceType(req.AudienceType),
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		TotalQuantity:     req.TotalQuantity,
		MaxUsagePerUser:   req.MaxUsagePerUser,
		IsActive:          true,
	}

	// 4. Gọi DB
	err := s.repository.CreateVoucher(ctx, params)
	if err != nil {
		// Ở đây bạn có thể check lỗi (ví dụ: lỗi duplicate `voucher_code`)
		return &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("lỗi khi tạo voucher: %w", err),
		}
	}

	return nil
}

// --- Hàm Sửa Voucher (Partial Update) ---

func (s *service) UpdateVoucher(ctx context.Context, voucherID string, user_id string, user_type string, req services.UpdateVoucherRequest) *assets_services.ServiceError {
	// 1. Map DTO (Request) sang Params (sqlc)
	// Vì dùng `sqlc.narg`, struct Params của `UpdateVoucher` sẽ dùng `sql.Null*`
	// Chúng ta chỉ set giá trị `Valid: true` cho những trường KHÔNG PHẢI nil trong request

	params := db.UpdateVoucherParams{
		ID: voucherID,
	}
	voucher_db, err := s.repository.GetVoucherByID(ctx, voucherID)
	if err != nil {
		return &assets_services.ServiceError{
			Code: 404,
			Err:  fmt.Errorf("không tìm thấy voucher với ID %s: %w", voucherID, err)}
	}

	if voucher_db.OwnerType == db.VouchersOwnerTypePLATFORM && user_type == "ROLE_SELLER" {
		return &assets_services.ServiceError{
			Code: 403,
			Err:  fmt.Errorf("bạn không có quyền sửa voucher này"),
		}
	}
	if voucher_db.OwnerID != user_id && user_type == "ROLE_SELLER" {
		return &assets_services.ServiceError{
			Code: 403,
			Err:  fmt.Errorf("bạn không có quyền sửa voucher này"),
		}
	}
	// Ánh xạ các trường string
	if req.Name != nil {
		params.Name = sql.NullString{String: *req.Name, Valid: true}
	}
	if req.VoucherCode != nil {
		params.VoucherCode = sql.NullString{String: *req.VoucherCode, Valid: true}
	}
	if req.DiscountType != nil {
		params.DiscountType = db.NullVouchersDiscountType{VouchersDiscountType: db.VouchersDiscountType(*req.DiscountType), Valid: true}
	}
	if req.AppliesToType != nil {
		params.AppliesToType = db.NullVouchersAppliesToType{VouchersAppliesToType: db.VouchersAppliesToType(*req.AppliesToType), Valid: true}
	}
	if req.AudienceType != nil {
		params.AudienceType = db.NullVouchersAudienceType{VouchersAudienceType: db.VouchersAudienceType(*req.AudienceType), Valid: true}
	}

	// Ánh xạ các trường float64
	if req.DiscountValue != nil {
		params.DiscountValue = sql.NullString{String: fmt.Sprintf("%.2f", *req.DiscountValue), Valid: true}
	}
	if req.MaxDiscountAmount != nil {
		params.MaxDiscountAmount = sql.NullString{String: fmt.Sprintf("%.2f", *req.MaxDiscountAmount), Valid: true}
	}
	if req.MinPurchaseAmount != nil {
		params.MinPurchaseAmount = sql.NullString{String: fmt.Sprintf("%.2f", *req.MinPurchaseAmount), Valid: true}
	}

	// Ánh xạ các trường int32
	if req.TotalQuantity != nil {
		params.TotalQuantity = sql.NullInt32{Int32: *req.TotalQuantity, Valid: true}
	}
	if req.MaxUsagePerUser != nil {
		params.MaxUsagePerUser = sql.NullInt32{Int32: *req.MaxUsagePerUser, Valid: true}
	}

	// Ánh xạ các trường time
	if req.StartDate != nil {
		params.StartDate = sql.NullTime{Time: *req.StartDate, Valid: true}
	}
	if req.EndDate != nil {
		params.EndDate = sql.NullTime{Time: *req.EndDate, Valid: true}
	}
	if req.IsActive != nil {

		params.IsActive = sql.NullBool{Bool: *req.IsActive, Valid: true}
	}
	// 2. Gọi DB
	// Nhờ `COALESCE` trong SQL, các trường `Valid: false` (mặc định) sẽ bị bỏ qua
	err = s.repository.UpdateVoucher(ctx, params)
	if err != nil {
		// Check lỗi (ví dụ: không tìm thấy voucher_id, hoặc duplicate voucher_code mới)
		return &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("lỗi khi cập nhật voucher: %w", err),
		}
	}

	return nil
}

// --- Hàm Lấy Danh Sách Voucher Cho User ---

func (s *service) ListVouchersForUser(ctx context.Context, userID string, filter services.VoucherFilterRequest) (map[string]interface{}, *assets_services.ServiceError) {
	var publicVouchers []db.Vouchers
	var assignedVouchers []db.Vouchers
	var err error

	// Set default sort_by nếu không được truyền vào
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	// Validate sort_by
	validSortBy := map[string]bool{
		"discount_asc":  true,
		"discount_desc": true,
		"created_at":    true,
	}
	if !validSortBy[filter.SortBy] {
		return nil, &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid sort_by value. Allowed: discount_asc, discount_desc, created_at"),
		}
	}

	// Validate owner_type nếu có
	if filter.OwnerType != nil {
		validOwnerType := map[string]bool{
			"PLATFORM": true,
			"SHOP":     true,
		}
		if !validOwnerType[*filter.OwnerType] {
			return nil, &assets_services.ServiceError{
				Code: 400,
				Err:  fmt.Errorf("invalid owner_type value. Allowed: PLATFORM, SHOP"),
			}
		}
	}

	// Validate applies_to_type nếu có
	if filter.AppliesToType != nil {
		validAppliesToType := map[string]bool{
			"ORDER_TOTAL":  true,
			"SHIPPING_FEE": true,
		}
		if !validAppliesToType[*filter.AppliesToType] {
			return nil, &assets_services.ServiceError{
				Code: 400,
				Err:  fmt.Errorf("invalid applies_to_type value. Allowed: ORDER_TOTAL, SHIPPING_FEE"),
			}
		}
	}

	// 1. Lấy voucher CÔNG KHAI với bộ lọc
	if filter.OwnerType != nil || filter.ShopID != nil || filter.AppliesToType != nil || filter.SortBy != "created_at" {
		// Sử dụng query có filter
		params := db.GetPublicVouchersWithFilterParams{
			SortBy: filter.SortBy,
		}
		if filter.OwnerType != nil {
			params.OwnerType = db.NullVouchersOwnerType{
				VouchersOwnerType: db.VouchersOwnerType(*filter.OwnerType),
				Valid:             true,
			}
		}
		if filter.ShopID != nil {
			params.ShopID = sql.NullString{String: *filter.ShopID, Valid: true}
		}
		if filter.AppliesToType != nil {
			params.AppliesToType = db.NullVouchersAppliesToType{
				VouchersAppliesToType: db.VouchersAppliesToType(*filter.AppliesToType),
				Valid:                 true,
			}
		}

		publicVouchers, err = s.repository.GetPublicVouchersWithFilter(ctx, params)
	} else {
		// Sử dụng query không có filter (mặc định)
		publicVouchers, err = s.repository.GetPublicVouchers(ctx)
	}

	if err != nil {
		return nil, &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("lỗi khi lấy voucher công khai: %w", err),
		}
	}

	// 2. Lấy voucher RIÊNG (Assigned) với bộ lọc
	if filter.OwnerType != nil || filter.ShopID != nil || filter.AppliesToType != nil || filter.SortBy != "created_at" {
		// Sử dụng query có filter
		params := db.GetAssignedVouchersByUserWithFilterParams{
			UserID: userID,
			SortBy: filter.SortBy,
		}
		if filter.OwnerType != nil {
			params.OwnerType = db.NullVouchersOwnerType{
				VouchersOwnerType: db.VouchersOwnerType(*filter.OwnerType),
				Valid:             true,
			}
		}
		if filter.ShopID != nil {
			params.ShopID = sql.NullString{String: *filter.ShopID, Valid: true}
		}
		if filter.AppliesToType != nil {
			params.AppliesToType = db.NullVouchersAppliesToType{
				VouchersAppliesToType: db.VouchersAppliesToType(*filter.AppliesToType),
				Valid:                 true,
			}
		}

		assignedVouchers, err = s.repository.GetAssignedVouchersByUserWithFilter(ctx, params)
	} else {
		// Sử dụng query không có filter (mặc định)
		assignedVouchers, err = s.repository.GetAssignedVouchersByUser(ctx, userID)
	}

	if err != nil {
		return nil, &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("lỗi khi lấy voucher được gán cho người dùng: %w", err),
		}
	}

	// 3. Gộp cả hai danh sách và loại bỏ trùng lặp (nếu có)
	// Dùng map để đảm bảo ID voucher là duy nhất
	combinedVouchers := make([]db.Vouchers, 0)
	combinedVouchers = append(combinedVouchers, publicVouchers...)
	combinedVouchers = append(combinedVouchers, assignedVouchers...)

	// 3.5 check người dùng còn được dùng voucher này không
	for i, v := range combinedVouchers {
		// Mặc định là hợp lệ
		isValid := true

		// Gọi Repository đếm số lần user này đã dùng voucher này
		// (Sử dụng query CountVoucherUsageByUser đã có trong sqlc)
		usageCount, err := s.repository.CountVoucherUsageByUser(ctx, db.CountVoucherUsageByUserParams{
			VoucherID: v.ID,
			UserID:    userID,
		})

		if err != nil {
			// Nếu lỗi DB, log lại và tạm thời coi như user chưa dùng (hoặc return lỗi tuỳ chính sách)
			// Ở đây tôi chọn continue để không làm gãy cả danh sách
			fmt.Printf("Error checking usage for voucher %s: %v\n", v.ID, err)
			usageCount = 0
		}

		// Logic kiểm tra: Đã dùng >= Giới hạn cho phép => Hết lượt (Invalid)
		if int32(usageCount) >= v.MaxUsagePerUser {
			isValid = false
		}

		// (Tuỳ chọn) Kiểm tra thêm điều kiện tổng quan: Voucher đã hết lượt dùng toàn hệ thống chưa?
		if v.UsedQuantity >= v.TotalQuantity {
			isValid = false
		}
		combinedVouchers[i].IsActive = isValid

	}

	// 4. Chuyển map về slice

	// result := assets_services.NormalizeListSQLNulls(combinedVouchers, "data")
	return map[string]interface{}{
		"data": combinedVouchers,
	}, nil
}

// --- Hàm Lấy Danh Sách Voucher Cho Admin/Seller Quản Lý ---
func (s *service) ListVouchersForManagement(ctx context.Context, ownerID string, ownerType string, filter services.VoucherManagementFilterRequest) (map[string]interface{}, *assets_services.ServiceError) {
	// 1. Validate ownerType
	if ownerType != "PLATFORM" && ownerType != "SHOP" {
		return nil, &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("owner_type không hợp lệ. Chỉ chấp nhận: PLATFORM, SHOP"),
		}
	}

	// 2. Validate và set default cho pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100 // Giới hạn tối đa 100 items/trang
	}

	// 3. Validate sort_by
	validSortBy := map[string]bool{
		"created_at_desc": true,
		"created_at_asc":  true,
		"start_date_desc": true,
		"start_date_asc":  true,
		"end_date_desc":   true,
		"end_date_asc":    true,
	}
	if filter.SortBy == "" {
		filter.SortBy = "created_at_desc" // Mặc định: mới nhất trước
	}
	if !validSortBy[filter.SortBy] {
		return nil, &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("sort_by không hợp lệ. Allowed: created_at_desc, created_at_asc, start_date_desc, start_date_asc, end_date_desc, end_date_asc"),
		}
	}

	// 4. Validate các filters
	if filter.DiscountType != nil {
		validDiscountTypes := map[string]bool{
			"PERCENTAGE":   true,
			"FIXED_AMOUNT": true,
		}
		if !validDiscountTypes[*filter.DiscountType] {
			return nil, &assets_services.ServiceError{
				Code: 400,
				Err:  fmt.Errorf("discount_type không hợp lệ. Allowed: PERCENTAGE, FIXED_AMOUNT"),
			}
		}
	}

	if filter.AppliesToType != nil {
		validAppliesToTypes := map[string]bool{
			"ORDER_TOTAL":  true,
			"SHIPPING_FEE": true,
		}
		if !validAppliesToTypes[*filter.AppliesToType] {
			return nil, &assets_services.ServiceError{
				Code: 400,
				Err:  fmt.Errorf("applies_to_type không hợp lệ. Allowed: ORDER_TOTAL, SHIPPING_FEE"),
			}
		}
	}

	if filter.AudienceType != nil {
		validAudienceTypes := map[string]bool{
			"PUBLIC":   true,
			"ASSIGNED": true,
		}
		if !validAudienceTypes[*filter.AudienceType] {
			return nil, &assets_services.ServiceError{
				Code: 400,
				Err:  fmt.Errorf("audience_type không hợp lệ. Allowed: PUBLIC, ASSIGNED"),
			}
		}
	}

	if filter.Status != nil {
		validStatuses := map[string]bool{
			"ACTIVE":   true,
			"EXPIRED":  true,
			"UPCOMING": true,
			"DEPLETED": true,
		}
		if !validStatuses[*filter.Status] {
			return nil, &assets_services.ServiceError{
				Code: 400,
				Err:  fmt.Errorf("status không hợp lệ. Allowed: ACTIVE, EXPIRED, UPCOMING, DEPLETED"),
			}
		}
	}

	// 5. Build params cho sqlc queries
	vouchers, total, err := s.getVouchersForManagementUsingSqlc(ctx, ownerID, ownerType, filter)
	if err != nil {
		return nil, &assets_services.ServiceError{
			Code: 500,
			Err:  fmt.Errorf("lỗi khi lấy danh sách voucher: %w", err),
		}
	}

	// 6. Tính toán pagination metadata
	totalPages := (total + int64(filter.PageSize) - 1) / int64(filter.PageSize)

	result := map[string]interface{}{
		"data": vouchers,
		"pagination": map[string]interface{}{
			"current_page": filter.Page,
			"page_size":    filter.PageSize,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	}

	return result, nil
}

// getVouchersForManagementUsingSqlc sử dụng sqlc generated queries
func (s *service) getVouchersForManagementUsingSqlc(ctx context.Context, ownerID string, ownerType string, filter services.VoucherManagementFilterRequest) ([]map[string]interface{}, int64, error) {
	// 1. Build params cho Count query
	if db.VouchersOwnerType(ownerType) == db.VouchersOwnerTypePLATFORM {
		ownerID = s.env.PlatformOwnerID
	}
	countParams := db.CountVouchersForManagementParams{
		OwnerID:   ownerID,
		OwnerType: db.VouchersOwnerType(ownerType),
	}

	// Add optional filters với LIKE pattern cho search
	if filter.VoucherCode != nil && *filter.VoucherCode != "" {
		voucherCodePattern := "%" + *filter.VoucherCode + "%"
		countParams.VoucherCode = sql.NullString{String: voucherCodePattern, Valid: true}
	}
	if filter.Name != nil && *filter.Name != "" {
		namePattern := "%" + *filter.Name + "%"
		countParams.Name = sql.NullString{String: namePattern, Valid: true}
	}
	if filter.DiscountType != nil {
		countParams.DiscountType = db.NullVouchersDiscountType{
			VouchersDiscountType: db.VouchersDiscountType(*filter.DiscountType),
			Valid:                true,
		}
	}
	if filter.AppliesToType != nil {
		countParams.AppliesToType = db.NullVouchersAppliesToType{
			VouchersAppliesToType: db.VouchersAppliesToType(*filter.AppliesToType),
			Valid:                 true,
		}
	}
	if filter.AudienceType != nil {
		countParams.AudienceType = db.NullVouchersAudienceType{
			VouchersAudienceType: db.VouchersAudienceType(*filter.AudienceType),
			Valid:                true,
		}
	}
	if filter.IsActive != nil {
		countParams.IsActive = sql.NullBool{Bool: *filter.IsActive, Valid: true}
	}
	if filter.Status != nil {
		countParams.Status = *filter.Status
	}

	// 2. Get total count
	total, err := s.repository.CountVouchersForManagement(ctx, countParams)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting vouchers: %w", err)
	}

	// 3. Build params cho List query
	offset := int32((filter.Page - 1) * filter.PageSize)
	limitVal := int32(filter.PageSize)

	// Prepare status value
	var statusValue interface{} = nil
	if filter.Status != nil {
		statusValue = *filter.Status
	}

	// 4. Call appropriate query based on sort_by
	var vouchersDB []db.Vouchers

	switch filter.SortBy {
	case "created_at_asc":
		vouchersDB, err = s.repository.ListVouchersForManagementBySortCreatedAtAsc(ctx, db.ListVouchersForManagementBySortCreatedAtAscParams{
			OwnerID:       ownerID,
			OwnerType:     db.VouchersOwnerType(ownerType),
			VoucherCode:   countParams.VoucherCode,
			Name:          countParams.Name,
			DiscountType:  countParams.DiscountType,
			AppliesToType: countParams.AppliesToType,
			AudienceType:  countParams.AudienceType,
			IsActive:      countParams.IsActive,
			Status:        statusValue,
			Limit:         limitVal,
			Offset:        offset,
		})
	case "start_date_desc":
		vouchersDB, err = s.repository.ListVouchersForManagementBySortStartDateDesc(ctx, db.ListVouchersForManagementBySortStartDateDescParams{
			OwnerID:       ownerID,
			OwnerType:     db.VouchersOwnerType(ownerType),
			VoucherCode:   countParams.VoucherCode,
			Name:          countParams.Name,
			DiscountType:  countParams.DiscountType,
			AppliesToType: countParams.AppliesToType,
			AudienceType:  countParams.AudienceType,
			IsActive:      countParams.IsActive,
			Status:        statusValue,
			Limit:         limitVal,
			Offset:        offset,
		})
	case "start_date_asc":
		vouchersDB, err = s.repository.ListVouchersForManagementBySortStartDateAsc(ctx, db.ListVouchersForManagementBySortStartDateAscParams{
			OwnerID:       ownerID,
			OwnerType:     db.VouchersOwnerType(ownerType),
			VoucherCode:   countParams.VoucherCode,
			Name:          countParams.Name,
			DiscountType:  countParams.DiscountType,
			AppliesToType: countParams.AppliesToType,
			AudienceType:  countParams.AudienceType,
			IsActive:      countParams.IsActive,
			Status:        statusValue,
			Limit:         limitVal,
			Offset:        offset,
		})
	case "end_date_desc":
		vouchersDB, err = s.repository.ListVouchersForManagementBySortEndDateDesc(ctx, db.ListVouchersForManagementBySortEndDateDescParams{
			OwnerID:       ownerID,
			OwnerType:     db.VouchersOwnerType(ownerType),
			VoucherCode:   countParams.VoucherCode,
			Name:          countParams.Name,
			DiscountType:  countParams.DiscountType,
			AppliesToType: countParams.AppliesToType,
			AudienceType:  countParams.AudienceType,
			IsActive:      countParams.IsActive,
			Status:        statusValue,
			Limit:         limitVal,
			Offset:        offset,
		})
	case "end_date_asc":
		vouchersDB, err = s.repository.ListVouchersForManagementBySortEndDateAsc(ctx, db.ListVouchersForManagementBySortEndDateAscParams{
			OwnerID:       ownerID,
			OwnerType:     db.VouchersOwnerType(ownerType),
			VoucherCode:   countParams.VoucherCode,
			Name:          countParams.Name,
			DiscountType:  countParams.DiscountType,
			AppliesToType: countParams.AppliesToType,
			AudienceType:  countParams.AudienceType,
			IsActive:      countParams.IsActive,
			Status:        statusValue,
			Limit:         limitVal,
			Offset:        offset,
		})
	default: // created_at_desc
		vouchersDB, err = s.repository.ListVouchersForManagementBySortCreatedAtDesc(ctx, db.ListVouchersForManagementBySortCreatedAtDescParams{
			OwnerID:       ownerID,
			OwnerType:     db.VouchersOwnerType(ownerType),
			VoucherCode:   countParams.VoucherCode,
			Name:          countParams.Name,
			DiscountType:  countParams.DiscountType,
			AppliesToType: countParams.AppliesToType,
			AudienceType:  countParams.AudienceType,
			IsActive:      countParams.IsActive,
			Status:        statusValue,
			Limit:         limitVal,
			Offset:        offset,
		})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("error fetching vouchers: %w", err)
	}

	// 5. Convert to response format with calculated fields
	vouchers := make([]map[string]interface{}, 0, len(vouchersDB))
	for _, v := range vouchersDB {
		status := s.calculateVoucherStatus(v)

		voucherMap := map[string]interface{}{
			"id":                  v.ID,
			"name":                v.Name,
			"voucher_code":        v.VoucherCode,
			"owner_type":          v.OwnerType,
			"owner_id":            v.OwnerID,
			"discount_type":       v.DiscountType,
			"discount_value":      v.DiscountValue,
			"max_discount_amount": v.MaxDiscountAmount,
			"applies_to_type":     v.AppliesToType,
			"min_purchase_amount": v.MinPurchaseAmount,
			"audience_type":       v.AudienceType,
			"start_date":          v.StartDate,
			"end_date":            v.EndDate,
			"total_quantity":      v.TotalQuantity,
			"used_quantity":       v.UsedQuantity,
			"remaining_quantity":  v.TotalQuantity - v.UsedQuantity,
			"max_usage_per_user":  v.MaxUsagePerUser,
			"is_active":           v.IsActive,
			"status":              status,
			"created_at":          v.CreatedAt,
			"updated_at":          v.UpdatedAt,
		}

		vouchers = append(vouchers, voucherMap)
	}

	return vouchers, total, nil
}

// calculateVoucherStatus tính toán trạng thái hiện tại của voucher
func (s *service) calculateVoucherStatus(v db.Vouchers) string {
	now := time.Now()

	// Check depleted first
	if v.UsedQuantity >= v.TotalQuantity {
		return "DEPLETED"
	}

	// Check if expired
	if v.EndDate.Before(now) {
		return "EXPIRED"
	}

	// Check if upcoming
	if v.StartDate.After(now) {
		return "UPCOMING"
	}

	// Check if active
	if v.IsActive && v.StartDate.Before(now) && v.EndDate.After(now) {
		return "ACTIVE"
	}

	return "INACTIVE"
}

//================================================================
// 4. HÀM CHECK VOUCHER (CẬP NHẬT)
//================================================================

// --- Hàm Check Nhiều Voucher Cùng Lúc ---
// (Hàm này chạy song song để tối ưu tốc độ check)
type CheckResult struct {
	Voucher db.Vouchers `json:"Voucher"`
	IsValid bool        `json:"isValid"`
	Reason  string      `json:"reason"`
}

type VoucherCheckInput struct {
	VoucherID       string  `json:"voucher_id"`
	TotalOrderPrice float64 `json:"total_order_price"`
	// VoucherType     voucherCheckInputType `json:"voucher_type"`
}

func (s *service) CheckVouchers(ctx context.Context, userID string, voucherIDs []string) ([]CheckResult, *assets_services.ServiceError) {

	resultsChan := make(chan CheckResult, len(voucherIDs))
	var wg sync.WaitGroup

	for _, voucher := range voucherIDs {
		wg.Add(1)

		// Thực hiện check mỗi voucher trong một goroutine riêng
		go func(voucher string) {
			defer wg.Done()

			// Gọi hàm check MỘT voucher (xem hàm bên dưới)
			isValid, reason, voucherDB := s.checkSingleVoucher(ctx, userID, voucher)

			resultsChan <- CheckResult{
				Voucher: voucherDB,
				IsValid: isValid,
				Reason:  reason,
			}
		}(voucher)
	}

	// Đợi tất cả goroutine hoàn thành
	wg.Wait()
	close(resultsChan)

	// Thu thập kết quả
	results := make([]CheckResult, 0, len(voucherIDs))
	for res := range resultsChan {
		results = append(results, res)
	}

	return results, nil
}

// Ensure checkSingleVoucher correctly uses the TotalOrderPrice passed in VoucherCheckInput
func (s *service) checkSingleVoucher(ctx context.Context, userID string, voucherInput string) (isValid bool, reason string, voucher db.Vouchers) {

	// 1. Get voucher info (check existence, active, date, total quantity)
	voucher, err := s.repository.GetVoucherByIDForValidation(ctx, voucherInput)
	if err != nil {
		// ... (error handling as before) ...
		if err == sql.ErrNoRows {
			return false, fmt.Sprintf("Voucher không tồn tại, đã hết hoặc hết hạn %s", voucherInput), voucher
		}
		log.Printf("Lỗi DB khi check voucher %s: %v", voucherInput, err)
		return false, "Lỗi hệ thống", voucher
	}

	// 2. Check Audience (as before)
	if voucher.AudienceType == "ASSIGNED" {
		// ... (check user_vouchers status as before) ...
		userVoucher, err := s.repository.GetUserVoucherStatus(ctx, db.GetUserVoucherStatusParams{
			VoucherID: voucherInput,
			UserID:    userID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return false, "Bạn không sở hữu voucher này.", voucher
			}
			log.Printf("Lỗi DB khi check user_voucher %s: %v", voucherInput, err)
			return false, "Lỗi hệ thống", voucher
		}
		if userVoucher.Status != "AVAILABLE" {
			return false, fmt.Sprintf("Voucher đã %s.", userVoucher.Status), voucher
		}
	} else if voucher.AudienceType != "PUBLIC" {
		return false, "Loại voucher không hợp lệ.", voucher
	}

	// 3. Check Personal Usage Limit (as before)
	count, err := s.repository.CountVoucherUsageByUser(ctx, db.CountVoucherUsageByUserParams{
		VoucherID: voucherInput,
		UserID:    userID,
	})
	if err != nil {
		// ... (error handling as before) ...
		log.Printf("Lỗi DB khi đếm usage voucher %s: %v", voucherInput, err)
		return false, "Lỗi hệ thống", voucher
	}
	if count >= int64(voucher.MaxUsagePerUser) {
		return false, fmt.Sprintf("Bạn đã sử dụng hết số lần cho voucher %s.", voucher.Name), voucher
	}

	// 5. If all checks pass, the voucher is valid *for this context*
	return true, "OK", voucher
}

// =================================================================
// HÀM SỬ DỤNG VOUCHER (GỌI KHI TẠO ĐƠN)
// =================================================================

// UseVoucher xử lý việc sử dụng MỘT voucher.
// Hàm này gọi wrapper ExecTS để đảm bảo tất cả các lệnh DB
// (check, update, insert) diễn ra trong cùng 1 transaction.
func (s *service) UseVoucher(ctx context.Context, input services.UseVoucherInput) error {

	// Gọi hàm transaction wrapper
	return s.repository.ExecTS(ctx, func(tx db.Querier) error {
		// 1. Check voucher có hợp lệ hay không (sử dụng tx)
		isValid, reason := s.checkSingleVoucherTx(ctx, tx, input.UserID, input.VoucherID)
		if !isValid {
			return fmt.Errorf("voucher ID %s không hợp lệ: %s", input.VoucherID, reason)
		}

		// 2. Lấy thông tin voucher (sử dụng tx)
		voucher, err := tx.GetVoucherByID(ctx, input.VoucherID)
		if err != nil {
			return fmt.Errorf("không thể lấy thông tin voucher %s: %w", input.VoucherID, err)
		}

		// 3. Tăng số lượng đã dùng (Atomically) (sử dụng tx)
		// *** SỬA LỖI Ở ĐÂY ***
		// Hàm :execrows của sqlc trả về (sql.Result, error)
		_, err = tx.IncrementVoucherUsage(ctx, input.VoucherID)
		if err != nil {
			return fmt.Errorf("lỗi DB khi cập nhật số lượng voucher %s: %w", input.VoucherID, err)
		}
		// Gọi RowsAffected() từ sql.Result
		// rowsAffected, err := res.RowsAffected()
		// if err != nil {
		// 	return fmt.Errorf("lỗi khi lấy RowsAffected: %w", err)
		// }
		// if rowsAffected != 1 {
		// 	// Đây là lỗi RACE CONDITION
		// 	return fmt.Errorf("voucher %s đã hết lượt ngay khi bạn sử dụng", voucher.VoucherCode)
		// }

		// 4. Ghi lại lịch sử sử dụng (sử dụng tx)
		historyParams := db.CreateVoucherUsageHistoryParams{
			VoucherID:      input.VoucherID,
			UserID:         input.UserID,
			DiscountAmount: fmt.Sprintf("%.2f", input.DiscountAmount),
		}
		if err := tx.CreateVoucherUsageHistory(ctx, historyParams); err != nil {
			return fmt.Errorf("lỗi DB khi ghi lịch sử voucher %s: %w", input.VoucherID, err)
		}

		// 5. Nếu là voucher ASSIGNED, cập nhật trạng thái trong ví (sử dụng tx)
		if voucher.AudienceType == "ASSIGNED" {
			statusParams := db.SetUserVoucherStatusParams{
				VoucherID: input.VoucherID,
				UserID:    input.UserID,
				Status:    "USED", // Chuyển sang USED
			}
			_, err := tx.SetUserVoucherStatus(ctx, statusParams)
			if err != nil {
				return fmt.Errorf("lỗi DB khi cập nhật ví user_voucher %s: %w", input.VoucherID, err)
			}
			// // rowsAffected, _ := res.RowsAffected() // Check nếu cập nhật thành công
			// if rowsAffected != 1 {
			// 	// Lỗi logic: User dùng voucher assigned mà không có trong ví (status=AVAILABLE)
			// 	return fmt.Errorf("lỗi không nhất quán: voucher %s không 'AVAILABLE' trong ví", input.VoucherID)
			// }
		}

		return nil // Thành công -> ExecTS sẽ Commit
	})
}

// =================================================================
// HÀM TRẢ VOUCHER (GỌI KHI HỦY ĐƠN)
// =================================================================

// RollbackVoucher xử lý việc hoàn trả MỘT voucher.
// Hàm này cũng gọi ExecTS để đảm bảo tính nhất quán.
func (s *service) RollbackVoucher(ctx context.Context, input services.RollbackVoucherInput) error {

	return s.repository.ExecTS(ctx, func(tx db.Querier) error {
		// 1. Lấy thông tin voucher (để check audience_type) (sử dụng tx)
		voucher, err := tx.GetVoucherByID(ctx, input.VoucherID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("Cảnh báo: Rollback voucher %s nhưng voucher không còn tồn tại.", input.VoucherID)
				return nil // Bỏ qua, không báo lỗi
			}
			return fmt.Errorf("lỗi DB khi lấy voucher %s: %w", input.VoucherID, err)
		}

		// 2. Tìm và Xóa lịch sử sử dụng (sử dụng tx)
		history, err := tx.GetVoucherUsageHistory(ctx, db.GetVoucherUsageHistoryParams{
			VoucherID: input.VoucherID,
			UserID:    input.UserID,
		})

		if err != nil {
			if err == sql.ErrNoRows {
				// Lịch sử đã được rollback rồi. Không cần làm gì thêm.
				log.Printf("Cảnh báo: Rollback voucher %s cho đơn  nhưng không tìm thấy lịch sử sử dụng.", input.VoucherID)
				return nil
			}
			return fmt.Errorf("lỗi DB khi tìm lịch sử voucher %s: %w", input.VoucherID, err)
		}

		// Nếu tìm thấy, tiến hành xóa (sử dụng tx)
		if _, err := tx.DeleteVoucherUsageHistory(ctx, history.ID); err != nil {
			return fmt.Errorf("lỗi DB khi xóa lịch sử voucher %s: %w", input.VoucherID, err)
		}

		// 3. Giảm số lượng đã dùng (cộng trả lại) (sử dụng tx)
		if _, err := tx.DecrementVoucherUsage(ctx, input.VoucherID); err != nil {
			// Lỗi này nghiêm trọng, cần log lại
			log.Printf("LỖI NGHIÊM TRỌNG: Không thể decrement voucher %s: %v", input.VoucherID, err)
			// Không return error để các voucher khác tiếp tục rollback (nếu gọi hàm này trong vòng lặp)
		}

		// 4. Nếu là voucher ASSIGNED, reset trạng thái trong ví (sử dụng tx)
		if voucher.AudienceType == "ASSIGNED" {
			resetParams := db.ResetUserVoucherStatusParams{
				VoucherID: input.VoucherID,
				UserID:    input.UserID,
			}
			if _, err := tx.ResetUserVoucherStatus(ctx, resetParams); err != nil {
				log.Printf("Cảnh báo: Lỗi DB khi reset ví user_voucher %s: %v", input.VoucherID, err)
			}
		}

		return nil // Thành công -> ExecTS sẽ Commit
	})
}

// =================================================================
// HÀM HELPER CHECK (ĐÃ SỬA ĐỂ DÙNG `tx db.Querier`)
// =================================================================

// checkSingleVoucherTx là hàm helper nội bộ, nhận vào 1 querier (có thể là DB hoặc TX)
func (s *service) checkSingleVoucherTx(ctx context.Context, q db.Querier, userID string, voucherID string) (isValid bool, reason string) {

	// 1. Lấy thông tin voucher (check tồn tại, active, ngày, số lượng tổng)
	voucher, err := q.GetVoucherByIDForValidation(ctx, voucherID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Sprintf("Voucher không tồn tại, đã hết hoặc hết hạn: %s", voucherID)
		}
		log.Printf("Lỗi DB khi check voucher %s: %v", voucherID, err)
		return false, "Lỗi hệ thống"
	}

	// 2. Check logic dựa trên ĐỐI TƯỢNG (Audience)
	if voucher.AudienceType == "ASSIGNED" {
		userVoucher, err := q.GetUserVoucherStatus(ctx, db.GetUserVoucherStatusParams{
			VoucherID: voucherID,
			UserID:    userID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return false, "Bạn không sở hữu voucher này."
			}
			log.Printf("Lỗi DB khi check user_voucher %s: %v", voucherID, err)
			return false, "Lỗi hệ thống"
		}

		if userVoucher.Status != "AVAILABLE" {
			return false, fmt.Sprintf("Voucher đã %s.", userVoucher.Status)
		}

	} else if voucher.AudienceType != "PUBLIC" {
		return false, "Loại voucher không hợp lệ."
	}
	// Nếu là "PUBLIC", bỏ qua check sở hữu

	// 3. Check LƯỢT SỬ DỤNG CÁ NHÂN (Áp dụng cho cả PUBLIC và ASSIGNED)
	count, err := q.CountVoucherUsageByUser(ctx, db.CountVoucherUsageByUserParams{
		VoucherID: voucherID,
		UserID:    userID,
	})
	if err != nil {
		log.Printf("Lỗi DB khi đếm usage voucher %s: %v", voucherID, err)
		return false, "Lỗi hệ thống"
	}

	if count >= int64(voucher.MaxUsagePerUser) {
		return false, "Bạn đã sử dụng hết số lần cho voucher này."
	}

	return true, "OK"
}

// =================================================================
// HÀM VALIDATE TẠO VOUCHER
// =================================================================

// validateCreateVoucherRequest kiểm tra tính hợp lệ của toàn bộ dữ liệu tạo voucher
func (s *service) validateCreateVoucherRequest(req services.CreateVoucherRequest) error {
	// 1. Validate Name
	if req.Name == "" {
		return fmt.Errorf("tên voucher không được để trống")
	}
	if len(req.Name) > 255 {
		return fmt.Errorf("tên voucher không được vượt quá 255 ký tự")
	}

	// 2. Validate VoucherCode
	if req.VoucherCode == "" {
		return fmt.Errorf("mã voucher không được để trống")
	}
	if len(req.VoucherCode) > 50 {
		return fmt.Errorf("mã voucher không được vượt quá 50 ký tự")
	}

	// 5. Validate DiscountType
	validDiscountTypes := map[string]bool{
		"PERCENTAGE":   true,
		"FIXED_AMOUNT": true,
	}
	if !validDiscountTypes[req.DiscountType] {
		return fmt.Errorf("discount_type không hợp lệ. Chỉ chấp nhận: PERCENTAGE, FIXED_AMOUNT")
	}

	// 6. Validate DiscountValue
	if req.DiscountValue <= 0 {
		return fmt.Errorf("discount_value phải lớn hơn 0")
	}

	// Nếu là PERCENTAGE, giá trị phải từ 0-100
	if req.DiscountType == "PERCENTAGE" {
		if req.DiscountValue > 100 {
			return fmt.Errorf("discount_value cho PERCENTAGE không được vượt quá 100")
		}
	}

	// 7. Validate MaxDiscountAmount
	// Nếu discount_type là PERCENTAGE, max_discount_amount BẮT BUỘC phải có
	if req.DiscountType == "PERCENTAGE" {
		if req.MaxDiscountAmount == nil {
			return fmt.Errorf("max_discount_amount bắt buộc phải có khi discount_type là PERCENTAGE")
		}
		if *req.MaxDiscountAmount <= 0 {
			return fmt.Errorf("max_discount_amount phải lớn hơn 0")
		}
	}

	// 8. Validate AppliesToType
	validAppliesToTypes := map[string]bool{
		"ORDER_TOTAL":  true,
		"SHIPPING_FEE": true,
	}
	if !validAppliesToTypes[req.AppliesToType] {
		return fmt.Errorf("applies_to_type không hợp lệ. Chỉ chấp nhận: ORDER_TOTAL, SHIPPING_FEE")
	}

	// 9. Validate MinPurchaseAmount
	if req.MinPurchaseAmount < 0 {
		return fmt.Errorf("min_purchase_amount không được âm")
	}

	// 10. Validate AudienceType
	validAudienceTypes := map[string]bool{
		"PUBLIC":   true,
		"ASSIGNED": true,
	}
	if !validAudienceTypes[req.AudienceType] {
		return fmt.Errorf("audience_type không hợp lệ. Chỉ chấp nhận: PUBLIC, ASSIGNED")
	}

	// 11. Validate UserUse (danh sách user_id) - BẮT BUỘC khi audience_type là ASSIGNED
	if req.AudienceType == "ASSIGNED" {
		if len(req.UserUse) == 0 {
			return fmt.Errorf("audience_type là ASSIGNED phải có danh sách user_use")
		}
		// Kiểm tra các user_id không rỗng
		for i, userID := range req.UserUse {
			if userID == "" {
				return fmt.Errorf("user_use[%d] không được để trống", i)
			}
		}
	}

	// 12. Validate StartDate và EndDate
	if req.StartDate.IsZero() {
		return fmt.Errorf("start_date không được để trống")
	}
	if req.EndDate.IsZero() {
		return fmt.Errorf("end_date không được để trống")
	}
	if req.EndDate.Before(req.StartDate) {
		return fmt.Errorf("end_date phải sau start_date")
	}
	if req.StartDate.Before(req.StartDate.Add(-24 * 365 * 10 * time.Hour)) {
		// Kiểm tra start_date không quá xa trong quá khứ (10 năm)
		return fmt.Errorf("start_date không hợp lệ")
	}

	// 13. Validate TotalQuantity
	if req.TotalQuantity <= 0 {
		return fmt.Errorf("total_quantity phải lớn hơn 0")
	}

	// 14. Validate MaxUsagePerUser
	if req.MaxUsagePerUser <= 0 {
		return fmt.Errorf("max_usage_per_user phải lớn hơn 0")
	}
	// max_usage_per_user không được lớn hơn total_quantity
	if req.MaxUsagePerUser > req.TotalQuantity {
		return fmt.Errorf("max_usage_per_user không được lớn hơn total_quantity")
	}

	return nil
}
