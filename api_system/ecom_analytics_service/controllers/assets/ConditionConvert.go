package controllers_assets

import (
	"errors"
	"strconv"
	"strings"

	services "github.com/TranVinhHien/ecom_analytics_service/services/entity"
)

func ParseConditions(search string) ([]services.Condition, error) {
	var conditions []services.Condition

	// Định nghĩa các toán tử cho phép, lưu ý thứ tự ưu tiên ">=, <=" trước ">" "<" "="
	operatorList := []string{">=", "<=", ">", "<", "="}

	// Tách các điều kiện riêng biệt bằng dấu phẩy
	parts := strings.Split(search, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		var chosenOp string
		var opIndex int = -1

		// Tìm toán tử có trong chuỗi với thứ tự ưu tiên đã định
		for _, op := range operatorList {
			if idx := strings.Index(part, op); idx != -1 {
				chosenOp = op
				opIndex = idx
				break
			}
		}

		if chosenOp == "" {
			return nil, errors.New("không tìm thấy toán tử hợp lệ trong: " + part)
		}

		// Lấy ra trường (field) và giá trị (value) dựa trên vị trí của toán tử
		field := strings.TrimSpace(part[:opIndex])
		valueStr := strings.TrimSpace(part[opIndex+len(chosenOp):])
		if field == "" || valueStr == "" {
			return nil, errors.New("thiếu field hoặc value trong: " + part)
		}

		// Cố gắng chuyển đổi giá trị sang số nếu có thể, nếu không thì giữ nguyên dưới dạng chuỗi.
		var value interface{} = valueStr
		if intVal, err := strconv.Atoi(valueStr); err == nil {
			value = intVal
		} else if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
			value = floatVal
		}

		cond := services.Condition{
			Field:    field,
			Operator: chosenOp,
			Value:    value,
		}
		conditions = append(conditions, cond)
	}

	return conditions, nil
}
