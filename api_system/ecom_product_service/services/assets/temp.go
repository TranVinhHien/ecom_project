package assets_services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"

	"github.com/google/uuid"
)

// normalizeValue is a recursive helper to normalize values, handling nulls deeply
func normalizeValue(input interface{}) interface{} {
	if input == nil {
		return nil
	}

	v := reflect.ValueOf(input)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	t := v.Type()
	kind := v.Kind()

	switch kind {
	case reflect.Struct:
		// Handle SQL null types explicitly
		switch val := v.Interface().(type) {
		case sql.NullString:
			if val.Valid {
				return val.String
			}
			return nil
		case sql.NullBool:
			if val.Valid {
				return val.Bool
			}
			return nil
		case sql.NullTime:
			if val.Valid {
				return val.Time
			}
			return nil
		case sql.NullInt64:
			if val.Valid {
				return val.Int64
			}
			return nil
		case sql.NullInt32:
			if val.Valid {
				return val.Int32
			}
			return nil
		case sql.NullFloat64:
			if val.Valid {
				return val.Float64
			}
			return nil
		}

		// Check if it's a null-like struct with "Valid" field
		validField := v.FieldByName("Valid")
		if validField.IsValid() && validField.Kind() == reflect.Bool {
			if !validField.Bool() {
				return nil
			}
			// Find and return the data field (first non-Valid field)
			for i := 0; i < v.NumField(); i++ {
				if t.Field(i).Name != "Valid" {
					// Recurse on the data field for deeper handling
					return normalizeValue(v.Field(i).Interface())
				}
			}
			return nil // No data field found
		}

		// Normal struct: recurse into fields and build map
		result := make(map[string]interface{})
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = fieldType.Name
			}
			// Skip unexported fields
			if fieldType.PkgPath != "" {
				continue
			}
			result[jsonTag] = normalizeValue(field.Interface())
		}
		return result

	case reflect.Slice, reflect.Array:
		length := v.Len()
		res := make([]interface{}, length)
		for i := 0; i < length; i++ {
			res[i] = normalizeValue(v.Index(i).Interface())
		}
		return res

	case reflect.Map:
		res := make(map[interface{}]interface{})
		iter := v.MapRange()
		for iter.Next() {
			res[iter.Key().Interface()] = normalizeValue(iter.Value().Interface())
		}
		return res

	default:
		// Primitives and other types (e.g., time.Time): return as is
		return v.Interface()
	}
}

// NormalizeListSQLNulls converts a slice or array to a map with proper null handling, recursively
func NormalizeListSQLNulls(input interface{}, key string) map[string]interface{} {
	var items []interface{}

	if input == nil {
		return map[string]interface{}{key: []map[string]interface{}{}}
	}

	rv := reflect.ValueOf(input)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return map[string]interface{}{key: []map[string]interface{}{}}
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		n := rv.Len()
		items = make([]interface{}, n)
		for i := 0; i < n; i++ {
			items[i] = rv.Index(i).Interface()
		}
	default:
		items = []interface{}{rv.Interface()}
	}

	list_result := make([]map[string]interface{}, 0, len(items))

	for _, item := range items {
		normalized := normalizeValue(item)
		if normalizedMap, ok := normalized.(map[string]interface{}); ok {
			list_result = append(list_result, normalizedMap)
		} else {
			// If the top-level item normalizes to a non-map (e.g., a primitive or null), handle accordingly
			list_result = append(list_result, map[string]interface{}{})
		}
	}
	return map[string]interface{}{key: list_result}
}

// NormalizeSQLNulls converts a single struct to a map with proper null handling, recursively
func NormalizeSQLNulls(input interface{}, key string) map[string]interface{} {
	if input == nil {
		return map[string]interface{}{key: map[string]interface{}{}}
	}

	// Normalize the input using the recursive helper
	normalized := normalizeValue(input)

	// Ensure the result is a map[string]interface{}
	result, ok := normalized.(map[string]interface{})
	if !ok {
		// If the normalized result is not a map (e.g., a primitive or null), wrap it in an empty map
		return map[string]interface{}{key: map[string]interface{}{}}
	}

	return map[string]interface{}{key: result}
}

// HideFields ẩn các field được chỉ định từ đối tượng hoặc mảng/slice và trả về dưới dạng map với key được cung cấp
func HideFields(obj interface{}, key string, fieldsToHide ...string) (map[string]interface{}, error) {
	// Convert object to JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Khởi tạo kết quả trả về
	result := make(map[string]interface{})

	// Xử lý trường hợp đầu vào là slice/array
	if reflect.TypeOf(obj).Kind() == reflect.Slice || reflect.TypeOf(obj).Kind() == reflect.Array {
		var slice []map[string]interface{}
		if err := json.Unmarshal(data, &slice); err != nil {
			return nil, err
		}

		// Ẩn các field chỉ định trong từng phần tử
		for i := range slice {
			for _, field := range fieldsToHide {
				delete(slice[i], field)
			}
		}

		// Gán slice vào key trong result
		result[key] = slice
		return result, nil
	}

	// Xử lý trường hợp đầu vào là đối tượng đơn (struct hoặc map)
	var singleObj map[string]interface{}
	if err := json.Unmarshal(data, &singleObj); err != nil {
		return nil, err
	}

	// Ẩn các field chỉ định
	for _, field := range fieldsToHide {
		delete(singleObj, field)
	}

	// Nếu key rỗng, trả về trực tiếp singleObj
	if key == "" {
		return singleObj, nil
	}

	// Gán singleObj vào key trong result
	result[key] = singleObj
	return result, nil
}

func SaveUploadedFile(fileHeader *multipart.FileHeader, destination string) error {
	// Mở file gốc từ form
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("mở file lỗi: %v", err)
	}
	defer src.Close()
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("không thể tạo thư mục lưu ảnh: %v", err)
	}
	// Tạo file đích để lưu
	dst, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("không thể tạo file: %v", err)
	}
	defer dst.Close()

	// Ghi nội dung từ src -> dst
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("ghi file lỗi: %v", err)
	}

	return nil
}

func SaveFile(fileHeader *multipart.FileHeader, dstDir string) (*string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// đảm bảo thư mục tồn tại
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return nil, err
	}
	file_name := fileHeader.Filename + "-" + uuid.NewString() + filepath.Ext(fileHeader.Filename)
	dstPath := filepath.Join(dstDir, file_name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}
	return &file_name, nil
}
func DeleteFile(dir, fileName string) error {
	filePath := filepath.Join(dir, fileName)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("không thể xóa file %s: %w", filePath, err)
	}
	return nil
}
