package assets_services

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
)

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
