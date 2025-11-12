package assets_services

import (
	"strings"
)

func ConvertSliceToQuotedString(slice []string) string {
	// Trích xuất Product_sku_id từ slice

	// Chuyển thành chuỗi đúng định dạng
	productSkuIDString := strings.Join(slice, ",")
	return productSkuIDString
}
