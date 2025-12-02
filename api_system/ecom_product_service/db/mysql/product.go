package db

import (
	"context"
	"database/sql"
	"strings"

	db "github.com/TranVinhHien/ecom_product_service/db/sqlc"
)

// Hàm hỗ trợ build WHERE clause chung cho cả List và Count
func buildWhereClause(params db.ListProductsAdvancedParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// 1. Điều kiện cứng
	// conditions = append(conditions, "p.delete_status = 'Active'")
	if params.DeleteStatus.Valid {
		conditions = append(conditions, "p.delete_status = ?")
		args = append(args, params.DeleteStatus.ProductDeleteStatus)
	}
	// 2. Các điều kiện động (Chỉ thêm nếu có giá trị)
	if params.ShopID.Valid {
		conditions = append(conditions, "p.shop_id = ?")
		args = append(args, params.ShopID.String)
	}
	if params.CategoryID.Valid {
		conditions = append(conditions, "p.category_id = ?")
		args = append(args, params.CategoryID.String)
	}
	if params.BrandID.Valid {
		conditions = append(conditions, "p.brand_id = ?")
		args = append(args, params.BrandID.String)
	}
	if params.PriceMin.Valid {
		conditions = append(conditions, "p.min_price >= ?")
		args = append(args, params.PriceMin.Float64)
	}
	if params.PriceMax.Valid {
		conditions = append(conditions, "p.min_price <= ?")
		args = append(args, params.PriceMax.Float64)
	}

	// Xử lý Keyword (interface{} -> string)
	if kw, ok := params.Keyword.(sql.NullString); ok && kw.Valid {
		// Tối ưu: Dùng LIKE
		conditions = append(conditions, "p.name LIKE ?")
		args = append(args, "%"+kw.String+"%")
	}

	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

// ============================================================
// 1. HÀM LIST PRODUCTS (Dynamic)
// ============================================================
func (q *SQLStore) ListProductsDynamic(ctx context.Context, params db.ListProductsAdvancedParams) ([]db.ListProductsAdvancedRow, error) {
	// A. Build SELECT
	baseQuery := `
		SELECT 
			p.id, p.name, p.` + "`key`" + `, p.description, p.short_description,
			p.brand_id, p.category_id, p.shop_id, p.image, p.media,
			p.delete_status, p.product_is_permission_return, p.product_is_permission_check,
			p.create_date, p.update_date, p.create_by, p.update_by,
			p.total_sold, p.min_price, p.max_price,
			(SELECT ps.id FROM product_sku ps WHERE ps.product_id = p.id ORDER BY ps.price ASC LIMIT 1) AS min_price_sku_id,
			(SELECT ps.id FROM product_sku ps WHERE ps.product_id = p.id ORDER BY ps.price DESC LIMIT 1) AS max_price_sku_id
		FROM product p
	`

	// B. Build WHERE
	whereClause, args := buildWhereClause(params)

	// C. Build ORDER BY (Whitelist để tránh SQL Injection)
	orderBy := "ORDER BY p.create_date DESC" // Mặc định
	if sortStr, ok := params.Sort.(sql.NullString); ok && sortStr.Valid {
		switch sortStr.String {
		case "best_sell":
			orderBy = "ORDER BY p.total_sold DESC"
		case "price_asc":
			orderBy = "ORDER BY p.min_price ASC"
		case "price_desc":
			orderBy = "ORDER BY p.min_price DESC"
		case "name_asc":
			orderBy = "ORDER BY p.name ASC"
		case "name_desc":
			orderBy = "ORDER BY p.name DESC"
		}
	}

	// D. Build LIMIT/OFFSET
	limitOffset := " LIMIT ? OFFSET ?"
	args = append(args, params.Limit, params.Offset)

	// E. Final Query Assembly
	finalQuery := baseQuery + " " + whereClause + " " + orderBy + limitOffset

	// F. Execute
	rows, err := q.connPool.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// G. Scan Results
	var items []db.ListProductsAdvancedRow
	for rows.Next() {
		var i db.ListProductsAdvancedRow
		if err := rows.Scan(
			&i.ID, &i.Name, &i.Key, &i.Description, &i.ShortDescription,
			&i.BrandID, &i.CategoryID, &i.ShopID, &i.Image, &i.Media,
			&i.DeleteStatus, &i.ProductIsPermissionReturn, &i.ProductIsPermissionCheck,
			&i.CreateDate, &i.UpdateDate, &i.CreateBy, &i.UpdateBy,
			&i.TotalSold, &i.MinPrice, &i.MaxPrice,
			&i.MinPriceSkuID, &i.MaxPriceSkuID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// ============================================================
// 2. HÀM COUNT PRODUCTS (Dynamic)
// ============================================================
func (q *SQLStore) CountProductsDynamic(ctx context.Context, params db.ListProductsAdvancedParams) (int64, error) {
	// A. Build SELECT
	baseQuery := "SELECT COUNT(*) FROM product p"

	// B. Build WHERE (Tái sử dụng logic để đảm bảo đồng nhất)
	whereClause, args := buildWhereClause(params)

	// C. Final Query
	finalQuery := baseQuery + " " + whereClause

	// D. Execute
	var count int64
	err := q.connPool.QueryRowContext(ctx, finalQuery, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
